package ec2_events_dashboard

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ec2Regions []string

func init() {
	ec2Regions = []string{
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"eu-west-1",
		"eu-central-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"sa-east-1",
	}
}

// NewClientsFromCreds returns a slice of *ec2.EC2 for each combination of strings
// in creds and regions in AWS. Each string in creds is expected to look
// like: "AccessKeyID:SecretAccessKey". If any are not, an error is returned
func NewClientsFromCreds(creds []string) ([]*ec2.EC2, error) {
	ec2Clients := []*ec2.EC2{}

	for _, cred := range creds {
		for _, region := range ec2Regions {
			// TODO: use a regex to verify input key format is valid?
			// http://rubular.com/r/HrybQ6uiQf
			split := strings.Split(cred, ":")
			ec2Clients = append(ec2Clients, newEc2Client(split[0], split[1], region))
		}
	}

	return ec2Clients, nil
}

func newEc2Client(accessKeyID string, secretAccessKey string, region string) *ec2.EC2 {
	config := aws.NewConfig().
		WithRegion(region).
		WithCredentials(credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""))

	return ec2.New(session.New(config))
}
