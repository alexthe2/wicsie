package agents

type NoSpreading struct {
}

func CreateNoSpreading() *NoSpreading {
	return &NoSpreading{}
}

func (_ *NoSpreading) ModifyHealth(_ []*Agent) {

}

func (_ *NoSpreading) Infect(_ *Agent) {

}
