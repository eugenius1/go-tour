package main

import (
	"golang.org/x/tour/tree"
	"testing"
)

func TestWalk(t *testing.T) {
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	expected := 0
	for got := range ch {
		expected++
		if expected > 10 {
			t.Fatalf("chan sent more than the expected number of values, got %d", got)
		}
		if got != expected {
			t.Errorf("Got %d, expected %d", got, expected)
		}
	}
	if expected != 10 {
		t.Errorf("Got %d values, expected 10", expected)
	}
}

func TestSame(t *testing.T) {
	tree1 := tree.New(1)
	tree2 := tree.New(2)
	treeSingle := &tree.Tree{Left: nil, Value: 42, Right: nil}
	cases := []struct {
		t1, t2   *tree.Tree
		expected bool
	}{
		{tree1, tree1, true},
		{tree1, tree2, false},
		{nil, tree1, false},
		{tree2, nil, false},
		{tree2, treeSingle, false},
	}

	for _, c := range cases {
		got := Same(c.t1, c.t2)
		if got != c.expected {
			t.Errorf("\nt1=%v\nt2=%v\n...returned %t, expected %t", c.t1, c.t2, got, c.expected)
		}
	}
}
