package validation

import (
	"be-2/src/util"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		errs := err.(validator.ValidationErrors)
		return util.FailedResponse(http.StatusBadRequest, translate(errs))
	}

	return nil
}

func translate(errs validator.ValidationErrors) map[string]string {
	errors := map[string]string{}
	for _, e := range errs {
		key := createErrorKey(e.Field())
		errors[key] = getTagMessage(e)
	}

	return errors
}

func createErrorKey(key string) string {
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchAllCap.ReplaceAllString(key, "${1}_${2}")
	return strings.ToLower(snake)
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
