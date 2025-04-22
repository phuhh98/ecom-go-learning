package dtos

// DOCLoginResponse represents a login response for swagger
// @Description Login response with tokens and user data
type DOCLoginResponse struct {
	AccessToken  string          `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string          `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64           `json:"expires_in" example:"3600"`
	User         DOCUserResponse `json:"user"`
}

// DOCUserResponse represents a user data response for swagger
// @Description User data returned in authentication responses
type DOCUserResponse struct {
	ID       uint     `json:"id" example:"1"`
	Name     string   `json:"name" example:"John Doe"`
	Email    string   `json:"email" example:"john.doe@example.com"`
	Roles    []string `json:"roles" example:"[\"user\"]"`
	Verified bool     `json:"verified" example:"true"`
}

// DOCTokenResponse represents a token refresh response for swagger
// @Description Response after refreshing authentication tokens
type DOCTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int64  `json:"expires_in" example:"3600"`
}

// DOCRegistrationResponse represents a registration response for swagger
// @Description Response after successful user registration
type DOCRegistrationResponse struct {
	ID    uint   `json:"id" example:"1"`
	Email string `json:"email" example:"john.doe@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

// DOCTokenPair represents authentication tokens for swagger
// @Description Authentication token pair with access and refresh tokens
type DOCTokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
}

// DOCRefreshTokenRequest represents a refresh token request for swagger
// @Description Request to refresh an access token
type DOCRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"`
}

// DOCLogoutRequest represents a logout request for swagger
// @Description Request to logout and invalidate a refresh token
type DOCLogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// DOCVerifyTokenRequest represents a token verification request for swagger
// @Description Request to verify an access token
type DOCVerifyTokenRequest struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"`
}
