
Key Points
========

- Optimize for correctness not performance
- Go only puts function variables on the stack if it can prove a variable **not used** after the function returns
- Sharing down typically stay on the stack
- Sharing up typically escapes to the heap
- Does not Guess, use the tooling

Check memory - escape analysis
========

`go build -gcflags="-m -l"`

When are values constructed on the heap
=========

- When a value could possibly be referenced after the function that constucted the value returns
- When the compiler determines a value is too large to fit on the stack
- When the compiler does not know the size of a value at compile time

Some commonly allocated values
===========

- values shared with **Pointers**
- values stored in **Interface** variables
- **Func literal** variables
  - Variables captured by a Closure
- Backing data for **Maps, Channels, Slices, and Strings**
  - Strings are effectively a special []byte

Explais `io.Reader`

```go
type Reader interface {
  Read(p []byte) (n int, err error)
}
```

Instead of 


```go
type Reader interface {
  Read(n int) (p []byte, err error)
}
```



