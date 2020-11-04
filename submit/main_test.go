package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("Invalid Role", func(t *testing.T) {
		validateRole = func(cfg aws.Config, role string) bool {
			return false
		}
		defer func() {
			validateRole = validateRoleFunc
		}()

		expected := events.APIGatewayProxyResponse{
			Body:       "{\"message\":\"couldn't assume role: mocked-as-invalid\"}",
			StatusCode: 422,
		}

		entry := entry{
			Role:  "mocked-as-invalid",
			Email: "test@test.com",
		}

		body, _ := json.Marshal(entry)

		event := events.APIGatewayProxyRequest{
			Body: string(body),
		}
		response, _ := handler(event)

		assert.Equal(t, response, expected, "should be error response")
	})

	t.Run("Valid Role", func(t *testing.T) {
		validateRole = func(cfg aws.Config, role string) bool {
			return true
		}
		defer func() {
			validateRole = validateRoleFunc
		}()

		expected := events.APIGatewayProxyResponse{
			StatusCode: 200,
		}

		entry := entry{
			Role:  "mocked-as-valid",
			Email: "test@test.com",
		}

		body, _ := json.Marshal(entry)

		event := events.APIGatewayProxyRequest{
			Body: string(body),
		}
		response, _ := handler(event)

		assert.Equal(t, response, expected, "should be error response")
	})
}
