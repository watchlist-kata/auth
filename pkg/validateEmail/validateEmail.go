package validateEmail

import "regexp"

func IsValidEmail(email string) bool {
	return regexp.MustCompile(`([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z]+)`).MatchString(email)
}
