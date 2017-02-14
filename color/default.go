package color

import "github.com/fatih/color"

type DefaultColoredOutput struct{}

func NewDefaultColoredOutput() DefaultColoredOutput {
	return DefaultColoredOutput{}
}

func (f DefaultColoredOutput) Cyan(format string, a ...interface{}) string {
	return color.CyanString(format, a...)
}

func (f DefaultColoredOutput) Yellow(format string, a ...interface{}) string {
	return color.YellowString(format, a...)
}

func (f DefaultColoredOutput) Blue(format string, a ...interface{}) string {
	return color.BlueString(format, a...)
}

func (f DefaultColoredOutput) Green(format string, a ...interface{}) string {
	return color.GreenString(format, a...)
}

func (f DefaultColoredOutput) Red(format string, a ...interface{}) string {
	return color.RedString(format, a...)
}
