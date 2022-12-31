package optimize

import "testing"

func TestGetLimitSize(t *testing.T) {
	testVars := [][5]int{
		{100, 200, 100 /* -> */, 50, 100},
		{1000, 2000, 1000 /* -> */, 500, 1000},
		{200, 100, 100 /* -> */, 100, 50},
		{100, 100, 100 /* -> */, 100, 100},
	}
	for _, vars := range testVars {
		w := vars[0]
		h := vars[1]
		l := vars[2]
		newW := vars[3]
		newH := vars[4]
		resW, resH := getLimitSize(w, h, l)
		if resW != newW {
			t.Error("width mismatch:", resW, newW)
		}
		if resH != newH {
			t.Error("height mismatch:", resH, newH)
		}
	}
}
