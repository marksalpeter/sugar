package sugar

import (
	"fmt"
)

type logger struct {
	stack []interface{}	
}

type Log func(i ...interface{})

func (l *logger) log(i ...interface{}) {
	if i != nil {
		l.stack = append(l.stack, i...)
	}
}

func (l *logger) print() {
	for i, s := 0, len(l.stack); i < s; i++ {
		var tag string
		if i == s - 1 {
			tag = "┖"
		} else {
			tag = "┠"
		}
		fmt.Printf(" %s %+s \n", yellowColor(tag), yellowColor(l.stack[i]))
	}
}