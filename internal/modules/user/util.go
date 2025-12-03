package user

import (
	"errors"
	"regexp"
	"time"
	"unicode"

	"github.com/example/internal/config/myconstant"
	"github.com/example/internal/i18n"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUsername(username *string) error {
	if username == nil || *username == "" {
		return errors.New(i18n.KEY_VALIDATE_EMPTY_USERNAME)
	}
	for _, r := range *username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return errors.New(i18n.KEY_VALIDATE_NOT_USERNAME)
		}
	}
	return nil
}

// return TRUE if user eligible to get new passcode
func ValidateGetPasscodeTime(user *User) (bool, *time.Time, error) {
	if user == nil {
		return false, nil, errors.New(i18n.KEY_INTERNAL_ERROR_CAN_NOT_VALIDATE_PASSCODE_NO_USER)
	}

	if user.LastGetPasscodeTime == nil {
		return true, nil, nil
	}

	releaseTime := user.LastGetPasscodeTime.Add(myconstant.PASSCODE_RESEND_DURATION)

	if releaseTime.After(time.Now()) {
		return false, &releaseTime, nil
	}

	return true, nil, nil
}

func PasswordStrengthCheck(pass *string) error {
	if pass == nil {
		return errors.New("password cannot be nil")
	}

	password := *pass

	// Length check
	if len(password) < 8 || len(password) > 20 {
		return errors.New("password must be between 8 and 20 characters")
	}

	// Must start with a letter
	if matched, _ := regexp.MatchString(`^[A-Za-z]`, password); !matched {
		return errors.New("password must start with a letter")
	}

	// Must contain at least one lowercase, one uppercase, one number, one special char
	var (
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
	)

	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func HashPassword(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))
	return err == nil
}
