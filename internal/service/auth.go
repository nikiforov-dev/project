package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"my_app/internal/model"
	"my_app/internal/repository"
	"strings"
	"time"
)

type AuthorizationService struct {
	repository *repository.Repository
}

var (
	jwtAccessSecret  = []byte("M1WnC2b0QSUKglLx/qb/LN5x7sanGkCF+pZSetVObo2MerWLf4rwJNwuQFc6Ta1n\nMKzvNIlln6HZXJGNrGDF2Gky1Yrlie3iL6A71G3CppT4ogt2xSQaOzNd0WKfGG1T\ndqlBU7M8ZOJR51uXiM9TfInEEf4ZhE+8YX9vExWsH2HokfLBiDJs6hBQZRVJlEnO\nvAqucZOGumvxGdagmsobvSySPOIhVXUSxaNqs8qSK2TSnkfM91xwUglOJPqBTG0w\nnk2J/JmyCa7rWaU1RZzH+qNMlscNkkSa+ocL+AujLj/EGi128mkYOHm87GaNrzls\n5Q5c0pzqVX7zPX6FtIlaYw==")
	jwtRefreshSecret = []byte("4+dqSq9OnL4DF5RyNGprs41KKcPsaN/5ivvQl1lvJRaNx+wQ9GYv69VCsDnruG1y\n3lOh+6eVbygnojxET4TqFeZudSsK0rQhEhhgrj9G7ObwtcxwNACAwAGCdIp7HK/N\nskB+AWBGuF9DIM+bYMBZ3GN8sDSWzcGkQtZduKE5MzzMaF+p9eXKECu4XMhQBZKE\nsefV6P8X27jOBf5wJgKNPePMDiMedF9C1HqyLGZ27WbdLFfcBqi5mcLctq39hdiH\na2tHj8Sk6//THzjw7hkJ5rTiSgmftNYXwWrLrH0Np/IWZrx43FLhyAEX/le/BrVS\n477NcmzoFtoIMCcdC+b8jw==")
)

func NewAuthorizationService(repos *repository.Repository) *AuthorizationService {
	return &AuthorizationService{repository: repos}
}

func (s *AuthorizationService) CreateUser(input model.CreateUserInput) (int, error) {
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return 0, err
	}

	input.Password = hashedPassword
	return s.repository.User.CreateUser(input)
}

func (s *AuthorizationService) UserAuthorize(input model.SignInUserInput) (int, error) {
	user, err := s.repository.User.GetUserByLogin(input.Login)
	if err != nil {
		return 0, err
	}

	if !CheckPasswordHash(input.Password, user.HashedPassword) {
		return 0, errors.New("invalid password")
	}

	return user.ID, nil
}

func (s *AuthorizationService) GenerateJwtTokenPair(userId int) (model.JwtTokenPair, error) {
	accessTokenExp := time.Now().Add(15 * time.Minute) // 15 minutes
	refreshTokenExp := time.Now().Add(24 * time.Hour)  // 24 hours

	accessTokenClaims := jwt.MapClaims{
		"user_id": userId,
		"exp":     accessTokenExp.Unix(),
	}
	refreshTokenClaims := jwt.MapClaims{
		"user_id": userId,
		"exp":     refreshTokenExp.Unix(),
	}

	generateTokenString := func(claims jwt.MapClaims, secret []byte) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(secret)
	}

	accessTokenString, err := generateTokenString(accessTokenClaims, jwtAccessSecret)
	if err != nil {
		return model.JwtTokenPair{}, err
	}
	refreshTokenString, err := generateTokenString(refreshTokenClaims, jwtRefreshSecret)
	if err != nil {
		return model.JwtTokenPair{}, err
	}

	return model.JwtTokenPair{
		AccessToken:  &accessTokenString,
		RefreshToken: &refreshTokenString,
	}, nil
}

func (s *AuthorizationService) GetUserIdFromAuthHeader(authHeader string) (int, error) {
	const BearerSchema = "Bearer "
	if !strings.HasPrefix(authHeader, BearerSchema) {
		return 0, errors.New("not a Bearer token")
	}

	tokenString := authHeader[len(BearerSchema):]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtAccessSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("token invalid or expired")
	}

	userId, ok := claims["user_id"].(int)
	if !ok {
		return 0, errors.New("user_id not found in claims")
	}

	return userId, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
