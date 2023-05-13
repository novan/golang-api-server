package util

// Context variables
const (
	CONTEXT_SESSION     = "ctx_sess"
	CONTEXT_TOKEN       = "ctx_token"
	CONTEXT_ROUTE_NAME  = "ROUTE-NAME"
	CONTEXT_TARGET_HOST = "TARGET-HOST"
	CONTEXT_TARGET_PATH = "TARGET-PATH"
)

// User Type
type UserType string

const (
	USERTYPE_USER  UserType = "USER"
	USERTYPE_ADMIN UserType = "ADMIN"
)

func (u UserType) ToString() string {
	switch u {
	case USERTYPE_ADMIN:
		return "Admin"
	default:
		return "User"
	}
}

// Time format
const (
	TIMEFORMAT_DATE             = "2006-01-02"
	TIMEFORMAT_DATE_SHORT       = "02-Jan-2006"
	TIMEFORMAT_DATE_FORMAL      = "02 January 2006"
	TIMEFORMAT_DATETIME         = "2006-01-02 15:04:05"
	TIMEFORMAT_DATETIME_CLASSIC = "1/2/2006 03:04:05 PM"
	TIMEFORMAT_DATETIME_DOTNET  = "2006-01-02T15:04:05Z"
	TIMEFORMAT_TIME             = "15:04"
)
