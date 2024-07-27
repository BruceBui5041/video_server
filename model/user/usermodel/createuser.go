package usermodel

// CreateUser represents the data needed to create a new user
type CreateUser struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	AuthType          string `json:"auth_type"`
	AuthProviderID    string `json:"auth_provider_id,omitempty"`
	AuthProviderToken string `json:"auth_provider_token,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
}
