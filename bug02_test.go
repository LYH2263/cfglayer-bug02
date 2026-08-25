package cfglayer_test

import (
	"testing"
)

func TestBug02_ListLayersIsolate(t *testing.T) {
	m := newMerger(t)
	push(t, m, "alpha", map[string]string{"x": "1"})
	ls := m.ListLayers()
	if len(ls) == 0 {
		t.Fatal("empty")
	}
	ls[0].ID = "mutated"
	ls2 := m.ListLayers()
	if ls2[0].ID == "mutated" {
		t.Fatal("list alias")
	}
}
