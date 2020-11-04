package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/trek10inc/awsets"
	"github.com/trek10inc/awsets/resource"
)

var (
	tableName = os.Getenv("TABLE_NAME")
)

type entry struct {
	Role  string `json:"role"`
	Email string `json:"email"`
}

func handler(sqsEvent events.SQSEvent) error {
	records := len(sqsEvent.Records)
	if records != 1 {
		return fmt.Errorf("can only process one record per invocation, number: %d", records)
	}

	record := sqsEvent.Records[0]
	var entry entry
	json.Unmarshal([]byte(record.Body), &entry)

	localConfig, err := config.LoadDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed config.LoadDefaultConfig: %w", err)
	}

	clientConfig, clientAccount, err1 := getClientSession(localConfig, entry.Role)
	if err1 != nil {
		return fmt.Errorf("failed getClientSession: %w", err1)
	}

	resources, err2 := runAwsets(clientConfig, entry)
	if err2 != nil {
		return fmt.Errorf("failed runAwsets: %w", err2)
	}

	inCF, total := calcScores(resources)
	saveResults(localConfig, entry, clientAccount, inCF, total)

	return nil
}

func saveResults(localConfig aws.Config, entry entry, clientAccount string, inCF int, total int) {
	ddb := dynamodb.NewFromConfig(localConfig)

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*ddbTypes.AttributeValue{
			"pk": {
				S: aws.String(entry.Email),
			},
			"sk": {
				S: aws.String(clientAccount),
			},
			"role": {
				S: aws.String(entry.Role),
			},
			"validResources": {
				N: aws.String(fmt.Sprint(total)),
			},
			"cfResources": {
				N: aws.String(fmt.Sprint(inCF)),
			},
		},
	}

	ddb.PutItem(context.Background(), input)
}

func calcScores(resources *resource.Group) (int, int) {
	inCF := 0
	total := 0
	for _, resource := range resources.Resources {
		if resource.Type != "cloudformation/stack" {
			total++
			for name := range resource.Tags {
				if name == "aws:cloudformation:stack-id" {
					inCF++
				}
			}
		}
	}
	return inCF, total
}

func runAwsets(clientConfig aws.Config, entry entry) (*resource.Group, error) {
	regions, err := awsets.Regions(clientConfig)
	if err != nil {
		return &resource.Group{}, fmt.Errorf("failed to list regions: %w", err)
	}
	listers := awsets.Listers(nil, nil)
	if err != nil {
		return &resource.Group{}, fmt.Errorf("failed to create awsets ctx: %w", err)
	}
	return awsets.List(clientConfig, regions, listers, nil)
}

func getClientSession(localConfig aws.Config, role string) (aws.Config, string, error) {
	clientConfig := localConfig.Copy()
	stss := sts.NewFromConfig(localConfig)

	assumeOutput, err1 := stss.AssumeRole(context.Background(), &sts.AssumeRoleInput{
		RoleArn:         &role,
		RoleSessionName: aws.String("validateRole"),
	})

	if err1 != nil {
		return aws.Config{}, "", fmt.Errorf("failed stss.AssumeRole: %w", err1)
	}

	c := assumeOutput.Credentials
	creds := credentials.NewStaticCredentialsProvider(*c.AccessKeyId, *c.SecretAccessKey, *c.SessionToken)
	clientConfig.Credentials = creds

	clientAccount, err2 := stss.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})

	if err2 != nil {
		return aws.Config{}, "", fmt.Errorf("failed stss.AssumeRole: %w", err2)
	}

	return clientConfig, *clientAccount.Account, nil
}

func main() {
	lambda.Start(handler)
}
