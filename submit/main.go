package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	validateRole = validateRoleFunc
	queueURL     = os.Getenv("QUEUE")
)

type entry struct {
	Role  string `json:"role"`
	Email string `json:"email"`
}

type errorBody struct {
	Message string `json:"message"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	awsCfg, err := config.LoadDefaultConfig()
	if err != nil {
		panic(err)
	}

	var entry entry
	json.Unmarshal([]byte(request.Body), &entry)

	if validateRole(awsCfg, entry.Role) == false {
		message := fmt.Sprintf("couldn't assume role: %s", entry.Role)
		return errorResponse(message), nil
	}

	if entry.enqueue(awsCfg) == false {
		return errorResponse("error saving entry"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func validateRoleFunc(awsCfg aws.Config, role string) bool {
	stss := sts.NewFromConfig(awsCfg)
	_, err := stss.AssumeRole(context.Background(), &sts.AssumeRoleInput{
		RoleArn:         &role,
		RoleSessionName: aws.String("validateRole"),
	})

	if err != nil {
		log.Println("assume error ", err)
		return false
	}

	return true
}

func errorResponse(message string) events.APIGatewayProxyResponse {
	errorMessage := &errorBody{
		Message: message,
	}
	body, _ := json.Marshal(errorMessage)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 422,
	}
}

func (e entry) enqueue(awsCfg aws.Config) bool {
	sqss := sqs.NewFromConfig(awsCfg)

	body, errMarshal := json.Marshal(e)

	if errMarshal != nil {
		log.Println("enqueue Marshal ", errMarshal)
		return false
	}

	stringBody := string(body)

	_, err := sqss.SendMessage(context.Background(), &sqs.SendMessageInput{
		MessageBody: &stringBody,
		QueueUrl:    &queueURL,
	})

	if err != nil {
		log.Println("enqueue sendMessage", err)
		return false
	}

	return true
}

func main() {
	lambda.Start(handler)
}
