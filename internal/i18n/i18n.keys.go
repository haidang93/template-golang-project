package i18n

const (

	// Invalid token
	KEY_AUTHORIZATION_INVALID_TOKEN = "key_authorization_invalid_token"

	// Invalid credential
	KEY_AUTHORIZATION_INVALID_CREDENTIAL = "key_authorization_invalid_credential"

	// Wrong email or password
	KEY_AUTHORIZATION_WRONG_EMAIL_OR_PASSWORD = "key_authorization_wrong_email_or_password"

	// Access token not found
	KEY_AUTHORIZATION_NO_TOKEN = "key_authorization_no_token"

	// Unexpected signing method
	KEY_AUTHORIZATION_UNEXPECTED_SIGNING_METHOD = "key_authorization_unexpected_signing_method"

	// You account's email is not verified
	KEY_AUTHORIZATION_EMAIL_NOT_VERIFIED = "key_authorization_email_not_verified"

	// This user is not active
	KEY_BAD_REQUEST_USER_NOT_ACTIVE = "key_bad_request_user_not_active"

	// This account's email has already verified
	KEY_BAD_REQUEST_ALREADY_VERIFIED = "key_bad_request_already_verified"

	// Failed to update record
	KEY_INTERNAL_ERROR_UPDATE_FAILED = "key_internal_error_update_failed"

	// Failed to create file metadata
	KEY_INTERNAL_ERROR_CREATE_FILE_META_FAILED = "key_internal_error_create_file_meta_failed"

	// User is not created
	KEY_INTERNAL_ERROR_USER_NOT_CREATED = "key_internal_error_user_not_created"

	// Can't find account secret code. Please contact support
	KEY_INTERNAL_ERROR_NO_SECRET_CODE = "key_internal_error_no_secret_code"

	// The secret code is incorrect, Please try again
	KEY_INTERNAL_ERROR_WRONG_PASSCODE = "key_internal_error_wrong_passcode"

	// Can't validate passcode. User not found
	KEY_INTERNAL_ERROR_CAN_NOT_VALIDATE_PASSCODE_NO_USER = "key_internal_error_can_not_validate_passcode_no_user"

	// This user does not exists
	KEY_NOT_FOUND_USER = "key_not_found_user"

	// This community does not exists
	KEY_NOT_FOUND_COMMUNITY = "key_not_found_community"

	// This member does not exists
	KEY_NOT_FOUND_COMMUNITY_MEMBER = "key_not_found_community_member"

	// Failed to get member information
	KEY_NOT_FOUND_MEMBER_INFOR = "key_not_found_member_infor"

	// This username is already taken
	KEY_DUPLICATE_CONSTRAINT_USERNAME = "key_duplicate_constraint_username"

	// This email is already taken
	KEY_DUPLICATE_CONSTRAINT_EMAIL = "key_duplicate_constraint_email"

	// This community id is already taken
	KEY_DUPLICATE_CONSTRAINT_COMMUNITY_IDENTIFIER = "key_duplicate_constraint_community_identifier"

	// Duplicate value violates unique constraint
	KEY_DUPLICATE_CONSTRAINT_UNIQUE_VALUE = "key_duplicate_constraint_unique_value"

	// Invalid community type
	KEY_VALIDATE_COMMUNITY_TYPE = "key_validate_community_type"

	// Invalid status
	KEY_VALIDATE_STATUS = "key_validate_status"

	// File size must not exceeds 5MB
	KEY_VALIDATE_FILE_SIZE_5MB = "key_validate_file_size_5mb"

	// This is a private community
	KEY_VALIDATE_PRIVATE_COMMUNITY = "key_validate_private_community"

	// Community indentifier invalid
	KEY_VALIDATE_COMMUNITY_ID = "key_validate_community_id"

	// Contact type is invalid
	KEY_VALIDATE_CONTACT_TYPE = "key_validate_contact_type"

	// Member role is invalid
	KEY_VALIDATE_MEMBER_ROLE = "key_validate_member_role"

	// Please enter email
	KEY_VALIDATE_EMPTY_EMAIL = "key_validate_empty_email"

	// please enter username
	KEY_VALIDATE_EMPTY_USERNAME = "key_validate_empty_username"

	// please enter phone number
	KEY_VALIDATE_EMPTY_PHONE = "key_validate_empty_phone"

	// Email invalid
	KEY_VALIDATE_NOT_EMAIL = "key_validate_not_email"

	// username invalid
	KEY_VALIDATE_NOT_USERNAME = "key_validate_not_username"

	// Wrong password
	KEY_VALIDATE_WRONG_PASSWORD = "key_validate_wrong_password"

	// Please upload no more than 5 documents
	KEY_VALIDATE_MAX_5_DOC = "key_validate_max_5_doc"

	// You don't have permission for this action
	KEY_FORBIDDEN_PERMISSION_REQUIRED = "key_forbidden_permission_required"

	// You are not a member of this community
	KEY_FORBIDDEN_NOT_A_MEMBER = "key_forbidden_not_a_member"

	// You can't join this community. Please contact community admin for more information
	KEY_FORBIDDEN_CAN_NOT_JOIN_COMMUNITY = "key_forbidden_can_not_join_community"

	// You can't accept this invitation. Please contact community admin for more information
	KEY_FORBIDDEN_CAN_NOT_ACCEPT_INVITATION = "key_forbidden_can_not_accept_invitation"

	// You don't have permission to delete this member
	KEY_FORBIDDEN_DELETE_MEMBER = "key_forbidden_delete_member"

	// This announcement is private
	KEY_FORBIDDEN_ANNOUNCEMENT = "key_forbidden_announcement"

	// Please wait {second} second(s) before you can request again
	KEY_FORBIDDEN_WAIT_X_SECOND = "key_forbidden_wait_x_second"

	// Member record already exists
	KEY_RECORD_EXISTS_MEMBER = "key_record_exists_member"

	// You are already a member
	KEY_INFORMATION_MEMBER_EXISTS = "key_information_member_exists"

	// Redis client is not initialized
	KEY_ERROR_REDIS_NOT_INITIALIZED = "key_error_redis_not_initialized"
)
