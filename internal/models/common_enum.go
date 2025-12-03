package models

const (
	COMMON_STATUS_OPEN    string = "open"
	COMMON_STATUS_CLOSED  string = "closed"
	COMMON_STATUS_DELETED string = "deleted"
	COMMON_STATUS_PENDING string = "pending"
)

func IsValidCommonStatus(status *string) bool {
	switch *status {
	case COMMON_STATUS_OPEN,
		COMMON_STATUS_CLOSED,
		COMMON_STATUS_DELETED,
		COMMON_STATUS_PENDING:
		return true
	}
	return false
}

const (
	VERIFICATION_STATUS_CREATED   string = "created"
	VERIFICATION_STATUS_PENDING   string = "pending"
	VERIFICATION_STATUS_VERIFIED  string = "verified"
	VERIFICATION_STATUS_RETRACTED string = "retracted"
)

func IsValidVerificationStatus(status string) bool {
	switch status {
	case VERIFICATION_STATUS_CREATED,
		VERIFICATION_STATUS_PENDING,
		VERIFICATION_STATUS_VERIFIED,
		VERIFICATION_STATUS_RETRACTED:
		return true
	}
	return false
}
