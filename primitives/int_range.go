package primitives

type IntRange struct {
	Start, End int
}

func (r IntRange) Length() int {
	return r.End - r.Start
}

func ZeroIntRange() IntRange {
	return IntRange{0, 0}
}
