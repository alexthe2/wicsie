package agents

type Change interface {
	ModifyHealth(agents []*Agent)
}
