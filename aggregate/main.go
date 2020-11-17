package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	tableName = os.Getenv("TABLE_NAME")
	sess, _   = session.NewSessionWithOptions(session.Options{})
)

type aggregate struct {
	Type   string
	Values map[string]int
}

type entry struct {
	Role           string `json:"role"`
	Email          string `json:"email"`
	ValidResources int    `json:"validResources"`
	CfResources    int    `json:"cfResources"`
}

func handler(ddbEvents events.DynamoDBEvent) error {
	for _, record := range ddbEvents.Records {
		entry := entry{}
		unmarshalStreamImage(record.Change.NewImage, &entry)

		fmt.Println(entry)

		// topPercentage(entry)
	}
	return nil
}

func topPercentage(entry entry) error {
	top10, err := getAgg("topPercentage")

	if err != nil {
		return err
	}

	score := entry.CfResources / entry.ValidResources
	newTop10 := copyMap(top10)

	if val, ok := newTop10[entry.Email]; ok {
		if score > val {
			newTop10[entry.Email] = score
			err = updateAgg("topPercentage", top10, newTop10)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	newTop10[entry.Email] = score

	newTop10 = trimMap(newTop10, 10)

	updateAgg("topPercentage", top10, newTop10)

	return nil
}

// update aggregate dynamodb
func updateAgg(name string, orig map[string]int, new map[string]int) error {
	fmt.Println(new)

	return nil
}

// get aggregate dynamodb record
func getAgg(name string) (map[string]int, error) {
	ddb := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String("AGGS"),
			},
			"sk": {
				S: aws.String(name),
			},
		},
	}

	result, err := ddb.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return map[string]int{}, err
			default:
				panic(aerr.Error())
			}
		}
	}

	if result.Item == nil {
		return map[string]int{}, nil
	}

	agg := aggregate{}
	err1 := dynamodbattribute.UnmarshalMap(result.Item, &agg)

	if err1 != nil {
		return map[string]int{}, err1
	}

	return agg.Values, nil
}

func main() {
	lambda.Start(handler)
}
