//
// Sugar is a wrapper around testing.T that makes tests more beautiful and readible in the terminal, and more elegant 
// and syntactically clear in your test files.
//
// Terminal Improvements
//
// Sugars terminal output emphasizes readability and plain english descriptions of test cases. To make problems easier to spot,
// passing tests are silet by default (unless the -v flag is passed), and test results are displayed in different colors:
// red for failing, green for passing, and yellow for warnings and logs. Finally, sugars Logger also supports
// nested logging which can especially come in handy if you're parsing collections of data recursively
//
// Example
//
// This is some typical testing.T output:
//  === RUN   TestSugarMust
//  --- PASS: TestSugarMust (0.00s)
//  	examples_test.go:13: this must be true or t.FailNow() will be called
//  === RUN   TestSugarWarn
//  --- PASS: TestSugarWarn (0.00s)
//  	examples_test.go:17: this won't cause a failure, it just prints a warning
//  === RUN   TestSugarAssert
//  --- FAIL: TestSugarAssert (0.00s)
//  	examples_test.go:21: this must be true or t.Fail() will be called!
//  === RUN   TestSugarLog
//  --- PASS: TestSugarLog (0.00s)
//  	examples_test.go:25: t.Logf works like fmt.Printf
//  	examples_test.go:26: here are some structs &{1} &{2} &{3}
//  	examples_test.go:27:      and a nested log
//
// This is what the same tests look like with sugar:
//  === RUN   TestSugar
//  PASS	    4.634µs	this must be true or t.FailNow() will be called
//  WARN	    1.182µs	this won't cause a failure, it just prints a warning
//  FAIL	    1.224µs	this must be true or t.Fail() will be called!
//  PASS	    9.291µs	and finally this is a full demonstration of the logger
//   ┠ by default, sugar.Log works like fmt.Printf
//   ┠ but, if you just pass in structs, it will print them with their field names
//   ┠ &{Field:1}
//   ┠ &{Field:2}
//   ┠ &{Field:3}
//   ┖  ┖ finally, its possible to nest logs by createing a new logger
//  --- FAIL: TestSugar (0.00s)
// 
// Test Improvements
// 
// Tests in go can sometimes feel disjointed or hard to read. It can also sometimes be very tempting to test more than one thing
// within a single test function. Without sugar, this becomes unmanagable, quickly. Sugar addresses this problem by orgainizing 
// your tests in functions labeled by plain english descriptions of what they are trying to accomplish. 
//
// Example
//
// This is what a standard test might look like with testing.T:
//  func TestModel(t *testing.T) {
//  	var model Model
//  	if err := db.Find(&model).Error; err != nil {
//  		t.Fatal(err)
//  	}
//  	if model.Field != "this" {
//  		t.Fatal("model.Field != this")
//  	}
//  	if model.OtherField != "that" {
//  		t.Fatal("model.OtherField != that")
//  	}
//  }
// 
// Although a comprable test with sugar is slightly more verbose, it has greater code clarity and maintainability:
//  func TestModelWithSugar(t *testing.T) {
//  	var model Model
//  	sugar.New(t).
//  	Must("the model exisits in the db", func (log sugar.Log) bool {
//  		if err := db.Find(&model).Error; err != nil {
//  			log(err)
//  			return false
//  		}
//  	}).
//  	Assert("the models data is valid", func (log sugar.Log) bool {
//  		if model.Field != "this" {
//  			log("model.Field != this")
//  			 return false
//  		}
//  		if model.OtherField != "that" {
//  			log("model.OtherField != that")
//  			return false
//  		}
//  		return true
//  	})
//  }
//
// How to use sugar in the MainTest func
//
//  func TestMain
//
// Author: Mark Salpeter
//
package sugar

import (
	"fmt"
	"time"
	"testing"
	"os"
	"io"
)

type Sugar interface {
	// Flags a test as failed but the test continues execution
	Assert(string, Test) Sugar
	
	// Warns that something is wrong but the test will pass
	Warn(string, Test) Sugar
	
	// Flags a test as failed and prevents subsequent tests from running
	Must(string, Test) Sugar
	
	// Prints a title on the screen to delinate between groups of tests
	Title(string) Sugar
}

// Tests are the basis for all testing with sugar. If a test returns false, that means that it failed. 
// If a test returns true that means that it passed.
type Test func (Log) bool

type sugar struct {
	t   *testing.T
	out io.Writer
}

// creates a new sugar interface
// t is manditory 
// you can optionally pass outputs other than os.Stdout
func New (t *testing.T, outs ...io.Writer) Sugar {
	if t == nil {
		panic("t is <nil>")
	} else if outs != nil  {
		return &sugar{ t: t,  out: io.MultiWriter(outs...) }		
	} else {
		return &sugar{ t: t,  out: os.Stdout }		
	}
}

// writes a failure message, and marks the test as a failure if isPassed() returns false, but continues execution of the test
func (s *sugar) Assert(name string, isPassed Test) Sugar {
	startTime := time.Now()
	l := NewLogger()
	if isPassed(l.Log) {
		if testing.Verbose() {
			fmt.Fprintf(s.out, "%s	%20s	%s\n", greenColor("PASS"), cyanColor(time.Now().Sub(startTime)), name)
			fmt.Fprint(s.out,l)
		}
	} else {
		fmt.Fprintf(s.out,"%s	%20s	%s\n", redColor("FAIL"), cyanColor(time.Now().Sub(startTime)), name)
		fmt.Fprint(s.out,l)
		if s.t != nil {
			s.t.Fail()
		} else {
			os.Exit(0)
		}
	}
	return s
}

// writes a warning message if isPassed() returns false, but continues execution of the test and does not mark it as having failed
func (s *sugar) Warn(name string, isPassed Test) Sugar {
	startTime := time.Now()
	l := NewLogger()
	if isPassed(l.Log) {
		if testing.Verbose() {
			fmt.Fprintf(s.out,"%s	%20s	%s\n", greenColor("PASS"), cyanColor(time.Now().Sub(startTime)), name)
			fmt.Fprint(s.out,l)
		}
	} else {
		fmt.Fprintf(s.out,"%s	%20s	%s\n", yellowColor("WARN"), cyanColor(time.Now().Sub(startTime)), name)
		fmt.Fprint(s.out,l)
	}
	return s
}

// writes a warning message and fails the test immediatel if isPassed() returns false. the test will not continue to execute
func (s *sugar) Must(name string, isPassed Test) Sugar {
	startTime := time.Now()
	l := NewLogger()
	if isPassed(l.Log) {
		if testing.Verbose() {
			fmt.Fprintf(s.out,"%s	%20s	%s\n", greenColor("PASS"), cyanColor(time.Now().Sub(startTime)), name)
			fmt.Fprint(s.out,l)
		}
	} else {
		fmt.Fprintf(s.out,"%s	%20s	%s\n", redColor("FATAL"), cyanColor(time.Now().Sub(startTime)), name)
		fmt.Fprint(s.out,l)
		if s.t != nil {	
			s.t.FailNow()
		} else {
			os.Exit(0)
		}
	}
	return s
}

// draws a colorized heading
func (s *sugar) Title(title string) Sugar {
	if testing.Verbose() {
		fmt.Fprintf(s.out,"========= %s =========\n", title)
	}
	return s
}