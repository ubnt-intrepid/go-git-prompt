package color

type Colored interface {
	Blue(format string, a ...interface{}) string
	Cyan(format string, a ...interface{}) string
	Yellow(format string, a ...interface{}) string
	Green(format string, a ...interface{}) string
	Red(format string, a ...interface{}) string
}
