package errors

import "fmt"

type ExistError struct {
	Type string
	Id   string
}

func (error ExistError) Error() string {
	return fmt.Sprintf("%s %s exist", error.Type, error.Id)
}

type NotExistError struct {
	Type string
	Id   string
}

func (error NotExistError) Error() string {
	return fmt.Sprintf("%s %s doesn't exist", error.Type, error.Id)
}
