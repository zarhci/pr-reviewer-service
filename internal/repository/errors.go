package repository

import "errors"

var (
	// common
	ErrNotFound = errors.New("NOT_FOUND")
	ErrExists   = errors.New("ALREADY_EXISTS")

	// user
	ErrUserNotFound = ErrNotFound

	// team
	ErrTeamNotFound = ErrNotFound
	ErrTeamExists   = ErrExists

	// pull request
	ErrPRNotFound = ErrNotFound
	ErrPRExists   = ErrExists
)
