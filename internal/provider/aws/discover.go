package aws

import (
	"github.com/twistlock/cloud-discovery/internal/shared"
	"log"
)

func Discover(username, password string) shared.CloudDiscoveryResults {
	var discoverers []shared.Discoverer
	opt := AWSOptions{
		AccessKeyID:username,
		SecretAccessKey:password,
	}

	for _, region := range eksRegions {
		opt.Region = region
		discoverers = append(discoverers, NewEKSClient(opt))
	}
	for _, region := range ecsRegions {
		opt.Region = region
		discoverers = append(discoverers, NewECSClient(opt))
	}
	for _, region := range lambdaRegions {
		opt.Region = region
		discoverers = append(discoverers, NewLambdaClient(opt))
		discoverers = append(discoverers, NewECRClient(opt))
	}


	var result shared.CloudDiscoveryResults
	for _, discoverer := range discoverers {
		info, err := discoverer.Discover()
		if err != nil {
			log.Println(err.Error())
		} else if len(info.Assets) == 0 {
			log.Println("No results", info.Type, info.Region)
		} else {
			result.Results = append(result.Results, *info)
		}
	}
	return result
}


// eksRegions - Amazon known EKS regions list
// https://docs.aws.amazon.com/general/latest/gr/rande.html#eks_region
var eksRegions = []string{
	"us-east-1", // N. Virginia
	"us-west-2", // Oregon
}

var ecsRegions = []string{
	"us-east-2",      // US East (Ohio)
	"us-east-1",      // US East (N. Virginia)
	"us-west-1",      // US West (N. California)
	"us-west-2",      // US West (Oregon)
	"ap-northeast-1", // Asia Pacific (Tokyo)
	"ap-northeast-2", // Asia Pacific (Seoul)
	"ap-south-1",     // Asia Pacific (Mumbai)
	"ap-southeast-1", // Asia Pacific (Singapore)
	"ap-southeast-2", // Asia Pacific (Sydney)
	"ca-central-1",   // Canada (Central)
	"cn-north-1",     // China (Beijing)
	"cn-northwest-1", // China (Ningxia)
	"eu-central-1",   // EU (Frankfurt)
	"eu-west-1",      // EU (Ireland)
	"eu-west-2",      // EU (London)
	"eu-west-3",      // EU (Paris)
	"sa-east-1",      // South America (São Paulo)
}


// awsRegions - Amazon known Lambda regions list
// https://docs.aws.amazon.com/general/latest/gr/rande.html#lambda_region
var lambdaRegions = []string{
	"us-east-2",
	"us-east-1",
	"us-west-1",
	"us-west-2",
	"ap-northeast-2",
	"ap-south-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-northeast-1",
	"ca-central-1",
	"eu-central-1",
	"eu-west-1",
	"eu-west-2",
	"sa-east-1",
}