package model

import "errors"

var (
	ErrUnexpectedTokenSignMethod = errors.New("unexpected token sign method")
	ErrInvalidToken              = errors.New("invalid token")
	ErrCardAlreadyExists         = errors.New("card with this number already exists")
	ErrUserAlreadyExists         = errors.New("user with this login already exists")
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidPassword           = errors.New("invalid password")
)
