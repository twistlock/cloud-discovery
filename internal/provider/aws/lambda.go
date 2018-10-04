package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/twistlock/cloud-discovery/internal/shared"
)

type lambdaClient struct {
	opt AWSOptions
}

// NewServerlessClientAWS create a new aws serverless client
func NewLambdaClient(opt AWSOptions) *lambdaClient {
	return &lambdaClient{
		opt: opt,
	}
}

// Functions retrieves a list of functions for the settings provided to the client
func (s *lambdaClient) Discover() (*shared.CloudDiscoveryResult, error) {
	var res []shared.CloudAsset
	session, err := CreateAWSSession(&s.opt)
	if err != nil {
		return nil, err
	}
	client := lambda.New(session)
	if err != nil {
		return nil, err
	}

	var nextMarker string // A string returned by aws api if there are more functions to quiry after list functions is called
	truncated := true     // A flag to indicate if the list result is partial and more quries are needed
	for truncated {
		functionVersion := "ALL" // Scan all versions of a function
		input := &lambda.ListFunctionsInput{FunctionVersion: &functionVersion}
		if nextMarker != "" {
			input.Marker = &nextMarker
		}
		functions, err := client.ListFunctions(input)
		if err != nil {
			return nil, err
		}

		if functions == nil {
			return nil, fmt.Errorf("received nil function list from AWS %s", s.opt.Region)
		}
		for _, f := range functions.Functions {
			res = append(res, shared.CloudAsset{ID:         *f.FunctionName})
		}

		if functions.NextMarker == nil {
			truncated = false
			break
		}
		nextMarker = *functions.NextMarker
	}

	return &shared.CloudDiscoveryResult{Assets:res, Region: s.opt.Region, Type:"Lambda"}, nil
}
