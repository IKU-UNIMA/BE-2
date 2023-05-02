package validation

import (
	"be-2/src/util"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	cv.Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	if err := cv.Validator.Struct(i); err != nil {
		errs := err.(validator.ValidationErrors)
		return util.FailedResponse(http.StatusBadRequest, translate(errs))
	}

	return nil
}

func translate(errs validator.ValidationErrors) map[string]string {
	errors := map[string]string{}
	for _, e := range errs {
		errors[e.Field()] = getTagMessage(e)
	}

	return errors
}

func getTagMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "field ini wajib diisi"
	case "email":
		return "email harus berupa alamat email yang valid"
	}

	return ""
}
