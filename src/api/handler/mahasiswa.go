package handler

import (
	"be-2/src/api/request"
	"be-2/src/api/response"
	"be-2/src/config/database"
	"be-2/src/model"
	"be-2/src/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const getMahasiswaQuery = "SELECT mahasiswa.id, nama, nin, akun.email, id_prodi, jenis_kelamin FROM mahasiswa JOIN akun where mahasiswa.id = akun.id"

type mahasiswaQueryParam struct {
	Nim   string `query:"nim"`
	Prodi int    `query:"prodi"`
	Nama  string `query:"nama"`
	Page  int    `query:"page"`
}

func GetAllMahasiswaHandler(c echo.Context) error {
	queryParams := &mahasiswaQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Mahasiswa{}
	condition := ""

	if queryParams.Prodi != 0 {
		condition = fmt.Sprintf("id_prodi = %d", queryParams.Prodi)
	}

	if queryParams.Nama != "" {
		if queryParams.Prodi != 0 {
			condition += " AND UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
		} else {
			condition = "UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
		}
		if queryParams.Nim != "" {
			if queryParams.Prodi != 0 {
				condition = " AND nim = " + queryParams.Nim
			} else {
				condition = "nim = " + queryParams.Nim
			}
		}
	}

	if err := db.WithContext(ctx).Preload("Prodi").Raw(getMahasiswaQuery).
		Offset(util.CountOffset(queryParams.Page)).Limit(20).
		Where(condition).Find(&data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{Page: queryParams.Page, Data: data})
}

func GetMahasiswaByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.Mahasiswa{}

	condition := getMahasiswaQuery + fmt.Sprintf(" AND mahasiswa.id = %d", id)
	if err := db.WithContext(ctx).Preload("Prodi").Raw(condition).First(data).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertMahasiswaHandler(c echo.Context) error {
	request := &request.Mahasiswa{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()
	akun := &model.Akun{}
	akun.Email = request.Email
	akun.Role = string(util.MAHASISWA)
	password := util.GeneratePassword()
	akun.Password = util.HashPassword(password)

	if err := tx.WithContext(ctx).Create(akun).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(c, http.StatusBadRequest, []string{"email sudah digunakan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	mahasiswa := request.MapRequest()
	mahasiswa.ID = akun.ID

	if err := tx.WithContext(ctx).Create(mahasiswa).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(c, http.StatusBadRequest, []string{"NIM sudah digunakan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, []string{err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, map[string]string{"password": password})
}

func EditMahasiswaHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err})
	}

	request := &request.Mahasiswa{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Mahasiswa), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Table("akun").Where("id", id).Update("email", request.Email).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(c, http.StatusBadRequest, []string{"email sudah digunakan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	result := request.MapRequest()
	if err := tx.WithContext(ctx).Where("id", id).Omit("password").Updates(result).Error; err != nil {
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
				return util.FailedResponse(c, http.StatusBadRequest, []string{"NIP sudah digunakan"})
			}

			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, []string{err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteMahasiswaHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Akun), id)
	if query.Error != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(c, http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
