package agents

type Status int

const (
	Healthy Status = iota
	Incubated
	Infected
	UnknownInfected
	Cured
	Vaccinated
)

func (status Status) GetColor() (float64, float64, float64) {
	switch status {
	case Healthy:
		return 0, 1, 0

	case Infected, UnknownInfected:
		return 1, 0, 0

	case Cured:
		return 0.8, 0.8, 0.8

	case Vaccinated:
		return 0, 0.2, 1

	case Incubated:
		return 1, 0.71, 0.75
	}

	return 0, 0, 0
}
