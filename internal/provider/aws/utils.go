package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"time"
)

// AWSSession creates a session for AWS services handling
func CreateAWSSession(opt *AWSOptions) (*session.Session, error) {
	var creds *credentials.Credentials
	cfg := defaults.Config()
	if opt.UseAWSRole {
		// Get Role credentials from EC2 metadata endpoint. EC2RoleProvider retrieves credentials from the EC2 service, and keeps track if
		// those credentials are expired.
		endpoint, err := endpoints.DefaultResolver().EndpointFor(ec2metadata.ServiceName, opt.Region)
		if err != nil {
			return nil, err
		}
		creds = credentials.NewCredentials(&ec2rolecreds.EC2RoleProvider{
			Client:       ec2metadata.NewClient(*cfg, defaults.Handlers(), endpoint.URL, endpoint.SigningRegion),
			ExpiryWindow: 5 * time.Minute,
		})
	} else {
		if opt.SecretAccessKey == "" {
			return nil, fmt.Errorf("missing secret key in AWS settings")
		}

		creds = credentials.NewStaticCredentials(opt.AccessKeyID, opt.SecretAccessKey, "")
	}

	return session.NewSession(cfg.WithCredentials(creds).WithRegion(opt.Region))
}

func Discover(username, password string, emitFn func(result shared.CloudDiscoveryResult)) {
	var discoverers []shared.Discoverer
	opt := AWSOptions{
		AccessKeyID:     username,
		SecretAccessKey: password,
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
	}
	for _, region := range lambdaRegions {
		opt.Region = region
		discoverers = append(discoverers, NewECRClient(opt))
	}

	for _, discoverer := range discoverers {
		result, err := discoverer.Discover()
		if err != nil {
			log.Debugf(err.Error())
		} else if len(result.Assets) > 0 {
			emitFn(*result)
		}
	}
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
