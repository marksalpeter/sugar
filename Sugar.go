package sugar

import (
	"fmt"
	"time"
	"testing"
	"os"
)

type Sugar interface {
	Assert(string, Test) Sugar
	Warn(string, Test) Sugar
	Must(string, Test) Sugar
	Title(string) Sugar
}

type Test func (Log) bool

type sugar struct {
	t *testing.T
}

// if no testing.T is initialized, the test will just call os.Exit(0). this is handy for TestMain functions
func New (t *testing.T) Sugar {
	return &sugar{ t: t }
}

// writes a failure message, and marks the test as a failure if isPassed() returns false, but continues execution of the test
func (s *sugar) Assert(name string, isPassed Test) Sugar {
	startTime := time.Now()
	l := NewLogger()
	if isPassed(l.Log) {
		if testing.Verbose() {
			fmt.Printf("%s	%20s	%s\n", greenColor("PASS"), cyanColor(time.Now().Sub(startTime)), name)
			fmt.Print(l)
		}
	} else {
		fmt.Printf("%s	%20s	%s\n", redColor("FAIL"), cyanColor(time.Now().Sub(startTime)), name)
		fmt.Print(l)
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
			fmt.Printf("%s	%20s	%s\n", greenColor("PASS"), cyanColor(time.Now().Sub(startTime)), name)
			fmt.Print(l)
		}
	} else {
		fmt.Printf("%s	%20s	%s\n", yellowColor("WARN"), cyanColor(time.Now().Sub(startTime)), name)
		fmt.Print(l)
	}
	return s
}

// writes a warning message and fails the test immediatel if isPassed() returns false. the test will not continue to execute
func (s *sugar) Must(name string, isPassed Test) Sugar {
	startTime := time.Now()
	l := NewLogger()
	if isPassed(l.Log) {
		if testing.Verbose() {
			fmt.Printf("%s	%20s	%s\n", greenColor("PASS"), cyanColor(time.Now().Sub(startTime)), name)
			fmt.Print(l)
		}
	} else {
		fmt.Printf("%s	%20s	%s\n", redColor("FATAL"), cyanColor(time.Now().Sub(startTime)), name)
		fmt.Print(l)
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
		fmt.Printf("\n%20s\n", grayUnderlineColor(">>> " + title + " <<<"))
	}
	return s
}