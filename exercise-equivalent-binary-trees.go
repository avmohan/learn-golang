package main

import ("fmt"
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
// The channel quit is used to listen for a signal
// to quit
func WalkSub(t *tree.Tree, ch, quit chan int) {
	if t == nil {
		return
	}
	WalkSub(t.Left, ch, quit)
	select {
	case ch <- t.Value:
		// pass (value pushed into ch)
	case <-quit:
		// this unblocks when channel quit is closed
		return
	}
	WalkSub(t.Right, ch, quit)

}

func Walk(t *tree.Tree, ch, quit chan int) {
	// Need to separate this because close should only be done once
	// at the root level
	WalkSub(t, ch, quit)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	c2 := make(chan int)
	quit := make(chan int)
	defer close(quit)

	go Walk(t1, c1, quit)
	go Walk(t2, c2, quit)
	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2
		if !ok1 || !ok2 {
			return !ok1 && !ok2
		}
		if v1 != v2 {
			return false
		}
	}
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
