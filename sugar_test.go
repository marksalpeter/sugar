package sugar_test

import (
	"github.com/marksalpeter/sugar"
	"testing"
)

type SubStruct struct {
	ID uint
	Is bool
}

type Struct struct {
	Field      string
	SubStructs []SubStruct
	*SubStruct
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
		nestedLogger.Log("it is also possible to nest logs by creating a new logger")
		log(nestedLogger)

		log("finally, log.Compare will compare any two interfaces and log the differences")
		return log.Compare([]Struct{
			Struct{Field: "this equals the other one"},
			Struct{Field: "this does not"},
		}, []Struct{
			Struct{Field: "this equals the other one"},
			Struct{Field: "equal this one"},
		})

	})

}

func TestCopy(t *testing.T) {

	s := sugar.New(t)
	var a Struct

	s.Must("set up test struct", func(log sugar.Log) bool {
		a.Field = "field"
		a.SubStructs = []SubStruct{{
			ID: 1,
		}, {
			ID: 2,
		}}
		a.SubStruct = &SubStruct{
			ID: 3,
			Is: true,
		}
		return true
	})

	s.Assert("copy copies slices, allocates and copies pointers to structs, and copies ints, bools, and strings", func(log sugar.Log) bool {
		var b Struct
		if err := sugar.Copy(&a, &b); err != nil {
			log(err)
			return false
		}
		return log.Compare(a, b)
	})

}
