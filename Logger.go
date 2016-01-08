package sugar

import (
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type logger struct {
	stack []interface{}
	lines []int
	out   io.Writer
}

// Log lines in yellow in the following format:
//  ┠ this is what a log looks line
//  ┠ &{Field:1}
//  ┖  ┖ finally, its possible to nest logs by createing a new logger
type Log func(s interface{}, args ...interface{})

// Logger will print nested logs. So, if you are logging recursively, create new Loggers and then log them after the function
// returns. You can optionally pass in different writers to be used to write the output of the writer the default output is
// io.Stdout
type Logger interface {
	// Logs lines and other loggers nested underneath a tests. Here is an example of the output:
	//  ┠ this is what a log looks line
	//  ┠ &{Field:1}
	//  ┖  ┖ finally, its possible to nest logs by createing a new logger
	Log(s interface{}, args ...interface{})

	// Compare compares interface `b` against interface `a` and logs all of the differences
	Compare(a, b interface{}) bool

	// Prints the log
	String() string
}

// NewLogger returns a logger
func NewLogger(outs ...io.Writer) Logger {
	if outs != nil {
		return &logger{out: io.MultiWriter(outs...)}
	}
	return &logger{out: os.Stdout}
}

// Log lines in yellow in the following format:
//  ┠ this is what a log looks line
//  ┠ &{Field:1}
//  ┖  ┖ finally, its possible to nest logs by createing a new logger
func (l *logger) Log(s interface{}, args ...interface{}) {
	// TODO: determine the part of the call stack that actually called this function, see #6
	_, _, line, _ := runtime.Caller(2)
	if args != nil {
		if str, ok := s.(string); ok {
			l.stack = append(l.stack, fmt.Sprintf(str, args...))
			l.lines = append(l.lines, line)
		} else {
			l.stack = append(l.stack, s)
			l.lines = append(l.lines, line)
			for _, arg := range args {
				l.stack = append(l.stack, arg)
				l.lines = append(l.lines, line)
			}
		}
	} else if s != nil {
		l.stack = append(l.stack, s)
		l.lines = append(l.lines, line)
	}
}

// Compare performs a deep reflection over two interfaces and logs any differences that it finds. It returns true if the two
// interfaces match eachother.
func (l *logger) Compare(a, b interface{}) bool {

	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	if aValue.Kind() == reflect.Ptr {
		// a and b point to non nil values
		if aValue.Elem().IsValid() && bValue.Elem().IsValid() {
			return l.Compare(aValue.Elem().Interface(), bValue.Elem().Interface())
		}
		// a and b are either both nil and they are the same pointer type
		return aValue.Elem().IsValid() == bValue.Elem().IsValid() && a == b
	} else if aValue.Kind() == reflect.Slice {
		// fail if there are different amounts of items in the slices
		if aValue.Len() != bValue.Len() {
			l.Log("%s", aValue.Type())
			nestedLogger := NewLogger()
			nestedLogger.Log("expected: %d", aValue.Len())
			nestedLogger.Log("found   : %d", bValue.Len())
			l.Log(nestedLogger)
			return false
		}
		// see if there is anything we expect that isn't in the structValues
		for i, len := 0, aValue.Len(); i < len; i++ {
			nestedLogger := NewLogger()
			if !l.Compare(aValue.Index(i).Interface(), bValue.Index(i).Interface()) {
				l.Log("%s failed at index %d", aValue.Type(), i)
				l.Log(nestedLogger)
				return false
			}
		}
		return true
	} else if aTime, ok := aValue.Interface().(time.Time); ok {
		bTime, _ := bValue.Interface().(time.Time)
		// compare times to the nearest second
		if time.Duration(math.Abs(float64(aTime.Sub(bTime)))) > time.Second {
			l.Log("expected: %v.(%s)", aTime, aValue.Type())
			l.Log("found   : %v.(%s)", bTime, bValue.Type())
			return false
		}
		return true
	} else if aValue.Kind() == reflect.Struct {
		// iterate over all of the field
		for i, len := 0, aValue.NumField(); i < len; i++ {
			aField := aValue.Type().Field(i)
			aFieldValue := aValue.FieldByName(aField.Name)
			bField := bValue.Type().Field(i)
			bFieldValue := bValue.FieldByName(bField.Name)
			nestedLogger := NewLogger()
			if !l.Compare(aFieldValue.Interface(), bFieldValue.Interface()) {
				l.Log("%s.%s", aValue.Type(), aField.Name)
				l.Log(nestedLogger)
				return false
			}
		}
		return true
	} else if a != b {
		l.Log("expected: %v.(%s)", a, aValue.Type())
		l.Log("found   : %v.(%s)", b, bValue.Type())
		return false
	}
	return true
}

func (l *logger) String() string {
	var result string
	for i, s := 0, len(l.stack); i < s; i++ {
		var tag string
		if i == s-1 {
			tag = "┖"
		} else {
			tag = "┠"
		}
		if nestedLogger, ok := l.stack[i].(Logger); ok && nestedLogger != l {
			nestedStrings := strings.Split(nestedLogger.String(), "\n")
			for j, l := 0, len(nestedStrings)-1; j < l; j++ {
				if j < l-1 || i < s-1 {
					result += fmt.Sprintf(" %s %s \n", yellowColor("┃"), nestedStrings[j])
				} else {
					result += fmt.Sprintf(" %s %s \n", yellowColor("┖"), nestedStrings[j])
				}
			}
		} else {
			result += fmt.Sprintf(" %s %+s %s \n",
				yellowColor(tag),
				yellowColor(l.stack[i]),
				grayColor(fmt.Sprintf("[line:%d]", l.lines[i])),
			)
		}
	}
	return result
}
