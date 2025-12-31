package metrics

type Noop struct {
}

func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) Inc(name string, labels map[string]string) {}
