package models

import (
	"fmt"
	"strings"
	"time"
)

type TimestampTime struct {
	time.Time
}

func (t *TimestampTime) MarshalJSON() ([]byte, error) {
	bin := make([]byte, 0, len("2019-10-12T07:20:50.52Z"))
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
