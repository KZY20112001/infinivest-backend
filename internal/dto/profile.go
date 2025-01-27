package dto

type ProfileRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	ProfileUrl string `json:"profile_url"`
}
