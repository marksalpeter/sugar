package sugar

import (
	"fmt"
	"strings"
	"io"
	"os"
	"runtime"
)

type logger struct {
	stack []interface{}	
	lines []int
	out io.Writer
}

// Logs lines in yellow in the following format:
//  ┠ this is what a log looks line
//  ┠ &{Field:1}
//  ┖  ┖ finally, its possible to nest logs by createing a new logger
type Log func(s interface{}, args ...interface{})

type Logger interface {
	// Logs lines and other loggers nested underneath a tests
	Log(s interface{}, args ...interface{})
	
	// Prints the log
	String() string
}

// Logger will print nested logs. So, if you are logging recursively, create new Loggers and then log them after the function returns.
// You can optionally pass in different writers to be used to write the output of the writer.
// the default output is io.Stdout
func NewLogger(outs ...io.Writer) Logger {
	if outs != nil {
		return  &logger{ out : io.MultiWriter(outs...) }
	} else {
		return &logger{ out : os.Stdout }		
	}
}

func (l *logger) Log(s interface{}, args ...interface{}) {
	// TODO: determine the part of the call stack that actually called this function
	_, file, line, _ := runtime.Caller(2)
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

func (l *logger) String() string {
	var result string
	for i, s := 0, len(l.stack); i < s; i++ {
		var tag string
		if i == s - 1 {
			tag = "┖"
		} else {
			tag = "┠"
		}
		if nestedLogger, ok := l.stack[i].(Logger); ok && nestedLogger != l {
			nestedStrings := strings.Split(nestedLogger.String(), "\n")
			for j, l := 0, len(nestedStrings) - 1; j < l; j++  {
				if j < l - 1 || i < s - 1 {
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