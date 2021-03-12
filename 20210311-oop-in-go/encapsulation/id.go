package encapsulation

import (
	"errors"
	"fmt"
)

func Birthday(id string) string {
	return id[6:14]
}

type ID string

func NewID(id string) (ID, error) {
	if len(id) != 18 {
		return "", errors.New(fmt.Sprintf("error id length:%v", len(id)))
	}
	return ID(id), nil
}

func (i ID) Birthday() string {
	return string(i[6:14])
}
