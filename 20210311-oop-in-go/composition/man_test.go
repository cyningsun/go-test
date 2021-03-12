package composition

import "testing"

const (
	peterid = "412717199109031697"
	samid   = "312717199109036148"
)

func TestMan_Birthday(t *testing.T) {
	peter, _ := NewMan("peter", peterid)
	sam, _ := NewMan("sam", samid)
	mans := []Man{peter, sam}
	for _, man := range mans {
		found := man.Birthday()
		wanted := "19910903"
		if found != wanted {
			t.Errorf("unexpected birthday, wanted:%v, found:%v", wanted, found)
		}
	}

}
