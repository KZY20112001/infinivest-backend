package dto

type ProfileRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Address    string `json:"address"`
	ProfileUrl string `json:"profile_url"`
}
