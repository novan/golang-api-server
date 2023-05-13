package validator

import (
	"regexp"
	"time"

	"github.com/novan/golang-api-server/util"
	"gopkg.in/go-playground/validator.v9"
)

func ValidateDate(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	_, err := time.Parse(util.TIMEFORMAT_DATE, val)
	return err == nil
}

func ValidateTime(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	rg := regexp.MustCompile("[0-9]{2}")
	if rg.MatchString(val) {
		ts := rg.FindAllString(val, -1)
		if len(ts) == 2 {
			hrs := util.AtoI(ts[0])
			mins := util.AtoI(ts[1])
			return hrs < 24 && mins < 60
		}
	}
	return false

}
