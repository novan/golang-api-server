package customtype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/novan/golang-api-server/util"
)

const dbTimeLayout = "15:04:00"
const timeLayout = "15:04"

// Date is the time.Time with JSON marshal and unmarshal capability
type Time struct {
	time.Time
}

func ParseTime(dt time.Time) Time {
	thisTime := Time{}
	thisTime.Time = dt
	return thisTime
}

func ParseStringToTime(src string) (Time, error) {
	t := Time{}
	parsed, err := time.Parse(timeLayout, src)
	if err != nil {
		return ParseTime(time.Now()), err
	}
	t.Time = parsed
	return t, nil
}

// UnmarshalJSON will unmarshal using the layout
func (t *Time) UnmarshalJSON(b []byte) error {
	src := string(b)
	if src == "" || src == "\"null\"" {
		return nil
	}
	parsed, err := time.Parse(`"`+timeLayout+`"`, src)
	if err == nil {
		t.Time = parsed
	} else {
		err = errors.New(fmt.Sprintf("Invalid date format: %s", util.ToString(src)))
	}
	return err
}

func (t *Time) UnmarshalParam(src string) error {
	if src == "" {
		return nil
	}
	ts, err := time.Parse(timeLayout, src)
	if err == nil {
		t.Time = ts
	} else {
		err = errors.New(fmt.Sprintf("Invalid date format: %s", util.ToString(src)))
	}
	return err
}

// MarshalJSON will marshal using 2006-01-02 layout
func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, nil
	}
	s := t.Format(`"` + timeLayout + `"`)
	return []byte(s), nil
}

func (t Time) Value() (driver.Value, error) {
	return t.Time.Format(timeLayout), nil
}

func (t *Time) Scan(value interface{}) error {
	var err error
	if reflect.TypeOf(value).Kind() == reflect.String {
		t.Time, err = time.Parse(dbTimeLayout, value.(string))
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		dateBytes := value.([]byte)
		t.Time, err = time.Parse(dbTimeLayout, string(dateBytes))
	} else if _, ok := value.(time.Time); ok {
		t.Time = value.(time.Time)
	} else {
		err = errors.New(fmt.Sprintf("Invalid date format: %s", util.ToString(value)))
	}

	return err
}

func (t Time) String() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func (t Time) DbString() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Format(dbTimeLayout)
}

func (t *Time) Hour() int {
	str := t.String()
	return util.AtoI(str[0:2])
}
