package service

type AccountsError string

func (e AccountsError) Error() string {
	return (string)(e)
}

const (
	InvalidEmail    AccountsError = "invalid email"
	AlreadyExist    AccountsError = "already exists"
	NoNameProvided  AccountsError = "no name provided"
	NoTokenProvided AccountsError = "no token provided"
	NoUser          AccountsError = "no such user"
	NoContent       AccountsError = "no content"
	InvalidName     AccountsError = "invalid name"
)
