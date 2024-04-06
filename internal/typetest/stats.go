package typetest

type Stats struct {
	KeystrokeCorrect int
	KeystrokeWrong   int
}

func NewStats() *Stats {
	return &Stats{0, 0}
}
