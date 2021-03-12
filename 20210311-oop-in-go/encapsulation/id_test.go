package encapsulation

import (
	"testing"
)

const (
	fakeid = "412717199109031697"
)

func Test_Birthday(t *testing.T) {
	found := Birthday(fakeid)
	wanted := "19910903"
	if found != wanted {
		t.Errorf("unexpected birthday, wanted:%v, found:%v", wanted, found)
	}
}

func TestID_Birthday(t *testing.T) {
	found := ID(fakeid).Birthday()
	wanted := "19910903"
	if found != wanted {
		t.Errorf("unexpected birthday, wanted:%v, found:%v", wanted, found)
	}
}
