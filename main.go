package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	layers := flag.String("layers", "", "A comma-separated list of ARN(s) of the layer(s) you want to inspect or export.")
	exportpath := flag.String("export", "", "[OPTIONAL] The directory path to export the layer metadata or content to, if not provided, the layer metadata will be printed to stdout, only.")
	flag.Parse()
	if *layers == "" {
		log.Fatalln("Need at least one ARN of a Lambda layer, sorry :(")
	}
	switch {
	case strings.Contains(*layers, ","): // multiple layer ARNs provided
		layerarns := strings.Split(*layers, ",")
		err := renderall(layerarns)
		if err != nil {
			log.Fatalf("Can't render provided Lambda layers: %v", err)
		}
	default: // a single layer ARN provided
		larns := *layers
		// look up metadata and content of the layer:
		linfo, larn, err := resolve(larns)
		if err != nil {
			log.Fatalf("Can't diagnose Lambda layer based on the ARN %s: %v", larns, err)
		}
		err = render(larn, linfo)
		if err != nil {
			log.Fatalf("Can't resolve Lambda layer location: %v", err)
		}
		if *exportpath != "" {
			absoluteep, err := download(*linfo.Content.Location, *exportpath)
			if err != nil {
				log.Fatalf("%v", err)
			}
			fmt.Printf("Content exported to: %v\n", absoluteep)
		}
	}
}

// resolve looks up metadata and content of a AWS Lambda layer by ARN
func resolve(larns string) (*lambda.GetLayerVersionByArnOutput, arn.ARN, error) {
	larn, err := arn.Parse(larns)
	if err != nil {
		return nil, larn, err
	}
	svc := lambda.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion(larn.Region))
	vo, err := svc.GetLayerVersionByArn(
		&lambda.GetLayerVersionByArnInput{
			Arn: aws.String(larn.String()),
		})
	return vo, larn, nil
}

// render displays info about a single Lambda layer
func render(larn arn.ARN, linfo *lambda.GetLayerVersionByArnOutput) error {
	fmt.Printf("Name: %v\n", strings.Split(larn.Resource, ":")[1])
	fmt.Printf("Version: %v\n", *linfo.Version)
	fmt.Printf("Description: %v\n", *linfo.Description)
	fmt.Printf("Created on: %v\n", *linfo.CreatedDate)
	message.NewPrinter(language.English).Printf("Size: %v kB\n", *linfo.Content.CodeSize/1024)
	lloc, err := url.Parse(*linfo.Content.Location)
	if err != nil {
		return err
	}
	q := lloc.Query()
	fmt.Printf("Location: %v://%v%v?versionId=%v\n", lloc.Scheme, lloc.Host, lloc.Path, q.Get("versionId"))
	return nil
}

// renderall displays tabular infos about multiple Lambda layers
func renderall(larnslist []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION\tCREATED ON\tSIZE (kB)")
	for _, larns := range larnslist {
		// look up metadata and content of the layer:
		laclean := strings.TrimSpace(larns)
		linfo, larn, err := resolve(laclean)
		if err != nil {
			log.Fatalf("Can't diagnose Lambda layer based on the ARN %s: %v", larns, err)
		}
		lname := fmt.Sprintf("%v\t", strings.Split(larn.Resource, ":")[1])
		lversion := fmt.Sprintf("%d\t", *linfo.Version)
		ldesc := fmt.Sprintf("%v\t", *linfo.Description)
		lcreatedon := fmt.Sprintf("%v\t", *linfo.CreatedDate)
		lsize := message.NewPrinter(language.English).Sprintf("%v", *linfo.Content.CodeSize/1024)
		fmt.Fprintln(w, lname+lversion+ldesc+lcreatedon+lsize)
	}
	w.Flush()
	return nil
}

// download dereferences the URL and writes its content into the path provided
func download(url, exportpath string) (string, error) {
	// download layer content:
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// create export file path and generate absolute version of it:
	absepathdir, err := filepath.Abs(exportpath)
	if err != nil {
		return "", fmt.Errorf("Can't generate absolute path of %s whilst exporting layer content: %v", absepathdir, err)
	}
	// write layer content to local (temp) ZIP file:
	lczipfile := filepath.Join(absepathdir, "layer-content.zip")
	os.MkdirAll(absepathdir, 0755)
	out, err := os.Create(lczipfile)
	if err != nil {
		return "", fmt.Errorf("Can't create ZIP file %s for exporting layer content: %v", lczipfile, err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Can't write to ZIP file %s for exporting layer content: %v", lczipfile, err)
	}
	// unzip layer contents into directory:
	lcdir := filepath.Join(absepathdir, "layer-content")
	err = unzip(lczipfile, lcdir)
	if err != nil {
		return "", fmt.Errorf("Can't unzip file %s for exporting layer content: %v", lczipfile, err)
	}
	// clean up by removing the ZIP file of the layer contents:
	err = os.Remove(lczipfile)
	if err != nil {
		return "", fmt.Errorf("Can't delete ZIP file %s whilst exporting layer content: %v", lczipfile, err)
	}
	return lcdir, nil
}

// unzip extracts src ZIP file into dest directory,
// and creates the directory if it doesn't exist,
// from https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func unzip(zipfile, dest string) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	os.MkdirAll(dest, 0755)
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}
