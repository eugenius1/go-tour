package main

import (
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	walkRecursively(t, ch)
}

func walkRecursively(t *tree.Tree, ch chan int) {
	if t != nil {
		walkRecursively(t.Left, ch)
		ch <- t.Value
		walkRecursively(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		// no more values in both trees
		if !ok1 && !ok2 {
			return true
		}

		// no more values in one tree or values don't match
		if ok1 != ok2 || v1 != v2 {
			return false
		}
	}
}
