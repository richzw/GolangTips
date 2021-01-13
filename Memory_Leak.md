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





