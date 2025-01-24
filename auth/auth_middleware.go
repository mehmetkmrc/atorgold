package auth

import (
	"atorgold/dto"
	"atorgold/models"
	"atorgold/response"
	"context"
	"errors"
	"fmt"
	"time"
	"aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"
)

const (
	AuthHeader    = "Authorization"
	AccessToken   = "access_token"
	AccessPublic  = "access_public"
	RefreshToken  = "refresh_token"
	RefreshPublic = "refresh_public"
	UserDetail    = "UserDetail"
	AuthType      = "Bearer"
	AuthPayload   = "Payload"
)


var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

type(
	PasetoToken struct {
		tokenTTL   time.Duration
		refreshTTL time.Duration
	} 

	Payload struct{
		ID string `json:"id"`
		IssuedAt time.Time `json:"issued_at"`
		ExpiredAt time.Time `json:"expired_at"`
	}

	
)

func IsAuthorized(c fiber.Ctx) error {

	fmt.Println("IsAuthorized middleware çalıştı")
	if !isValidToken(c) {
		return redirectToLogin(c, fiber.StatusUnauthorized, "authorization header is not provided or invalid")
	}

	if !isValidPublicKey(c) {
		return redirectToLogin(c, fiber.StatusUnauthorized, "public key is not provided")
	}

	token := getAccessToken(c)
	publicKey := getAccessPublicKey(c)


	DecodeToken := func (pt *PasetoToken, pasetoToken, publicKeyHex string) (*Payload, error) {
		publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKeyHex)
		if err != nil {
			return nil, err
		}
	
		parser := paseto.NewParser()
		parsedToken, err := parser.ParseV4Public(publicKey, pasetoToken, nil)
		if err != nil {
			return nil, err
		}
	
		payload := new(Payload)
		expiredAt, err := parsedToken.GetExpiration()
		if err != nil {
			return nil, err
		}

		Valid := func (payload *Payload)  error {
			if !time.Now().After(payload.ExpiredAt) {
				return ErrExpiredToken
			}
			return nil
		}

		err = Valid(payload)
		if err != nil {
			return nil, err
		}
	
		issuedAt, err := parsedToken.GetIssuedAt()
		if err != nil {
			return nil, err
		}
	
		id, err := parsedToken.GetString("id")
		if err != nil {
			return nil, err
		}
	
		payload = &Payload{
			ID:        id,
			IssuedAt:  issuedAt,
			ExpiredAt: expiredAt,
		}
	
		return payload, nil
	
	}

	paseto := &PasetoToken{}
	payload, err := DecodeToken(paseto ,token, publicKey)
	if err != nil {
		return redirectToLogin(c, fiber.StatusUnauthorized, "invalid access token")
	}

	c.Locals(AuthPayload, payload)
	return c.Next()
}

func GetUserDetail(c fiber.Ctx) error {
	payload := c.Locals(AuthPayload).(*Payload)

	GetByID := func (ctx context.Context, r *UserRepository, id string) (*models.User, error) {
		userQuery := struct {
			UserID        string
			Name      	  string
			Surname   	  string
			Email     	  string
			Password  	  string
			CreatedAt 	  time.Time
		}{}
		query := `SELECT CAST(user_id AS VARCHAR(64)) as UserID, 
		   first_name, 
		   last_name, 
		   email, 
		   password, 
		   created_at 
		FROM Users 
		WHERE user_id = $1 
			  AND password IS NOT NULL 
			  AND email IS NOT NULL;
		`
		err := r.db.QueryRow(ctx, query, id).Scan(&userQuery.UserID, &userQuery.Name, &userQuery.Surname, &userQuery.Email, &userQuery.Password, &userQuery.CreatedAt)
		if err != nil {
			return nil, err
		}
	
		userData := &models.User{
			UserID:    userQuery.UserID,
			Name:      userQuery.Name,
			Surname:   userQuery.Surname,
			Email:     userQuery.Email,
			Password:  userQuery.Password,
			CreatedAt: userQuery.CreatedAt,
		}
		return userData, nil
	}

	GetUserByID := func (ctx context.Context, q *UserRepository,id string) (*models.User, error) {
		userModel, err := GetByID(ctx, &UserRepository{} , id)
		if err != nil {
			return nil, err
		}
	
		return userModel, nil
	}

	userAggregate, err := GetUserByID(c.Context(), &UserRepository{}, payload.ID)
	if err != nil {
		return response.Error_Response(c, "error while trying to get user detail", err, nil, fiber.StatusBadRequest)
	}

	GetUserModelToDto := func (userData *models.User) *dto.GetUserResponse {
		return &dto.GetUserResponse{
			UserID:    userData.UserID,
			Name:      userData.Name,
			Surname:   userData.Surname,
			Email:     userData.Email,
			CreatedAt: userData.CreatedAt,
		}
	}

	userResponse := GetUserModelToDto(userAggregate)
	c.Locals(UserDetail, userResponse)
	return c.Next()
}

func redirectToLogin(c fiber.Ctx, statusCode int, message string) error {
	c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
	return c.Redirect().To("/")
}


