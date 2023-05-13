package customtype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/novan/golang-api-server/util"
)

const dateTimeLayout = "2006-01-02 15:04:05"

// CommerceTime is the time.Time with JSON marshal and unmarshal capability
type DateTime struct {
	time.Time
}

func ParseDateTime(dt time.Time) DateTime {
	thisDateTime := DateTime{}
	thisDateTime.Time = dt
	return thisDateTime
}

func ParseStringToDateTime(src string) (DateTime, error) {
	t := DateTime{}
	parsed, err := time.Parse(dateTimeLayout, src)
	if err != nil {
		return ParseDateTime(time.Now()), err
	}
	t.Time = parsed
	return t, nil
}

// UnmarshalJSON will unmarshal using 2006-01-02 layout
func (t *DateTime) UnmarshalJSON(b []byte) error {
	src := string(b)
	if src == "" || src == "\"null\"" {
		return nil
	}
	parsed, err := time.Parse(`"`+dateTimeLayout+`"`, src)
	if err == nil {
		t.Time = parsed
	} else {
		err = errors.New(fmt.Sprintf("Invalid datetime format: %s", util.ToString(src)))
	}
	return err
}

func (t *DateTime) UnmarshalParam(src string) error {
	if src == "" {
		return nil
	}
	ts, err := time.Parse(dateTimeLayout, src)
	if err == nil {
		t.Time = ts
	} else {
		err = errors.New(fmt.Sprintf("Invalid datetime format: %s", util.ToString(src)))
	}
	return err
}

// MarshalJSON will marshal using 2006-01-02 layout
func (t *DateTime) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, nil
	}
	s := t.Format(`"` + dateTimeLayout + `"`)
	return []byte(s), nil
}

func (t *DateTime) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return t.Time.Format(dateTimeLayout), nil
}

func (t *DateTime) Scan(value interface{}) error {
	var err error
	if reflect.TypeOf(value).Kind() == reflect.String {
		t.Time, err = time.Parse(dateTimeLayout, value.(string))
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		timeBytes := value.([]byte)
		t.Time, err = time.Parse(dateTimeLayout, string(timeBytes))
	} else if _, ok := value.(time.Time); ok {
		t.Time = value.(time.Time)
	} else {
		err = errors.New(fmt.Sprintf("Invalid datetime format: %s", util.ToString(value)))
	}

	return err
}

func (t *DateTime) String() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Format(dateTimeLayout)
}
