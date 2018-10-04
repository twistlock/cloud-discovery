package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"time"
)

// AWSOptions are options for getting credentials for AWS services
type AWSOptions struct {
	AccessKeyID     string // AccessKeyID is the access key ID to aws services
	SecretAccessKey string // SecretAccessKey is the secret access key to aws services
	Region          string // Region is the region in query
	UseAWSRole      bool   // UseAWSRole is a flag indicates if local IAM role should be used for authentication, than username and password would be ignored
}

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