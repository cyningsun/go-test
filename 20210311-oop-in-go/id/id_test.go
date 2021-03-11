package id

import (
	"testing"
)

const (
	id = "412717199109031697"
)

func Test_Birthday(t *testing.T) {
	found := Birthday(id)
	wanted := "19910903"
	if found != wanted {
		t.Errorf("unexpected birthday, wanted:%v, found:%v", wanted, found)
	}
}

func TestID_Birthday(t *testing.T) {
	found := ID(id).Birthday()
	wanted := "19910903"
	if found != wanted {
		t.Errorf("unexpected birthday, wanted:%v, found:%v", wanted, found)
	}
}
