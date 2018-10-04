package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/twistlock/cloud-discovery/internal/shared"
)

type ecsClient struct {
	opt AWSOptions
}

func NewECSClient(opt AWSOptions) *ecsClient {
	return &ecsClient{opt:opt}
}

// Clusters returns a list of endpoints representing EKS clusters
func (c *ecsClient) Discover() (*shared.CloudDiscoveryResult, error) {
	// Setup aws ecs service object
	sess, err := CreateAWSSession(&c.opt)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	// List regional ECS cluster arns
	out, err := svc.ListClusters(nil)
	if err != nil {
		return nil, err
	}
	arns := out.ClusterArns
	for out.NextToken != nil {
		out, err = svc.ListClusters(&ecs.ListClustersInput{NextToken: out.NextToken})
		if err != nil {
			return nil, err
		}
		arns = append(arns, out.ClusterArns...)
	}

	var clusters []shared.CloudAsset

	// For each cluster arn, fetch its information:
	// 1. Cluster name
	// 2. Number of instances (EC2 hosts)
	for _, arn := range arns {
		name, err := c.fetchClusterName(svc, arn)
		if err != nil {
			return nil, fmt.Errorf("error fetching ECS cluster name: %v", err)
		}

		// NOTE: the aws ecs ListContainerInstances api call returns a list of
		// ec2 hosts (not containers)

		var instances  []string
		out, err := svc.ListContainerInstances(&ecs.ListContainerInstancesInput{Cluster: arn})
		if err != nil {
			return nil, fmt.Errorf("error listing ECS cluster=%s instances: %v", *arn, err)
		}
		for _, arn := range out.ContainerInstanceArns {
			instances = append(instances, *arn)
		}
		for out.NextToken != nil {
			out, err = svc.ListContainerInstances(&ecs.ListContainerInstancesInput{Cluster: arn})
			if err != nil {
				return nil, fmt.Errorf("error listing ECS cluster=%s instances: %v", *arn, err)
			}
			for _, arn := range out.ContainerInstanceArns {
				instances = append(instances, *arn)
			}
		}

		clusters = append(clusters, shared.CloudAsset{ID: name, })
	}
	fmt.Println("Found AWS ECS clusters", len(clusters))
	return &shared.CloudDiscoveryResult{ Assets:clusters, Region: c.opt.Region, Type:"ECS" } , nil
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