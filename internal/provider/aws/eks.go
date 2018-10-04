package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/twistlock/cloud-discovery/internal/shared"
)

type eksClient struct {
	options AWSOptions
}

func NewEKSClient(options AWSOptions) *eksClient {
	return &eksClient{options: options}
}

// Clusters returns a list of endpoints representing EKS clusters
func (c *eksClient) Discover() (result *shared.CloudDiscoveryResult, err error) {
	// Setup aws eks service object
	sess, err := CreateAWSSession(&c.options)
	if err != nil {
		return nil, err
	}
	svc := eks.New(sess)

	// List regional EKS endpoints names
	var clusters []*string
	out, err := svc.ListClusters(nil)
	if err != nil {
		return nil, err
	}
	clusters = append(clusters, out.Clusters...)
	for out.NextToken != nil {
		out, err = svc.ListClusters(&eks.ListClustersInput{NextToken: out.NextToken})
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, out.Clusters...)
	}

	// For each cluster name, fetch its information
	var endpoints []shared.CloudAsset
	for _, cluster := range clusters {
		out, err := svc.DescribeCluster(&eks.DescribeClusterInput{Name: cluster})
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, shared.CloudAsset{ ID: aws.StringValue(out.Cluster.Name), })
	}
	fmt.Println("Found AWS EKS endpoints", len(endpoints))
	return &shared.CloudDiscoveryResult{Assets:endpoints, Region: c.options.Region, Type:"EKS"}, nil
}
