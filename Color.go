package sugar

import (
	"github.com/mgutz/ansi"
	"fmt"
)

var (
	green = ansi.ColorCode("green")
	yellow = ansi.ColorCode("yellow")
	red = ansi.ColorCode("red")
	cyan = ansi.ColorCode("cyan")
	gray= ansi.LightBlack
	grayUnderline = ansi.ColorCode("180+u")
	reset = ansi.ColorCode("reset")
)

func greenColor(input interface{}) string {
	return fmt.Sprintf("%s%+v%s", green, input, reset)
}

func redColor(input interface{}) string {
	return fmt.Sprintf("%s%+v%s", red, input, reset)
}

func yellowColor(input interface{}) string {
	return fmt.Sprintf("%s%+v%s", yellow, input, reset)
}

func cyanColor(input interface{}) string {
	return fmt.Sprintf("%s%+v%s", cyan, input, reset)
}

func grayColor(input interface{}) string {
	return fmt.Sprintf("%s%+v%s", gray, input, reset)
}

func grayUnderlineColor(input interface{}) string {
	return fmt.Sprintf("%s%+v%s", grayUnderline, input, reset)
}