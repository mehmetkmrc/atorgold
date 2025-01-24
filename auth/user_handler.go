package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"atorgold/database"
	"atorgold/dto"
	"atorgold/models"
	"atorgold/response"

	"github.com/gofiber/fiber/v3"
)

type (
	
	UserAccess struct {
		User 			*models.User
		AccessToken 	string 			`json:"access_token"`
		AccessPublic 	string 			`json:"access_public"`
		RefreshToken 	string 			`json:"refresh_token"`
		RefreshPublic 	string 			`json:"refresh_public"`
	}

	UserRepository struct{
		db *pgxpool.Pool
	}
)

func Login(c fiber.Ctx) error {
	reqBody := new(dto.UserLoginRequest)
	
	body := c.Body()
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}
	
	Login := func (ctx context.Context, email, password string) (*UserAccess, error) {
		r := &UserRepository{db: database.DBPool}
		GetUserPassword := func (ctx context.Context, r *UserRepository, email string) (string, error) {
			var password string
			query := `
			SELECT password 
			FROM users 
			WHERE email = $1;
			`
			err := r.db.QueryRow(ctx, query, email).Scan(&password)
			if err != nil {
				return "", err
			}
			return password, nil
		}
		
		userPassword, err := GetUserPassword(ctx, r, email)
		if err != nil {
			return nil, err
		}
		if strings.Compare(password, userPassword) != 0 {
			return nil, errors.New("Şifreler eşleşmiyor")
		}

		GetByEmail := func (ctx context.Context, r *UserRepository, email string) (*models.User, error){
			userQuery := struct {
				UserID		  string
				Name	  	  string
				Surname   	  string
				Email	  	  string
				Password  	  string
				CreatedAt 	  time.Time
			}{}
			query := `
			SELECT CAST(user_id AS VARCHAR(64)) as ID, 
			   first_name, 
			   last_name, 
			   email, 
			   password, 
			   created_at 
			FROM Users 
			WHERE Email = $1 
				  AND password IS NOT NULL 
				  AND email IS NOT NULL;
			`
			err := r.db.QueryRow(ctx, query, email).Scan(&userQuery.UserID, &userQuery.Name, &userQuery.Surname, &userQuery.Email, &userQuery.Password, &userQuery.CreatedAt)
			if err != nil{
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

		userModel, err := GetByEmail(ctx, r ,email)
		if err != nil {
			return nil, err
		}

		CreateToken := func (userID string, tokenTTL time.Duration) (string, string, *Payload, error) {
			duration := tokenTTL

			NewPayload := func (userID string, duration time.Duration) (*Payload, error) {
				payload := &Payload{
					ID:        userID,
					IssuedAt:  time.Now(),
					ExpiredAt: time.Now().Add(duration),
				}
				return payload, nil
			}

			payload, err := NewPayload(userID, duration)
			if err != nil {
				return "", "", nil, err
			}
		
			tokenPaseto := paseto.NewToken()
			tokenPaseto.SetExpiration(payload.ExpiredAt)
			tokenPaseto.SetIssuedAt(payload.IssuedAt)
			tokenPaseto.SetString("id", payload.ID)
			secretKey := paseto.NewV4AsymmetricSecretKey()
			publicKey := secretKey.Public().ExportHex()
			encrypted := tokenPaseto.V4Sign(secretKey, nil)
		
			return encrypted, publicKey, payload, nil
		}

		accessToken, publicKey, accessPayload, err := CreateToken(userModel.UserID, time.Hour*24)
		if err != nil {
			return nil, err
		}

		CreateRefreshToken := func (refreshTTL time.Duration, payload *Payload) (string, string, error) {
			tokenPaseto := paseto.NewToken()
			payload.ExpiredAt = payload.ExpiredAt.Add(refreshTTL)
			tokenPaseto.SetExpiration(payload.ExpiredAt)
			tokenPaseto.SetIssuedAt(payload.IssuedAt)
			tokenPaseto.SetString("id", payload.ID)
			secretKey := paseto.NewV4AsymmetricSecretKey()
			publicKey := secretKey.Public().ExportHex()
			encrypted := tokenPaseto.V4Sign(secretKey, nil)
			return encrypted, publicKey, nil
		}

		refreshToken, refreshPublicKey, err := CreateRefreshToken( time.Hour*24 ,accessPayload)
		if err != nil {
			return nil, err
		}

		NewUserAccess := func (user *models.User, accessToken, accessPublic, refreshToken, refreshPublic string) *UserAccess{
			return &UserAccess{
				User: user,
				AccessToken: accessToken,
				AccessPublic: accessPublic,
				RefreshToken: refreshToken,
				RefreshPublic: refreshPublic,
			}
		}
	
		sessionModel := NewUserAccess(userModel, accessToken, publicKey, refreshToken, refreshPublicKey)
	
		return sessionModel, nil
	}
	

	userData, err := Login(c.Context(), reqBody.Email, reqBody.Password)
	if err != nil {
		return response.Error_Response(c, "error while trying to login", err, nil, fiber.StatusBadRequest)
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

	userResponse := GetUserModelToDto(userData.User)
	bearerAccess := "Bearer " + userData.AccessToken
	fmt.Println(userData.AccessToken)
	c.Cookie(&fiber.Cookie{
		Name:     "id",
		Value:    userData.User.UserID,
		Expires:  time.Now().Add(3 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "name",
		Value:    userData.User.Name + " " + userData.User.Surname,
		Expires:  time.Now().Add(3 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     AccessToken,
		Value:    bearerAccess,
		Expires:  time.Now().Add(time.Hour * 3),
		HTTPOnly: true,
		Secure:   true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     AccessPublic,
		Value:    userData.AccessPublic,
		Expires:  time.Now().Add(time.Hour * 3),
		HTTPOnly: true,
		Secure:   true,
	})

	bearerRefresh := "Bearer " + userData.RefreshToken
	c.Cookie(&fiber.Cookie{
		Name:     RefreshToken,
		Value:    bearerRefresh,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     RefreshPublic,
		Value:    userData.RefreshPublic,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})
	

	return response.Success_Response(c, userResponse, "Kullanıcı başarıyla giriş yaptı.", fiber.StatusOK)
}

func Register(c fiber.Ctx) error {
	reqBody := new(dto.UserRegisterRequest)
	body := c.Body()
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	Register := func (ctx context.Context, first_name, last_name, email, phone, password string)(*models.User, error) {
		newUser := &models.User{
			UserID: uuid.New().String(),
			Name: first_name,
			Surname: last_name,
			Email: email,
			Phone: phone,
			Password: password,
		}

		Create := func(ctx context.Context, r *UserRepository, user *models.User) error {
	
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			user.CreatedAt = time.Now()
			user.Password = string(hashedPassword)
		
			query := `
			INSERT INTO users (user_id, first_name, last_name, email, phone, password, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7);
			`
			_, err = r.db.Exec(ctx, query, user.UserID, user.Name, user.Surname, user.Email, user.Phone, user.Password, user.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		}

		err := Create(ctx, &UserRepository{}, newUser)
		if err != nil {
			return nil, err
		}
		return newUser, nil
	}
	
	newUser, err := Register(c.Context(), reqBody.Name, reqBody.Surname, reqBody.Email, reqBody.Phone, reqBody.Password)
	if err != nil {
		return response.Error_Response(c, "Error while trying to register user", err, nil, fiber.StatusBadRequest)
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

	userResponse := GetUserModelToDto(newUser)

	return response.Success_Response(c, userResponse, "user registered succesfully", fiber.StatusCreated)
}
