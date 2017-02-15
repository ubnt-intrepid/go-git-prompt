package color

import "fmt"

/// Formatter ...
type ZshColoredOutput struct{}

func NewZshColoredOutput() Colored {
	return &ZshColoredOutput{}
}

/// Foreground ...
func Foreground(color string, format string, a ...interface{}) string {
	return "%F{" + color + "}" + fmt.Sprintf(format, a...) + "%f"
}

func (f ZshColoredOutput) Blue(format string, a ...interface{}) string {
	return Foreground("blue", format, a...)
}

func (f ZshColoredOutput) Cyan(format string, a ...interface{}) string {
	return Foreground("cyan", format, a...)
}

func (f ZshColoredOutput) Yellow(format string, a ...interface{}) string {
	return Foreground("yellow", format, a...)
}

func (f ZshColoredOutput) Green(format string, a ...interface{}) string {
	return Foreground("green", format, a...)
}

func (f ZshColoredOutput) Red(format string, a ...interface{}) string {
	return Foreground("red", format, a...)
}
