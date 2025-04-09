## Issues: newController is undefined
```html
Error:
./main.go:34:7: undefined: newController
```
## Cause: 
I ran command `go run main.go` but the `newController` function is written in `controller.go`. In this scenario 
the `controller.go` file is not compiled.

## Solution : 
Run all the go file.
```shell
$ go run .
```

