package errors

import "errors"

var (
	UserAlreadyExist = errors.New("user already exist")
	UserNotFound     = errors.New("user not found")
	InvalidUUID      = errors.New("invalid uuid")
	InvalidOrderBy   = errors.New("invalid order by")
)
