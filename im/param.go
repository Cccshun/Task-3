package im

// parameters for IM
const (
	SeedSize       = 10
	PopSize        = 40
	ActivationProb = 0.01
)

// parameters for algorithm
const (
	PC     = 0.5
	PM     = 0.5
	PL     = 0.3
	MaxGen = 150
)

// parameters for cascading failure
const (
	Alpha         = 0.5
	Beta          = 1.7
	NodeAttackPer = 0.0
	RepeatTime    = 5 // 重复次数
)
