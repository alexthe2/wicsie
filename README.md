# What if Covid started in Europe (WiCsiE)
What if Covid started in Europe - A simulation for Scientific Computing at the RUG

## Paper
The paper can is in the repository under WICSIE_PAPER.pdf

## What is WICSIE
WICSIE is a simulation that was created to simulate the spread of infcetious diseases throughout Europe

## How can I run WICSIE
WICSIE is a pure Go implementation, that is dependent on some other (open-source) Go packages. You can just download it and run it with
```bash
go run main.go
```

## Notes
WICSIE was more of a PoC than a fully fledged software. It does allow for a lot of modification (see .behaviour files in config) and healthConstants, it has not been finished due to the time Constraints of the course

At a high resolution of 15 million Agents (weight=1), WICSIE ran for 30hrs on the RUG's High Computation Cluster and produced "just" 200 days, this should be noted, when running your own simulations
