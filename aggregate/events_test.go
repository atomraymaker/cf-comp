package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalStreamImage(t *testing.T) {
	t.Run("unmarshalStreamImage", func(t *testing.T) {
		event := event()
		image := event.Records[0].Change.NewImage
		var entry entry
		unmarshalStreamImage(image, &entry)
		assert.Equal(t, "testRole", entry.Role)
		assert.Equal(t, "testRole", entry.Role)
		assert.Equal(t, 100, entry.ValidResources)
		assert.Equal(t, 20, entry.CfResources)
	})
}

func event() events.DynamoDBEvent {
	return events.DynamoDBEvent{
		Records: []events.DynamoDBEventRecord{
			{
				Change: events.DynamoDBStreamRecord{
					NewImage: map[string]events.DynamoDBAttributeValue{
						"Role":           events.NewStringAttribute("testRole"),
						"Email":          events.NewStringAttribute("testEmail"),
						"ValidResources": events.NewNumberAttribute("100"),
						"CfResources":    events.NewNumberAttribute("20"),
					},
				},
			},
		},
	}
}
