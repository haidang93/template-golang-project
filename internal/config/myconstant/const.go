package myconstant

import "time"

const (
	// file size values
	MAX_UPLOAD_FILE_SIZE int64 = 5 * 1024 * 1024

	// duration values
	PASSCODE_EXPIRATION_DURATION uint          = 60 * 30
	PASSCODE_RESEND_DURATION     time.Duration = time.Minute

	// context value
	CONTEXT_KEY_REDIS         string = "CONTEXT_KEY_REDIS"
	CONTEXT_KEY_ENV           string = "CONTEXT_KEY_ENV"
	CONTEXT_KEY_USER_TYPE     string = "CONTEXT_KEY_USER_TYPE"
	CONTEXT_KEY_USER_ID       string = "CONTEXT_KEY_USER_ID"
	CONTEXT_KEY_TRANSLATION   string = "CONTEXT_KEY_TRANSLATION"
	CONTEXT_KEY_LANGUAGE_CODE string = "CONTEXT_KEY_LANGUAGE_CODE"
)
