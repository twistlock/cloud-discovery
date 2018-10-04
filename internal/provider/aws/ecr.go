package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/twistlock/cloud-discovery/internal/shared"
)

type ecrClient struct {
	opt AWSOptions
}

func NewECRClient(opt AWSOptions) *ecrClient {
	return &ecrClient{opt:opt}
}

func (r *ecrClient) Discover() (*shared.CloudDiscoveryResult, error) {
	session, err := CreateAWSSession(&r.opt)
	ecr.New(session)
	client := ecr.New(session)
	if err != nil {
		return nil, err
	}

	result := shared.CloudDiscoveryResult{
		Region: r.opt.Region,
		Type:"ECR",
	}
	var nextToken *string
	for {
		out, err := client.DescribeRepositories(&ecr.DescribeRepositoriesInput{NextToken: nextToken})
		if err != nil {
			return nil, err
		} else if out == nil {
			return nil, fmt.Errorf("failed to query AWS repositories")
		}
		for _, repo := range out.Repositories {
			result.Assets = append(result.Assets, shared.CloudAsset{ID: aws.StringValue(repo.RepositoryName)})
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return &result, nil
}