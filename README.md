# learning-golang

## learn go with tests
From one of the most well written sources that I ever found for learning programming by [Chris James](https://quii.gitbook.io/learn-go-with-tests/)

## personal notes
- There can only be one package per folder
- Go interfaces are implicit
- When passing maps to a function or method, only the pointer is copied

### structs
The structs exersice is a simple demonstration of structs, methods, interfaces and using tables for test cases.

Test run:
```bash
cd learn-go-with-tests
go test ../learn-go-with-tests/structs
```

### pointers
The code snipped above is valid, but in go struct pointers are automatically dereferenced:

```go
//this is correct
func (w *Wallet) Balance() int {
	return (*w).balance
}

//but typically we write
func (w *Wallet) Balance() int {
	return w.balance
}
```
By convention we keep the pointer receiver consistent across all methods of a struct, but technically it is not required.

### maps

When looking up keywords in a map, we can get a value and a boolean inidicator if the keyword is present or not:

```go
value, exists := somemap[somekeyword]
```

Always initialize, as it will result in a panic runtime error due to its underlying pointer pointing to nil:

```go
// this declaration will cause panic at runtime
var m map[string]string
m["key"] = "value"
// panic: assignment to entry in nil map

// always initialize
var m = map[string]string{}
var m = make(map[string]string)
// etc. 
```

### about packages
We can rename packages to anything we want, which is useful in cases of naming conflicts:

```go
import(
	"github.com/devopsforgo/mypackage"
	jpackage "github.com/johnsiilver/mypackage"
)
```

Sometimes we need to import packages, but don't actually use them in code. Golang would not allow to build code that imports unused packages. In such cases we prepend the underscore to the package. This should always be done in the main package:

```go
package main
import (
    "fmt"
    _ "sync" //Just an example
)

```

Variables and functions starting with an upper case are available outside the package (public) and lower case resources are only available within the package (private).

```go
package calculator

func Sum(a,b int) int {...}
func connectDB() {...}
```
