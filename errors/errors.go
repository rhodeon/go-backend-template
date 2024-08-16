package errors

import "fmt"

type DuplicateDataError struct {
	Entity string
	Field  string
	Value  string
}

func NewDuplicateDataError(entity string, field string, value string) DuplicateDataError {
	return DuplicateDataError{
		Entity: entity,
		Field:  field,
		Value:  value,
	}
}

func (err DuplicateDataError) Error() string {
	return fmt.Sprintf("%s with %s %q already exists", err.Entity, err.Field, err.Value)
}
