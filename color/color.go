package color

import "fmt"

func BlueString(format string, a ...interface{}) string {
	return fmt.Sprintf("%%F{blue}%s%%f", fmt.Sprintf(format, a...))
}

func CyanString(format string, a ...interface{}) string {
	return fmt.Sprintf("%%F{cyan}%s%%f", fmt.Sprintf(format, a...))
}

func YellowString(format string, a ...interface{}) string {
	return fmt.Sprintf("%%F{yellow}%s%%f", fmt.Sprintf(format, a...))
}

func GreenString(format string, a ...interface{}) string {
	return fmt.Sprintf("%%F{green}%s%%f", fmt.Sprintf(format, a...))
}

func RedString(format string, a ...interface{}) string {
	return fmt.Sprintf("%%F{red}%s%%f", fmt.Sprintf(format, a...))
}
