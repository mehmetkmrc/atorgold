package response

import (
	"github.com/gofiber/fiber/v3"
)

type ValidationMessage struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Message          string               `json:"message"`
	ValidationErrors []*ValidationMessage `json:"validation_errors,omitempty"`
	Error            string               `json:"error,omitempty"`
	Status           int                  `json:"status"`
}

type SuccessResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type SuccessResponseWithoutData struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type RedirectionResponse struct {
	Url string `json:"url"`
}

func Success_Response(c fiber.Ctx, data interface{}, message string, status int) error {
	if data != nil {
		return c.Status(status).JSON(SuccessResponseBuilder(data, message, status))
	}

	return c.Status(status).JSON(SuccessResponseWithoutDataBuilder(message, status))
}

func Error_Response(c fiber.Ctx, message string, err error, validationErrors []*ValidationMessage, status int) error {
	if err != nil {
		return c.Status(status).JSON(ErrorResponseBuilder(message, err, status))
	}

	return c.Status(status).JSON(ValidationErrorsResponseBuilder(message, validationErrors, status))
}

func ErrorResponseBuilder(message string, err error, status int) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Error:   err.Error(),
		Status:  status,
	}
}

func ValidationErrorsResponseBuilder(message string, validationErrors []*ValidationMessage, status int) *ErrorResponse {
	return &ErrorResponse{
		Message:          message,
		ValidationErrors: validationErrors,
		Status:           status,
	}
}

func SuccessResponseBuilder[T any](data T, message string, status int) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		Data:    data,
		Message: message,
		Status:  status,
	}
}

func SuccessResponseWithoutDataBuilder(message string, status int) *SuccessResponseWithoutData {
	return &SuccessResponseWithoutData{
		Message: message,
		Status:  status,
	}
}

func LimitReachedFunc(c fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(ErrorResponse{
		Message: fiber.ErrTooManyRequests.Message,
		Status:  fiber.StatusTooManyRequests,
	})
}
