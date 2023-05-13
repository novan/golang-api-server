package customtype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/novan/golang-api-server/util"
)

const dateLayout = "2006-01-02"

// Date is the time.Time with JSON marshal and unmarshal capability
type Date struct {
	time.Time
}

func ParseDate(dt time.Time) Date {
	thisDate := Date{}
	thisDate.Time = dt
	return thisDate
}

func ParseStringToDate(src string) (Date, error) {
	t := Date{}
	parsed, err := time.Parse(dateLayout, src)
	if err != nil {
		return ParseDate(time.Now()), err
	}
	t.Time = parsed
	return t, nil
}

// UnmarshalJSON will unmarshal using 2006-01-02 layout
func (t *Date) UnmarshalJSON(b []byte) error {
	src := string(b)
	if src == "" || src == "\"null\"" {
		return nil
	}
	parsed, err := time.Parse(`"`+dateLayout+`"`, src)
	if err == nil {
		t.Time = parsed
	} else {
		err = errors.New(fmt.Sprintf("Invalid date format: %s", util.ToString(src)))
	}
	return err
}

func (t *Date) UnmarshalParam(src string) error {
	if src == "" {
		return nil
	}
	ts, err := time.Parse(dateLayout, src)
	if err == nil {
		t.Time = ts
	} else {
		err = errors.New(fmt.Sprintf("Invalid date format: %s", util.ToString(src)))
	}
	return err
}

// MarshalJSON will marshal using 2006-01-02 layout
func (t *Date) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, nil
	}
	s := t.Format(`"` + dateLayout + `"`)
	return []byte(s), nil
}

func (t Date) Value() (driver.Value, error) {
	return t.Time.Format(dateLayout), nil
}

func (t *Date) Scan(value interface{}) error {
	var err error
	if reflect.TypeOf(value).Kind() == reflect.String {
		t.Time, err = time.Parse(dateLayout, value.(string))
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		dateBytes := value.([]byte)
		t.Time, err = time.Parse(dateLayout, string(dateBytes))
	} else if _, ok := value.(time.Time); ok {
		t.Time = value.(time.Time)
	} else {
		err = errors.New(fmt.Sprintf("Invalid date format: %s", util.ToString(value)))
	}

	return err
}

func (t Date) String() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Format(dateLayout)
}
