package color

import "fmt"

/// Formatter ...
type ZshColor struct{}

func NewZshColor() ZshColor {
	return ZshColor{}
}

/// Foreground ...
func Foreground(color string, format string, a ...interface{}) string {
	return "%F{" + color + "}" + fmt.Sprintf(format, a...) + "%f"
}

func (f ZshColor) Blue(format string, a ...interface{}) string {
	return Foreground("blue", format, a...)
}

func (f ZshColor) Cyan(format string, a ...interface{}) string {
	return Foreground("cyan", format, a...)
}

func (f ZshColor) Yellow(format string, a ...interface{}) string {
	return Foreground("yellow", format, a...)
}

func (f ZshColor) Green(format string, a ...interface{}) string {
	return Foreground("green", format, a...)
}

func (f ZshColor) Red(format string, a ...interface{}) string {
	return Foreground("red", format, a...)
}
