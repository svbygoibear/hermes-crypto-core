package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TimestampTime struct {
	time.Time
}

func (t *TimestampTime) MarshalJSON() ([]byte, error) {
	bin := make([]byte, 0, len("0000-00-00T00:00:00.00Z"))
	bin = append(bin, fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))...)
	return bin, nil
}

func (t *TimestampTime) UnmarshalJSON(bin []byte) error {
	s := strings.Trim(string(bin), string([]byte{0, '"'}))
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

func (t *TimestampTime) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	if av == nil {
		return nil
	}

	switch av := av.(type) {
	case *types.AttributeValueMemberS:
		parsedTime, err := time.Parse(time.RFC3339, av.Value)
		if err != nil {
			return fmt.Errorf("failed to parse time string: %v", err)
		}
		*t = TimestampTime{Time: parsedTime}
		return nil
	default:
		return fmt.Errorf("unsupported AttributeValue type: %T", av)
	}
}
