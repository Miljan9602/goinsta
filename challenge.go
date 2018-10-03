package goinsta

type Challenges struct {
	inst *Instagram
}

func newChallenge(inst *Instagram) *Challenges {
	challenges := &Challenges{
		inst: inst,
	}
	return challenges
}
