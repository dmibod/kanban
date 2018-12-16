package log

type Logger interface {
	Info(string)
	Debug(string)
	Error(string)
}
