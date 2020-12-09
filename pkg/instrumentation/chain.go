package instrumentation

func NewChain(insts ...Instrumentator) Instrumentator {
	return &Chain{insts}

}

type Chain struct {
	insts []Instrumentator
}

func (c *Chain) OnStart() {
	for _, i := range c.insts {
		i.OnStart()
	}
}

func (c *Chain) OnError(err error) {
	for _, i := range c.insts {
		i.OnError(err)
	}
}

func (c *Chain) OnComplete() {
	for _, i := range c.insts {
		i.OnComplete()
	}
}
