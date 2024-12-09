package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/cache"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/global"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

type UserService struct {
	repo  repositories.UserRepo
	redis *cache.RedisCache
}

const (
	AccessToken  TokenType = "ACCESS"
	RefreshToken TokenType = "REFRESH"
)

var ctx = context.Background()

func NewUserService(ur repositories.UserRepo, client *cache.RedisCache) *UserService {
	return &UserService{
		repo: ur, redis: client,
	}
}

func (us *UserService) SignUp(userDto dto.AuthRequest) (*dto.TokenResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Email:        userDto.Email,
		PasswordHash: string(hash),
	}
	if err := us.repo.SignUp(&user); err != nil {
		return nil, err
	}
	return us.generateTokens(user.Email)
}

func (us *UserService) SignIn(userDto dto.AuthRequest) (*dto.TokenResponse, error) {
	user, err := us.repo.GetUser(userDto.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDto.Password))

	if err != nil {
		return nil, global.ErrNotFound
	}
	return us.generateTokens(user.Email)
}

func (us *UserService) RefreshRequest(tokenDto dto.RefreshRequest) (*dto.TokenResponse, error) {
	email, err := authenticateToken(tokenDto.RefreshToken, RefreshToken)
	if err != nil {
		return nil, err
	}

	_, err = us.GetUser(email)
	if err != nil {
		return nil, err
	}
	return us.generateTokens(email)
}

func (us *UserService) GetUser(email string) (*models.User, error) {
	return us.repo.GetUser(email)
}

func (us *UserService) generateTokens(email string) (*dto.TokenResponse, error) {

	accessToken, err := generateJWT(email, AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateJWT(email, RefreshToken)
	if err != nil {
		return nil, err
	}
	err = us.redis.Set(ctx, "accessToken:"+email, accessToken, time.Hour*24)
	if err != nil {
		return nil, err
	}

	err = us.redis.Set(ctx, "refreshToken:"+email, refreshToken, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateJWT(email string, tokenType TokenType) (string, error) {
	var t int64 = 0
	switch tokenType {
	case AccessToken:
		t = time.Now().Add(2 * time.Hour).Unix()
	case RefreshToken:
		t = time.Now().Add(8 * time.Hour).Unix()
	default:
		t = 0
	}
	claims := jwt.MapClaims{
		"email": email,
		"type":  tokenType,
		"exp":   t,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func authenticateToken(tokenString string, expectedType TokenType) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != string(expectedType) {
			return "", errors.New("invalid token type")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return "", errors.New("token has expired")
			}
		} else {
			return "", errors.New("invalid expiration time")
		}

		if email, ok := claims["email"].(string); ok {
			return email, nil
		}
		return "", errors.New("email claim is missing")
	}

	return "", errors.New("invalid token")
}
