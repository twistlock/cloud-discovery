package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"time"
)

type ecrClient struct {
	opt Options
}

// NewECRClient creates a new ECR client (Amazon Elastic Container Registry)
func NewECRClient(opt Options) *ecrClient {
	return &ecrClient{opt: opt}
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
		Type:   "ECR",
	}
	var nextToken *string
	for {
		out, err := client.DescribeRepositories(&ecr.DescribeRepositoriesInput{NextToken: nextToken})
		if err != nil {
			return nil, err
		} else if out == nil {
			return nil, fmt.Errorf("failed to describe repositories")
		}
		for _, repo := range out.Repositories {
			result.Assets = append(result.Assets, shared.CloudAsset{ID: aws.StringValue(repo.RepositoryName), Data: struct {
				ARN           *string    `json:"arn"`
				CreatedAt     *time.Time `json:"createdAt"`
				RegistryId    *string    `json:"registryId"`
				RepositoryArn *string    `json:"repositoryArn"`
				RepositoryUri *string    `json:"repositoryUri"`
				Version       *string    `json:"version"`
			}{
				CreatedAt:     repo.CreatedAt,
				RegistryId:    repo.RegistryId,
				RepositoryArn: repo.RepositoryArn,
				RepositoryUri: repo.RepositoryUri,
			}})
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return &result, nil
}
