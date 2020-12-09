package instrumentation

type apm interface {
	Start()
	Complete()
	Error()
}

func NewAPM(l apm) Instrumentator {
	return &APM{l}

}

type APM struct {
	apm apm
}

func (a *APM) OnStart() {
	a.apm.Start()
}

func (a *APM) OnError(err error) {
	a.apm.Error()
}

func (a *APM) OnComplete() {
	a.apm.Complete()
}
