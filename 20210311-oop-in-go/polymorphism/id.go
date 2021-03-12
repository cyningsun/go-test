package polymorphism

import (
	"errors"
	"fmt"
)

type ID interface {
	Birthday() string
}

type id string

func NewID(i string) (ID, error) {
	if len(i) != 18 {
		return nil, errors.New(fmt.Sprintf("error id length:%v", len(i)))
	}
	return id(i), nil
}

func (i id) Birthday() string {
	return string(i[6:14])
}
