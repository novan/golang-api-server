package errors

import "errors"

func UnauthorizedAccessError() error {
	return NewAppError(errors.New("unauthorized access"), NotAuthorized)
}

func InvalidTokenError() error {
	return NewAppError(errors.New("invalid access token"), BadRequest)
}

func InvalidRefreshTokenError() error {
	return NewAppError(errors.New("invalid refresh token"), BadRequest)
}

func ExpiredRefreshTokenError() error {
	return NewAppError(errors.New("this refresh token is expired"), BadRequest)
}

func NotMatchRefreshTokenError() error {
	return NewAppError(errors.New("this refresh token doesn't belong to the current user"), BadRequest)
}

func InvalidatedRefreshTokenError() error {
	return NewAppError(errors.New("this refresh token has been invalidated"), BadRequest)
}

func UsedRefreshTokenError() error {
	return NewAppError(errors.New("this refresh token has been used"), BadRequest)
}

func UserNotFoundError() error {
	return NewAppError(errors.New("user not found"), BadRequest)
}

func InvalidVerificationCodeError() error {
	return NewAppError(errors.New("invalid verification code"), BadRequest)
}

func ExpiredVerificationCodeError() error {
	return NewAppError(errors.New("verification code was expired"), BadRequest)
}

func AlreadyUsedVerificationCodeError() error {
	return NewAppError(errors.New("verification code was used"), BadRequest)
}

func GeneratePasswordError() error {
	return InternalServerError("Failed generating password")
}

func ParsingTokenError() error {
	return InternalServerError("Failed parsing the token or token might be invalid")
}

func InvalidUserLoginError() error {
	return NewAppError(errors.New("wrong email or password"), BadRequest)
}

func EmailAlreadyExistError() error {
	return NewAppError(errors.New("email sudah terdaftar"), BadRequest)
}
