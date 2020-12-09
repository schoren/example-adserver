package instrumentation

type Instrumentator interface {
	OnStart()
	OnError(error)
	OnComplete()
}
