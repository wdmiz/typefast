package stats

import "time"

type Stats struct {
	KeystrokeCorrect int
	KeystrokeWrong   int
}

func New() Stats {
	return Stats{0, 0}
}

func (s Stats) Accuracy() float64 {
	if s.KeystrokeCorrect <= 0 {
		return 0
	}
	return float64(s.KeystrokeCorrect) / float64(s.KeystrokeCorrect+s.KeystrokeWrong)
}

func (s Stats) WPMRaw(t time.Duration) float64 {
	if t.Minutes() <= 0 {
		return 0
	}
	return float64(s.KeystrokeCorrect+s.KeystrokeWrong) / 5.0 / t.Minutes()
}

func (s Stats) WPM(t time.Duration) float64 {
	if t.Minutes() <= 0 {
		return 0
	}
	return float64(s.KeystrokeCorrect) / 5.0 / t.Minutes()
}
