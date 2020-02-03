package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	numargs := len(os.Args[1:])
	if numargs < 1 {
		log.Fatalln("Need at least one ARN of a Lambda layer, sorry :(")
	}
	switch numargs {
	case 1:
		// the ARN of the Lambda layer has to be the first argument:
		larns := os.Args[1]
		// look up metadata and content:
		linfo, larn, err := resolve(larns)
		if err != nil {
			log.Fatalf("Can't diagnose Lambda layer based on the ARN %s: %v", larns, err)
		}
		err = render(larn, linfo)
		if err != nil {
			log.Fatalf("Can't resolve Lambda layer location: %v", err)
		}
	default:
		log.Println("render all layers")
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
	p := message.NewPrinter(language.English)
	p.Printf("Size: %v kB\n", *linfo.Content.CodeSize/1024)
	lloc, err := url.Parse(*linfo.Content.Location)
	if err != nil {
		return err
	}
	q := lloc.Query()
	fmt.Printf("Location: %v://%v%v?versionId=%v\n", lloc.Scheme, lloc.Host, lloc.Path, q.Get("versionId"))
	return nil
}
