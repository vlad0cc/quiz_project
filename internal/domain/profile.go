package domain

type Profile string

const (
	ProfileAD Profile = "AD"
	ProfileZI Profile = "ZI"
)

func (p Profile) IsValid() bool {
	return p == ProfileAD || p == ProfileZI
}

func (p Profile) Label() string {
	switch p {
	case ProfileAD:
		return "Анализ данных"
	case ProfileZI:
		return "Защита информации"
	default:
		return "Не определён"
	}
}
