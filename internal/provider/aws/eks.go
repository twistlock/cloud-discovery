package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"time"
)

type eksClient struct {
	options AWSOptions
}

// NewEKSClient creates a new EKS (Amazon Elastic Kubernetes Service) client
func NewEKSClient(options AWSOptions) *eksClient {
	return &eksClient{options: options}
}

func (c *eksClient) Discover() (result *shared.CloudDiscoveryResult, err error) {
	sess, err := CreateAWSSession(&c.options)
	if err != nil {
		return nil, err
	}
	svc := eks.New(sess)

	var clusters []*string
	var nextToken *string
	for {
		out, err := svc.ListClusters(&eks.ListClustersInput{NextToken: nextToken})
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, out.Clusters...)
		nextToken = out.NextToken
		if nextToken == nil {
			break
		}
	}
	// For each cluster name, fetch its information
	var assets []shared.CloudAsset
	for _, cluster := range clusters {
		out, err := svc.DescribeCluster(&eks.DescribeClusterInput{Name: cluster})
		if err != nil {
			return nil, err
		}
		assets = append(assets, shared.CloudAsset{ID: aws.StringValue(out.Cluster.Name),
			Data: struct {
				ARN       *string    `json:"arn"`
				CreatedAt *time.Time `json:"createdAt"`
				Endpoint  *string    `json:"endpoint"`
				RoleArn   *string    `json:"roleArn"`
				Status    *string    `json:"status"`
				Version   *string    `json:"version"`
			}{
				ARN:       out.Cluster.Arn,
				CreatedAt: out.Cluster.CreatedAt,
				RoleArn:   out.Cluster.RoleArn,
				Status:    out.Cluster.Status,
				Version:   out.Cluster.Version,
			},
		})
	}
	return &shared.CloudDiscoveryResult{Assets: assets, Region: c.options.Region, Type: "EKS"}, nil
}
