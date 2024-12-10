package global

type TokenType string

const (
	AccessToken  TokenType = "ACCESS"
	RefreshToken TokenType = "REFRESH"
)
