package common

const (
	EN = "en"
	PT = "pt"
)

func IsLangValid(lan string) bool {
	if lan != EN && lan != PT {
		return false
	}
	return true
}

func CheckLang(lan string) string {
	// if lan != EN && lan != PT {
	// 	return EN
	// }
	// return lan
	return EN
}
