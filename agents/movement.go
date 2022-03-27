package agents

type Movement interface {
	Move(agentsAround []Agent, me Agent) (float64, float64)
}
