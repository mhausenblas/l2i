package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
			exportdir, err := ioutil.TempDir(*exportpath, "l2i")
			if err != nil {
				log.Fatalf("Can't create export directory: %v", err)
			}
			ep, err := filepath.Abs(exportdir)
			if err != nil {
				log.Fatalf("Can't make export directory absolute: %v", err)
			}
			fmt.Printf("Content exported to: %v\n", ep)
			// tmpfn := filepath.Join(dir, "tmpfile")
			// if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
			// log.Fatalf("Can't export layer info: %v", err)
			// }
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

func flagpresent(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
