package logger

type Logger interface {
	Info(...interface{})
	Debug(...interface{})
	Error(...interface{})
	Infoln(...interface{})
	Debugln(...interface{})
	Errorln(...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
	Errorf(string, ...interface{})
}
