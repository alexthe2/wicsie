package agents

type Movement interface {
	Move(agentsAround []Agent, health Status) (float64, float64)
}
