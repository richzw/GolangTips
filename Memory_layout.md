
Structs provide a more flexible way of defining composite types, whose width is the sum of the width of the constituent types, plus padding

```go
type S struct {
        a uint16
        b uint32
}
var s S
fmt.Println(unsafe.Sizeof(s)) // prints 8, not 6
```

It’s not true that “a value must be aligned in memory to a multiple of its width.” Each type has another property, its alignment. 
Alignments are always powers of two. The alignment of a basic type is usually equal to its width, but the alignment of a struct is the maximum alignment of any field, and the alignment of an array is the alignment of the array element. The maximum alignment of any value is therefore the maximum alignment of any basic type. Even on 32-bit systems this is often 8 bytes, because atomic operations on 64-bit values typically require 64-bit alignment.

To be concrete, a struct containing 3 int32 fields has alignment 4 but width 12.

It is true that a value’s width is always a multiple of its alignment. One implication is that there is no padding between array elements.

- Type Alignment Guarantees in Go

In fact, each type has two alignment guarantees, one is for when it is used as field types of other (struct) types, the other is for other cases (when it is used for a variable declaration, array element type, etc). We call the former one the field alignment guarantee of that type, and call the latter one the general alignment guarantee of that type.

For a type `T`, we can call `unsafe.Alignof(t)` to get its general alignment guarantee, where t is a non-field value of type `T`, and call `unsafe.Alignof(x.t)` to get its field alignment guarantee, where `x` is a struct value and t is a field value of type `T`.

Calls to the functions in the `unsafe` standard code packages are **always evaluated at compile time**.

```
type                    size in bytes
------                  ------
byte, uint8, int8       1
uint16, int16           2
uint32, int32, float32  4
uint64, int64           8
float64, complex64      8
complex128              16
uint, int               implementation-specific,
                        generally 4 on 32-bit
                        architectures, and 8 on
                        64-bit architectures.
uintptr                 implementation-specific,
                        large enough to store
                        the uninterpreted bits
                        of a pointer value.
```

```
type T1 struct {
	a int8

	// On 64-bit architectures, to make field b
	// 8-byte aligned, 7 bytes need to be padded
	// here. On 32-bit architectures, to make
	// field b 4-byte aligned, 3 bytes need to be
	// padded here.

	b int64
	c int16

	// To make the size of type T1 be a multiple
	// of the alignment guarantee of T1, on 64-bit
	// architectures, 6 bytes need to be padded
	// here, and on 32-bit architectures, 2 bytes
	// need to be padded here.
}
// The size of T1 is 24 (= 1 + 7 + 8 + 2 + 6)
// bytes on 64-bit architectures and is 16
// (= 1 + 3 + 8 + 2 + 2) on 32-bit architectures.

type T2 struct {
	a int8

	// To make field c 2-byte aligned, one byte
	// needs to be padded here on both 64-bit
	// and 32-bit architectures.

	c int16

	// On 64-bit architectures, to make field b
	// 8-byte aligned, 4 bytes need to be padded
	// here. On 32-bit architectures, field b is
	// already 4-byte aligned, so no bytes need
	// to be padded here.

	b int64
}
// The size of T2 is 16 (= 1 + 1 + 2 + 4 + 8)
// bytes on 64-bit architectures, and is 12
// (= 1 + 1 + 2 + 8) on 32-bit architectures.
```

