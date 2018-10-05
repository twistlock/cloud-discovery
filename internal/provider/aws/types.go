package aws

// AWSOptions are options for getting credentials for AWS services
type AWSOptions struct {
	AccessKeyID     string // AccessKeyID is the access key ID to aws services
	SecretAccessKey string // SecretAccessKey is the secret access key to aws services
	Region          string // Region is the region in query
	UseAWSRole      bool   // UseAWSRole is a flag indicates if local IAM role should be used for authentication, than username and password would be ignored
}
