
- **Do NOT compare empty struct, per [spec](https://golang.org/ref/spec#Comparison_operators) Pointers to distinct zero-size variables may or may not be equal.**

```go
func f1() {
	var a, b struct{}
	print(&a, "\n", &b, "\n") // Prints same address
	fmt.Println(&a == &b)     // Comparison returns false
}

func f2() {
	var a, b struct{}
	fmt.Printf("%p\n%p\n", &a, &b) // Again, same address
	fmt.Println(&a == &b)          // ...but the comparison returns true
}
```

The reason is due to `escape analysis`, run `go run -gcflags '-m' play.go` to debug it

for more details, please refer to [this](https://stackoverflow.com/a/52421277/3011380)

- empty struct has width of zero

```go
type S struct {
        A struct{}
        B struct{}
}
var s S
fmt.Println(unsafe.Sizeof(s)) // prints 0
```

```go
var x [1000000000]struct{}
fmt.Println(unsafe.Sizeof(x)) // prints 0

var x = make([]struct{}, 100)
var y = x[:50]
fmt.Println(len(y), cap(y)) // prints 50 100

a := make([]struct{}, 10)
b := make([]struct{}, 20)
fmt.Println(&a == &b)       // false, a and b are different slices
fmt.Println(&a[0] == &b[0]) // true, their backing arrays are the same
```

- method receiver

```go
type S struct{}

func (s *S) addr() { fmt.Printf("%p\n", s) }

func main() {
        var a, b S
        a.addr() // 0x1beeb0
        b.addr() // 0x1beeb0
}
```


