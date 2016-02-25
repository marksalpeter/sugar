package sugar_test

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
		nestedLogger.Log("it is also possible to nest logs by creating a new logger")
		log(nestedLogger)

		log("finally, interface comparisons are built in to the logger as well")
		nestedLogger = sugar.NewLogger()
		if !nestedLogger.Compare([]Struct{
			Struct{Field: "this equals the other one"},
			Struct{Field: "this field does not"},
		}, []Struct{
			Struct{Field: "this equals the other one"},
			Struct{Field: "equal this one"},
		}) {
			log(nestedLogger)
		}

		return true
	})

}
