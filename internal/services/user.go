package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/commons"
	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/models"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(dto dto.AuthRequest) (*dto.TokenResponse, error)
	SignIn(dto dto.AuthRequest) (*dto.TokenResponse, error)
	RefreshRequest(dto dto.RefreshRequest) (*dto.TokenResponse, error)
	GetUser(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	generateTokens(id uint) (*dto.TokenResponse, error)
}

type userServiceImpl struct {
	repo repositories.UserRepo
}

func NewUserServiceImpl(ur repositories.UserRepo) *userServiceImpl {
	return &userServiceImpl{
		repo: ur,
	}
}

func (us *userServiceImpl) SignUp(dto dto.AuthRequest) (*dto.TokenResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Email:        dto.Email,
		PasswordHash: string(hash),
	}
	if err := us.repo.SignUp(&user); err != nil {
		return nil, err
	}

	return us.generateTokens(user.ID)
}

func (us *userServiceImpl) SignIn(dto dto.AuthRequest) (*dto.TokenResponse, error) {
	user, err := us.repo.GetUserByEmail(dto.Email)
	if err != nil {
		return nil, commons.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password))

	if err != nil {
		return nil, commons.ErrInvalidCredentials
	}
	return us.generateTokens(user.ID)
}

func (us *userServiceImpl) RefreshRequest(dto dto.RefreshRequest) (*dto.TokenResponse, error) {
	id, err := authenticateToken(dto.RefreshToken, commons.RefreshToken)

	if err != nil {
		return nil, err
	}

	_, err = us.GetUser(id)
	if err != nil {
		return nil, err
	}
	return us.generateTokens(id)
}

func (us *userServiceImpl) GetUser(id uint) (*models.User, error) {
	return us.repo.GetUser(id)
}

func (us *userServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	return us.repo.GetUserByEmail(email)
}

func (us *userServiceImpl) generateTokens(id uint) (*dto.TokenResponse, error) {
	accessToken, err := generateJWT(id, commons.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateJWT(id, commons.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateJWT(id uint, tokenType commons.TokenType) (string, error) {
	var t int64 = 0
	switch tokenType {
	case commons.AccessToken:
		t = time.Now().Add(2 * time.Hour).Unix()
	case commons.RefreshToken:
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
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func authenticateToken(tokenString string, expectedType commons.TokenType) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
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
