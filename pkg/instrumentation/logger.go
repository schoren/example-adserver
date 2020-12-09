package instrumentation

type logger interface {
	Printf(format string, v ...interface{})
}

func NewLogger(l logger) Instrumentator {
	return &Logger{l}

}

type Logger struct {
	log logger
}

func (l *Logger) OnStart() {
	l.log.Printf("Start")
}

func (l *Logger) OnError(err error) {
	l.log.Printf("Error: %s", err.Error())
}

func (l *Logger) OnComplete() {
	l.log.Printf("Start")
}
