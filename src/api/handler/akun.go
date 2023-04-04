package handler

import (
	"be-2/src/api/request"
	"be-2/src/api/response"
	"be-2/src/config/database"
	"be-2/src/model"
	"be-2/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	request := &request.Login{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &model.Akun{}

	if err := db.WithContext(ctx).First(data, "email", request.Email).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusUnauthorized, []string{"email atau password salah"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !util.ValidateHash(request.Password, data.Password) {
		return util.FailedResponse(c, http.StatusUnauthorized, []string{"email atau password salah"})
	}

	var bagian string
	var idProdi int
	if data.Role == string(util.ADMIN) {
		if err := db.WithContext(ctx).Table("admin").Select("bagian").Where("id", data.ID).Scan(&bagian).Error; err != nil {
			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	} else if data.Role == string(util.OPERATOR) || data.Role == string(util.MAHASISWA) {
		if err := db.WithContext(ctx).Table(data.Role).Select("id_prodi").Where("id", data.ID).Scan(&idProdi).Error; err != nil {
			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	}

	var nama string
	if err := db.WithContext(ctx).Table(data.Role).Select("nama").Where("id", data.ID).Scan(&nama).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	token := util.GenerateJWT(data.ID, idProdi, nama, data.Role, bagian)

	return util.SuccessResponse(c, http.StatusOK, response.Login{Token: token})
}

func ChangePasswordHandler(c echo.Context) error {
	request := &request.ChangePassword{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &model.Akun{}
	claims := util.GetClaimsFromContext(c)
	id := int(claims["id"].(float64))

	if err := db.WithContext(ctx).First(data, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, []string{"user tidak ditemukan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !util.ValidateHash(request.PasswordLama, data.Password) {
		return util.FailedResponse(c, http.StatusUnauthorized, []string{"password anda berbeda dengan yang lama"})
	}

	if request.PasswordBaru == "" {
		return util.FailedResponse(c, http.StatusBadRequest, []string{"password baru tidak boleh kosong"})
	}

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(request.PasswordBaru)).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func ResetPasswordHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Akun), "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, []string{"user tidak ditemukan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	password := util.GeneratePassword()

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(password)).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, map[string]string{"password": password})
}
