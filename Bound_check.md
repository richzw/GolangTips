
Go is a memory safe language. In array/slice element indexing and subslice operations, Go runtime will check whether or not the involved indexes are out of range. 
If an index is out of range, a panic will be produced to prevent the invalid index from doing harm. This is called **bounds check**. Bounds checks make our code run safely, on the other hand, they also make our code run a little slower.

Go compiler has used a new compiler backend, which based on SSA (static single-assignment form). SSA helps Go compilers effectively use optimizations like BCE (bounds check elimination) and CSE (common subexpression elimination)

```go
func fb2(s []int, index int) {
	_ = s[index:] // line 7: bounds check
	_ = s[:index] // line 8: bounds check eliminated!
}

func fc2() {
	s := []int{0, 1, 2, 3, 4, 5, 6}
	s = s[:4]
	index := rand.Intn(7)
	_ = s[index:] // line 15 bounds check
	_ = s[:index] // line 16: bounds check eliminated!
}
```

run `go build -gcflags="-d=ssa/check_bce/debug=1"`




Ref:

[1](https://docs.google.com/document/d/1vdAEAjYdzjnPA9WDOQ1e4e05cYVMpqSxJYZT33Cqw2g/edit)

