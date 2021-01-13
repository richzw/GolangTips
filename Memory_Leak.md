- Caused by Substrings

```go
var s0 string // a package-level variable

// A demo purpose function.
func f(s1 string) {
	s0 = s1[:50]
	// Now, s0 shares the same underlying memory block
	// with s1. Although s1 is not alive now, but s0
	// is still alive, so the memory block they share
	// couldn't be collected, though there are only 50
	// bytes used in the block and all other bytes in
	// the block become unavailable.
}

func demo() {
	s := createStringWithLengthOnHeap(1 << 20) // 1M bytes
	f(s)
}
```

solve it 
```go
import "strings"

func f(s1 string) {
	var b strings.Builder
	b.Grow(50)
	b.WriteString(s1[:50])
	s0 = b.String()
}
```

- Caused by Subslices

```go
var s0 []int

func g(s1 []int) {
	// Assume the length of s1 is much larger than 30.
	s0 = s1[len(s1)-30:]
}
```

solve it

```go
func g(s1 []int) {
	s0 = append([]int(nil), s1[len(s1)-30:]...)
	// Now, the memory block hosting the elements
	// of s1 can be collected if no other values
	// are referencing the memory block.
}
```

- Caused by Not Resetting Pointers in Lost Slice Elements

```go
func h() []*int {
	s := []*int{new(int), new(int), new(int), new(int)}
	// do something with s ...

	return s[1:3:3]
}
```

solve it

```go
func h() []*int {
	s := []*int{new(int), new(int), new(int), new(int)}
	// do something with s ...

	// Reset pointer values.
	s[0], s[len(s)-1] = nil, nil
	return s[1:3:3]
}
```

- Caused by Hanging Goroutines

- Caused by Not Stopping `time.Ticker` Values Which Are Not Used Any More

- Caused by Using Finalizers Improperly

```go
func memoryLeaking() {
	type T struct {
		v [1<<20]int
		t *T
	}

	var finalizer = func(t *T) {
		 fmt.Println("finalizer called")
	}

	var x, y T

	// The SetFinalizer call makes x escape to heap.
	runtime.SetFinalizer(&x, finalizer)

	// The following line forms a cyclic reference
	// group with two members, x and y.
	// This causes x and y are not collectable.
	x.t, y.t = &y, &x // y also escapes to heap.
}
```





