package util

import (
	"encoding/json"
	"time"
)

const DateFormat = "2006-01-02"

type Date time.Time

func (t *Date) UnmarshalJSON(data []byte) error {
	var value int64
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	*t = Date(time.Unix(value, 0))
	return err
}

func (t Date) String() string {
	return time.Time(t).Format(DateFormat)
}
