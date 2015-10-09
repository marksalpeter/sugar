package sugar

import (
	"fmt"
	"strings"
)

type logger struct {
	stack []interface{}	
}

type Log func(s interface{}, args ...interface{})

type Logger interface {
	Log(s interface{}, args ...interface{})
	String() string
}

// logger will print nested logs. so if you are logging on a stack, create new Loggers and add them to logs
func NewLogger() Logger {
	return &logger{}
}

func (l *logger) Log(s interface{}, args ...interface{}) {
	if args != nil {
		if str, ok := s.(string); ok {
			l.stack = append(l.stack, fmt.Sprintf(str, args...))	
		} else {
			l.stack = append(l.stack, s)
			l.stack = append(l.stack, args...)	
		}
	} else if s != nil {
		l.stack = append(l.stack, s)
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
			result += fmt.Sprintf(" %s %+s \n", yellowColor(tag), yellowColor(l.stack[i]))			
		}
	}
	return result
}