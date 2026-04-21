package domain

type Question struct {
	ID            int
	Text          string
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectOption string
	Profile       Profile
	Difficulty    int
	Weight        float64
}

func (q Question) Options() map[string]string {
	return map[string]string{
		"A": q.OptionA,
		"B": q.OptionB,
		"C": q.OptionC,
		"D": q.OptionD,
	}
}
