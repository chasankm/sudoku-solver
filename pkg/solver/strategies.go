package solver

type Strategy interface {
	Name() StrategyName
	Apply(*Board) (bool, error)
}

type strategyFunc struct {
	name  StrategyName
	apply func(*Board) error
}

func (s strategyFunc) Name() StrategyName {
	return s.name
}

func (s strategyFunc) Apply(board *Board) (bool, error) {
	before := board.totalMarks()
	if err := s.apply(board); err != nil {
		return false, err
	}
	return before != board.totalMarks(), nil
}

var orderedStrategies = []Strategy{
	strategyFunc{name: NakedQuadsStrategy, apply: (*Board).eliminateNQ},
	strategyFunc{name: NakedTriplesStrategy, apply: (*Board).eliminateNT},
	strategyFunc{name: NakedPairsStrategy, apply: (*Board).eliminateNP},
	strategyFunc{name: XYWingsStrategy, apply: (*Board).eliminateXYWings},
	strategyFunc{name: XYZWingsStrategy, apply: (*Board).eliminateXYZWings},
	strategyFunc{name: XWingsStrategy, apply: (*Board).eliminateXWings},
	strategyFunc{name: SwordFishStrategy, apply: (*Board).eliminateSwordFish},
	strategyFunc{name: HiddenSingleStrategy, apply: (*Board).eliminateHS},
	strategyFunc{name: HiddenQuadsStrategy, apply: (*Board).eliminateHQ},
	strategyFunc{name: HiddenTripletsStrategy, apply: (*Board).eliminateHT},
	strategyFunc{name: HiddenPairsStrategy, apply: (*Board).eliminateHP},
}

func (b *Board) applyStrategies() (bool, error) {
	for _, strategy := range orderedStrategies {
		changed, err := strategy.Apply(b)
		if err != nil {
			return false, err
		}
		if changed {
			b.addStrategy(strategy.Name())
			return true, nil
		}
	}
	return false, nil
}
