package auth

import (
	"atorgold/response"

	"github.com/goccy/go-json"
	"atorgold/dto"

	"github.com/gofiber/fiber/v3"
)

func LoginValidation(c fiber.Ctx) error {
	data := new(dto.UserLoginRequest)
	body := c.Body()
	err := json.Unmarshal(body, &data)
	if err != nil {
		return response.Error_Response(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}

	validationErrors := ValidateRequestByStruct(data)
	if len(validationErrors) > 0 {
		return response.Error_Response(c, "validation failed", nil, validationErrors, fiber.StatusUnprocessableEntity)
	}

	return c.Next()
}


func RegisterValidation(c fiber.Ctx) error {
	data := new(dto.UserRegisterRequest)
	body := c.Body()
	err := json.Unmarshal(body, &data)
	if err != nil {
		return response.Error_Response(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}

	// Check if Password matches ConfirmPassword
    if data.Password != data.ConfirmPassword {
        return response.Error_Response(c, "passwords do not match", nil, nil, fiber.StatusBadRequest)
    }
	validationErrors := ValidateRequestByStruct(data)
	if len(validationErrors) > 0 {
		return response.Error_Response(c, "validation failed", nil, validationErrors, fiber.StatusUnprocessableEntity)
	}

	return c.Next()
}