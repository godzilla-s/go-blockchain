package log

const (
	logLevel = 1
)

func Debug(args ...interface{}) {
	logPrint("[DEBUG]", args...)
}

func Warn(args ...interface{}) {
	logPrint("[WARN]")
}

func Error(args ...interface{}) {
	logPrint("[ERROR]")
}

func Info(args ...interface{}) {
	logPrint("[INFO]")
}
func logPrint(levelMsg string, args ...interface{}) {

}
