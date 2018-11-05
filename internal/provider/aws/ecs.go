package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/twistlock/cloud-discovery/internal/shared"
)

type ecsClient struct {
	opt Options
}

// NewECSClient creates a new ECS (Amazon Elastic Container Service) client
func NewECSClient(opt Options) *ecsClient {
	return &ecsClient{opt: opt}
}

func (c *ecsClient) Discover() (*shared.CloudDiscoveryResult, error) {
	sess, err := CreateAWSSession(&c.opt)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	// List regional ECS cluster arns
	var nextToken *string
	var arns []*string
	for {
		out, err := svc.ListClusters(&ecs.ListClustersInput{NextToken: nextToken})
		if err != nil {
			return nil, err
		}
		arns = append(arns, out.ClusterArns...)
		nextToken = out.NextToken
		if nextToken == nil {
			break
		}
	}

	// Fetch information for each resource identifier
	// 1. Cluster name
	// 2. Number of instances (EC2 hosts)
	var clusters []shared.CloudAsset
	for _, arn := range arns {
		name, err := c.fetchClusterName(svc, arn)
		if err != nil {
			return nil, fmt.Errorf("error fetching ECS cluster name: %v", err)
		}

		// Fetch list of hosts
		var hosts []string
		var nextToken *string
		for {
			out, err := svc.ListContainerInstances(&ecs.ListContainerInstancesInput{Cluster: arn, NextToken: nextToken})
			if err != nil {
				return nil, fmt.Errorf("error listing ECS cluster instance. =%v instances: %v", arn, err)
			}
			for _, arn := range out.ContainerInstanceArns {
				hosts = append(hosts, *arn)
			}
			nextToken = out.NextToken
			if nextToken == nil {
				break
			}
		}
		clusters = append(clusters, shared.CloudAsset{ID: name, Data: struct {
			Hosts []string `json:"hosts"`
		}{Hosts: hosts}})
	}
	return &shared.CloudDiscoveryResult{Assets: clusters, Region: c.opt.Region, Type: "ECS"}, nil
}

func (c *ecsClient) fetchClusterName(svc *ecs.ECS, arn *string) (string, error) {
	out, err := svc.DescribeClusters(&ecs.DescribeClustersInput{Clusters: []*string{arn}})
	if err != nil {
		return "", err
	}
	if len(out.Clusters) < 1 || out.Clusters[0].ClusterName == nil {
		return "", fmt.Errorf("unknown ECS cluster: %s", *arn)
	}
	return aws.StringValue(out.Clusters[0].ClusterName), nil
}
