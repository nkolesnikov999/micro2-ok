package model

import "errors"

var (
	ErrSessionNotFound    = errors.New("session not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserGetFailed      = errors.New("user get failed")
	ErrUserCreateFailed   = errors.New("user create failed")
	ErrInvalidLogin       = errors.New("invalid login")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrPasswordHashFailed = errors.New("password hash failed")
)
