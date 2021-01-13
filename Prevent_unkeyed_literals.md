
- To prevent unkeyed literals add `_ struct{}` field:

```go
type Point struct {
	X, Y float64
	_    struct{} // to prevent unkeyed literals
}
```

For `Point{X: 1, Y: 1}` everything will be fine, but for `Point{1,1}` you will get a compile error:

```
./file.go:1:11: too few values in Point literal
```

Moreover, 

```go
type Pointer struct {
	X, Y int
	_ struct{} // to prevent unkeyed literals
	// _ byte
	//_ func()
	//_ [0]func()
}
```
