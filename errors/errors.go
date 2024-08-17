package errors

import "fmt"

type RecordNotFoundError struct {
	Entity string
}

func NewRecordNotFoundErr(entity string) *RecordNotFoundError {
	return &RecordNotFoundError{
		Entity: entity,
	}
}

func (e *RecordNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}

type DuplicateDataError struct {
	entity string
	field  string
	value  string
}

func NewDuplicateDataError(entity string, field string, value string) *DuplicateDataError {
	return &DuplicateDataError{
		entity: entity,
		field:  field,
		value:  value,
	}
}

func (err *DuplicateDataError) Error() string {
	return fmt.Sprintf("%s with %s %q already exists", err.entity, err.field, err.value)
}
