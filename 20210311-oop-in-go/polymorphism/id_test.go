package polymorphism

import "testing"

const (
	fakeid = "412717199109031697"
)

func TestID_Birthday(t *testing.T) {
	var i ID
	i, _ = NewID(fakeid)
	found := i.Birthday()
	wanted := "19910903"
	if found != wanted {
		t.Errorf("unexpected birthday, wanted:%v, found:%v", wanted, found)
	}
}
