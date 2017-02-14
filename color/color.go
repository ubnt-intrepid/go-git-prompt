package color

import "fmt"

func Foreground(color string, format string, a ...interface{}) string {
	return fmt.Sprintf("%%F{%s}%s%%f", color, fmt.Sprintf(format, a...))
}

func BlueString(format string, a ...interface{}) string {
	return Foreground("blue", format, a...)
}

func CyanString(format string, a ...interface{}) string {
	return Foreground("cyan", format, a...)
}

func YellowString(format string, a ...interface{}) string {
	return Foreground("yellow", format, a...)
}

func GreenString(format string, a ...interface{}) string {
	return Foreground("green", format, a...)
}

func RedString(format string, a ...interface{}) string {
	return Foreground("red", format, a...)
}
