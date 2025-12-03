package auth

import (
	"net/mail"
	"regexp"
)

func ParseUsernameFromEmail(email *string) *string {
	if email == nil {
		return nil
	}

	atIndex := -1
	for i, ch := range *email {
		if ch == '@' {
			atIndex = i
			break
		}
	}

	var username string
	if atIndex == -1 {
		username = *email
	} else {
		username = (*email)[:atIndex]
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	cleaned := re.ReplaceAllString(username, "")

	return &cleaned
}

func IsValidEmail(email *string) bool {
	if email == nil || *email == "" {
		return false
	}
	_, err := mail.ParseAddress(*email)
	return err == nil
}
