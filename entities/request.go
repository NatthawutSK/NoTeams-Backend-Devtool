package entities

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	contextWrapperService interface {
		BindRi(data any) error
	}

	contextWrapper struct {
		Context   *fiber.Ctx
		validator *validator.Validate
	}

	errorResponse struct {
		FailedField string
		Tag         string
		Value       interface{}
	}
)

func ContextWrapper(ctx *fiber.Ctx) contextWrapperService {
	return &contextWrapper{
		Context:   ctx,
		validator: validator.New(),
	}
}

func (c *contextWrapper) BindRi(data any) error {
	if err := c.Context.BodyParser(data); err != nil {
		log.Printf("Error: BodyParser data failed: %s", err.Error())
		return errors.New("BodyParser data failed")
	}

	if errs := c.validator.Struct(data); errs != nil {
		errMsgs := make([]string, 0)
		validationErrors := []errorResponse{}
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem errorResponse

			log.Printf("Error At Struct: %s Field: %s \n", err.StructNamespace(), err.StructField())

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value

			validationErrors = append(validationErrors, elem)
		}

		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"%s: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, " and "))

	}

	return nil
}
