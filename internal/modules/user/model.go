package user

import "time"

const (
	USER_TYPE_NORMAL = "normal"
	USER_TYPE_ADMIN  = "admin"
)

type User struct {
	ID                *string    `json:"id" column:"id"`
	Email             *string    `json:"email" column:"email"`
	Username          *string    `json:"username" column:"username"`
	DisplayName       *string    `json:"displayName" column:"display_name"`
	FirstName         *string    `json:"firstName" column:"first_name"`
	LastName          *string    `json:"lastName" column:"last_name"`
	PhoneNumber       *string    `json:"phoneNumber" column:"phone_number"`
	Bio               *string    `json:"bio" column:"bio"`
	Status            *string    `json:"status" column:"status"`
	UserType          *string    `json:"userType" column:"user_type"`
	AcceptPolicy      *bool      `json:"acceptPolicy" column:"accept_policy"`
	ReceiveEmail      *bool      `json:"receiveEmail" column:"receive_email"`
	PreferredLanguage *string    `json:"preferredLanguage" column:"preferred_language"`
	ProfileCreated    *bool      `json:"profileCreated" column:"profile_created"`
	EmailVerified     *bool      `json:"emailVerified" column:"email_verified"`
	ProfileImageId    *string    `json:"-" column:"profile_image_id"`
	ProfileImageId50  *string    `json:"profileImageId50" column:"profile_image_id_50"`
	ProfileImageId100 *string    `json:"profileImageId100" column:"profile_image_id_100"`
	ProfileImageId250 *string    `json:"profileImageId250" column:"profile_image_id_250"`
	CoverImageId      *string    `json:"-" column:"cover_image_id"`
	CoverImageId250   *string    `json:"coverImageId250" column:"cover_image_id_250"`
	CoverImageId500   *string    `json:"coverImageId500" column:"cover_image_id_500"`
	CoverImageId800   *string    `json:"coverImageId800" column:"cover_image_id_800"`
	CoverImageId1200  *string    `json:"coverImageId1200" column:"cover_image_id_1200"`
	LastLogin         *time.Time `json:"lastLogin" column:"last_login"`
	CreateDate        *time.Time `json:"createDate" column:"create_date"`
	UpdateDate        *time.Time `json:"updateDate" column:"update_date"`

	// sensitive data
	Password            *string    `json:"-" column:"password"`
	Secret              *string    `json:"-" column:"secret"`
	LastGetPasscodeTime *time.Time `json:"-" column:"last_get_passcode_time"`

	// response data
	MatchScore *float32 `json:"-" column:"match_score"`
}

func (User) TableName() string {
	return "users"
}
