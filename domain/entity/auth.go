package entity

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rtbe/clean-rest-api/internal/config"
)

var (
	AdminRole = "ADMIN"
	UserRole  = "USER"
)

// TokenPairModel is an pair of access/refresh tokens.
//
// swagger:model
type TokenPairModel struct {
	// JWT access token
	//
	// required: true
	AccessToken string `json:"access_token,omitempty" validate:"required"`

	// JWT refresh token
	//
	// required: true
	RefreshToken string `json:"refresh_token,omitempty" validate:"required"`
}

var jwtSalt = config.New().JWTSalt

// Auth defines model for authentication.
type Auth struct {
	UserName string
	Password string
	Email    string
}

// AccessToken defines model for jwt access token.
type AccessToken struct {
	Token     string
	ExpiresAt int64
}

// RefreshToken defines a model for jwt refresh token.
type RefreshToken struct {
	UUID      string `bson:"_id"`
	UserID    string `bson:"user_id"`
	Token     string `bson:"token"`
	ExpiresAt int64  `bson:"expires_at"`
	Used      bool   `bson:"used"`
}

// TokenPair is an pair of access/refresh tokens.
//
// swagger:model
type TokenPair struct {
	// JWT access token
	//
	// required: true
	AccessToken string `json:"access_token,omitempty"`

	// JWT refresh token
	//
	// required: true
	RefreshToken string `json:"refresh_token,omitempty"`
}

// TokenPair defines a model for JWT token pair.
type JWTTokenPair struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

// AccessTokenClaims is a set of additional claims for jwt access token.
type AccessTokenClaims struct {
	User_id string
	//This field helps to bind access token to appropriate refresh token.
	Refresh_uuid string
	User_roles   []string
	jwt.StandardClaims
}

// RefreshTokenClaims is a set of additional claims for jwt refresh token.
type RefreshTokenClaims struct {
	User_id string
	UUID    string
	jwt.StandardClaims
}

// NewTokenPair creates a new pair of access and refresh jwt tokens.
func NewTokenPair(userID string, userRoles []string) (*JWTTokenPair, error) {
	refreshTokenExp := time.Now().Add(time.Minute * 15).Unix()
	refreshTokenUUID := uuid.New().String()

	refreshToken, err := createRefreshToken(userID, refreshTokenUUID, refreshTokenExp)
	if err != nil {
		return &JWTTokenPair{}, errors.Wrap(err, "creating refresh token")
	}

	accessTokenExp := time.Now().Add(time.Hour * 1).Unix()
	accessToken, err := createAccessToken(userID, refreshTokenUUID, userRoles, accessTokenExp)
	if err != nil {
		return &JWTTokenPair{}, errors.Wrap(err, "creating access token")
	}

	tokens := JWTTokenPair{
		AccessToken: AccessToken{
			Token:     accessToken,
			ExpiresAt: accessTokenExp,
		},
		RefreshToken: RefreshToken{
			UserID:    userID,
			UUID:      refreshTokenUUID,
			Token:     refreshToken,
			ExpiresAt: refreshTokenExp,
			Used:      false,
		},
	}

	return &tokens, nil
}

// createAccessToken creates a new jwt access token.
func createAccessToken(userID string, refreshUUID string, userRoles []string, expires int64) (string, error) {
	claims := AccessTokenClaims{
		User_id:      userID,
		Refresh_uuid: refreshUUID,
		User_roles:   userRoles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(jwtSalt))
	if err != nil {
		return "", errors.Wrap(err, "creating access token")
	}

	return signedToken, nil
}

// createRefreshToken creates a new jwt refresh token.
func createRefreshToken(userID, UUID string, expires int64) (string, error) {
	claims := RefreshTokenClaims{
		User_id: userID,
		UUID:    UUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(jwtSalt))
	if err != nil {
		return "", errors.Wrap(err, "creating refresh token")
	}

	return signedToken, nil
}

// ParseRefreshTokenClaims checks validity of refresh token and returns it's claims
func ParseRefreshTokenClaims(tokenString string) (*RefreshTokenClaims, error) {

	token, err := parseJWTToken(tokenString, &RefreshTokenClaims{})
	if err != nil {
		return &RefreshTokenClaims{}, errors.Wrap(err, "refresh token is not valid")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return &RefreshTokenClaims{}, errors.New("refresh token is not valid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return &RefreshTokenClaims{}, errors.New("refresh token is expired")
	}

	return claims, nil
}

// ParseAccessTokenClaims checks validity of access token and returns it's claims.
func ParseAccessTokenClaims(tokenString string) (*AccessTokenClaims, error) {
	token, err := parseJWTToken(tokenString, &AccessTokenClaims{})
	if err != nil {
		return &AccessTokenClaims{}, errors.Wrap(err, "access token is not valid")
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return &AccessTokenClaims{}, errors.New("access token is not valid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return &AccessTokenClaims{}, errors.New("access token is expired")
	}

	return claims, nil
}

// parseJWTToken parses token string into jwt token.
func parseJWTToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSalt), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "JWT token is not valid")
	}

	return token, nil
}
