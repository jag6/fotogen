package models

import "errors"

var (
	ErrNotFound   = errors.New("models: resource could not be found")
	ErrEmailTaken = errors.New("models: email address already in use")
	//ErrUsernameTaken = errors.New("models: username already in use")
)
