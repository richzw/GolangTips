
**Do NOT compare empty struct, per [spec](https://golang.org/ref/spec#Comparison_operators) Pointers to distinct zero-size variables may or may not be equal.**

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

