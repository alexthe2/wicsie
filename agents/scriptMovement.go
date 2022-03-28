package agents

type LuaMovement struct {
	gridMap *GridMap
}

func CreateLuaMovement(gridMap *GridMap) *LuaMovement {
	return &LuaMovement{}
}

func (movement *LuaMovement) Move(agentsAround []Agent, me Agent) (float64, float64) {

}
