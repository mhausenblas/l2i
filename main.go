package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func main() {
	if len(os.Args[1:]) < 1 {
		log.Fatalln("Need the ARN of the Lambda layer sorry :(")
	}
	// the ARN of the Lambda layer has to be the first argument:
	layerarn := os.Args[1]
	// look up metadata and content:
	linfo, err := resolve(layerarn)
	if err != nil {
		log.Fatalf("Can't diagnose Lambda layer based on the ARN %s", layerarn)
	}
	log.Printf("%v", linfo)
}

// resolve looks up metadata and content of a AWS Lambda layer by ARN
func resolve(layerarn string) (*lambda.GetLayerVersionByArnOutput, error) {
	layer, err := arn.Parse(layerarn)
	if err != nil {
		return nil, err
	}
	svc := lambda.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion(layer.Region))
	vo, err := svc.GetLayerVersionByArn(
		&lambda.GetLayerVersionByArnInput{
			Arn: aws.String(layer.String()),
		})
	return vo, nil
}
