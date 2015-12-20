Sugar is meant to simplify test output and and better organize test in syntax in go. For full documentation see the godoc link below.
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/golang/gddo)

### Why do I need this?
With sugar, your test output will look sexy, like this ;) :
![terminal](terminal.png?raw=true)

and test code will be organized more clearly like this:
```go
package main

import (
	"github.com/marksalpeter/sugar"
	"testing"
)

type Struct struct {
	Field string
}

func TestStruct(t *testing.T) {

	s := sugar.New(t)

	s.Must("this must be true or t.FailNow() will be called", func(_ sugar.Log) bool {
		return true
	})

	s.Warn("this won't cause a failure, it just prints a warning", func(_ sugar.Log) bool {
		return false
	})

	s.Assert("this must be true or t.Fail() will be called!", func(_ sugar.Log) bool {
		return false
	})

	s.Assert("and finally this is a full demonstration of the logger", func(log sugar.Log) bool {

		log("by default, %s works like fmt.Printf", "sugar.Log")

		log("but, if you just pass in structs, it will print them with their field names")
		log(&Struct{Field: "1"}, &Struct{Field: "2"}, &Struct{Field: "3"})

		nestedLogger := sugar.NewLogger()
		nestedLogger.Log("finally, its possible to nest logs by creating a new logger")
		log(nestedLogger)

		return true
	})

}
```
