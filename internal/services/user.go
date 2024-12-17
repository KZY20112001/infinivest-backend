package services

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/cache"
	"github.com/KZY20112001/infinivest-backend/internal/constants"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  repositories.UserRepo
	redis *cache.RedisCache
}

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
	return us.generateTokens(user.ID)
}

func (us *UserService) SignIn(userDto dto.AuthRequest) (*dto.TokenResponse, error) {
	user, err := us.repo.GetUserByEmail(userDto.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDto.Password))

	if err != nil {
		return nil, constants.ErrNotFound
	}
	return us.generateTokens(user.ID)
}

func (us *UserService) RefreshRequest(tokenDto dto.RefreshRequest) (*dto.TokenResponse, error) {
	id, err := authenticateToken(tokenDto.RefreshToken, constants.RefreshToken)

	if err != nil {
		return nil, err
	}

	_, err = us.GetUser(id)
	if err != nil {
		return nil, err
	}
	return us.generateTokens(id)
}

func (us *UserService) GetUser(id uint) (*models.User, error) {
	return us.repo.GetUser(id)
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	return us.repo.GetUserByEmail(email)
}

func (us *UserService) generateTokens(id uint) (*dto.TokenResponse, error) {
	accessToken, err := generateJWT(id, constants.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateJWT(id, constants.RefreshToken)
	if err != nil {
		return nil, err
	}
	err = us.redis.Set(ctx, "accessToken:"+strconv.FormatUint(uint64(id), 10), accessToken, time.Hour*24)
	if err != nil {
		return nil, err
	}

	err = us.redis.Set(ctx, "refreshToken:"+strconv.FormatUint(uint64(id), 10), refreshToken, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateJWT(id uint, tokenType constants.TokenType) (string, error) {
	var t int64 = 0
	switch tokenType {
	case constants.AccessToken:
		t = time.Now().Add(2 * time.Hour).Unix()
	case constants.RefreshToken:
		t = time.Now().Add(8 * time.Hour).Unix()
	default:
		t = 0
	}
	claims := jwt.MapClaims{
		"id":   strconv.FormatUint(uint64(id), 10),
		"type": tokenType,
		"exp":  t,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func authenticateToken(tokenString string, expectedType constants.TokenType) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != string(expectedType) {
			return 0, errors.New("invalid token type")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return 0, errors.New("token has expired")
			}
		} else {
			return 0, errors.New("invalid expiration time")
		}

		if idStr, ok := claims["id"].(string); ok {
			if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
				return uint(id), nil
			} else {

				return 0, err
			}
		}
		return 0, errors.New("ID claim is missing")
	}

	return 0, errors.New("invalid token")
}
