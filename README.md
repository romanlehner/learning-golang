# learning-golang

## learn go with tests
From one of the most well written sources that I ever found for learning programming by [Chris James](https://quii.gitbook.io/learn-go-with-tests/)

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

## personal notes
- There can only be one package per folder
- Go interfaces are implicit