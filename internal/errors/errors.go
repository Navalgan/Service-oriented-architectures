package errors

import "errors"

var (
	UserAlreadyExist = errors.New("user already exist")
	UserNotFound     = errors.New("user not found")
)
