package tokenutil

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/koropati/population-recap/domain"
	randomstr "github.com/koropati/population-recap/internal/reandomstr"
)

const (
	ErrorNotFoundInRedis = "invalid token (not found in redis)"
	ErrorInvalidToken    = "invalid token"
	RefreshToken         = "refresh_token"
	AccessToken          = "access_token"
	VerificationEmail    = "verification_email"
	ForgotPassword       = "forgot_password"
	AnonymousRole        = "anonymous"
)

func CreateAccessToken(user *domain.User, secret string, expiry int, accessTokenUsecase domain.AccessTokenUsecase) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &domain.JwtCustomClaims{
		Name: user.Name,
		ID:   user.ID.String(),
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	uuidData, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	// Set token in Redis with expiry
	err = accessTokenUsecase.Create(context.Background(), domain.AccessToken{
		ID:        uuidData,
		Token:     t,
		UserID:    user.ID,
		Revoked:   false,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: exp,
	})
	if err != nil {
		return "", err
	}

	return t, err
}

func CreateRefreshToken(user *domain.User, secret string, expiry int, accessToken string, refreshTokenUsecase domain.RefreshTokenUsecase) (refreshToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claimsRefresh := &domain.JwtCustomRefreshClaims{
		Name: user.Name,
		ID:   user.ID.String(),
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	uuidData, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	// Set token in Redis with expiry
	err = refreshTokenUsecase.Create(context.Background(), domain.RefreshToken{
		ID:        uuidData,
		Token:     rt,
		UserID:    user.ID,
		Revoked:   false,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: exp,
	})
	if err != nil {
		return "", err
	}
	return rt, err
}

func CreateForgotToken(user *domain.User, expiry int, forgotPasswordTokenUsecase domain.ForgotPasswordTokenUsecase) (forgotToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	uuidData, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	randStr := randomstr.New(true, true, true, false)
	tokenData := randStr.GenerateRandomString(24)

	err = forgotPasswordTokenUsecase.Create(context.Background(), domain.ForgotPasswordToken{
		ID:        uuidData,
		Token:     tokenData,
		UserID:    user.ID,
		Revoked:   false,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: exp,
	})
	if err != nil {
		return "", err
	}

	return tokenData, nil
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := ParseJWTToken(requestToken, secret)
	if err != nil {
		return false, err
	}
	return true, nil
}

func RevokeToken(accessToken string, secret string, accessTokenUsecase domain.AccessTokenUsecase, refreshTokenUsecase domain.RefreshTokenUsecase) error {
	// Ekstrak ID pengguna dari token
	_, _, err := ExtractIDFromToken(accessToken, secret, AccessToken, accessTokenUsecase, refreshTokenUsecase)
	if err != nil {
		return err
	}

	err = accessTokenUsecase.Revoke(context.Background(), accessToken)
	if err != nil {
		return err
	}

	err = refreshTokenUsecase.RevokeByPairToken(context.Background(), accessToken)
	if err != nil {
		return err
	}

	return nil
}

func RevokeAll(accessToken string, secret string, accessTokenUsecase domain.AccessTokenUsecase, refreshTokenUsecase domain.RefreshTokenUsecase) error {
	userId, _, err := ExtractIDFromToken(accessToken, secret, AccessToken, accessTokenUsecase, refreshTokenUsecase)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	err = accessTokenUsecase.RevokeByUserID(context.Background(), userID)
	if err != nil {
		return err
	}

	err = refreshTokenUsecase.RevokeByUserID(context.Background(), userID)
	if err != nil {
		return err
	}

	return nil

}

func ExtractIDFromToken(requestToken string, secret string, tokenType string, accessTokenUsecase domain.AccessTokenUsecase, refreshTokenUsecase domain.RefreshTokenUsecase) (userId string, userRole string, err error) {

	isExists := IsExistTokenInDB(requestToken, tokenType, accessTokenUsecase, refreshTokenUsecase)
	if !isExists {
		return "", AnonymousRole, fmt.Errorf(ErrorNotFoundInRedis)
	}

	// Parse JWT token
	claims, err := ParseJWTToken(requestToken, secret)
	if err != nil {
		return "", AnonymousRole, err
	}

	// Extract claims
	userID, _ := claims["id"].(string)
	userRole, _ = claims["role"].(string)

	return userID, userRole, nil
}

func ExtractDataFromToken(requestToken string, secret string, tokenType string, accessTokenUsecase domain.AccessTokenUsecase, refreshTokenUsecase domain.RefreshTokenUsecase) (domain.TokenData, error) {

	exists := IsExistTokenInDB(requestToken, tokenType, accessTokenUsecase, refreshTokenUsecase)

	if !exists {
		return domain.TokenData{}, fmt.Errorf(ErrorNotFoundInRedis)
	}

	// Parse JWT token
	claims, err := ParseJWTToken(requestToken, secret)
	if err != nil {
		return domain.TokenData{}, err
	}

	userId, _ := claims["id"].(string)

	userID, err := uuid.Parse(userId)
	if err != nil {
		return domain.TokenData{}, err
	}

	data := domain.TokenData{
		Name:      claims["name"].(string),
		UserID:    userID,
		Role:      claims["role"].(string),
		TokenType: tokenType,
		Token:     requestToken,
	}

	return data, nil
}

func IsExistTokenInDB(tokenData string, tokenType string, accessTokenUsecase domain.AccessTokenUsecase, refreshTokenUsecase domain.RefreshTokenUsecase) bool {
	var exist bool = false
	if tokenType == AccessToken {
		exist = accessTokenUsecase.IsValid(context.Background(), tokenData)
	} else {
		exist = refreshTokenUsecase.IsValid(context.Background(), tokenData)
	}

	return exist
}

func ParseJWTToken(requestToken string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf(ErrorInvalidToken)
	}

	return claims, nil
}
