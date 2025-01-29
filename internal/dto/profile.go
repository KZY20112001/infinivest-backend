package dto

type ProfileRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name,omitempty"`
	Address    string `json:"address"`
	ProfileUrl string `json:"profile_url,omitempty"`
	ProfileID  string `json:"profile_id,omitempty"`
}
