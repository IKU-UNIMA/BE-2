package handler

import (
	"be-2/src/api/request"
	"be-2/src/api/response"
	"be-2/src/config/database"
	"be-2/src/config/env"
	"be-2/src/config/storage"
	"be-2/src/model"
	"be-2/src/util"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type prestasiQueryParam struct {
	Prodi    int `query:"prodi"`
	Semester int `query:"semester"`
	Nim      int `query:"nim"`
	Page     int `query:"page"`
}

func GetAllPrestasiHandler(c echo.Context) error {
	queryParams := &prestasiQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := []response.Prestasi{}
	limit := 20
	condition := ""

	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	id := int(claims["id"].(float64))
	idProdi := int(claims["id_prodi"].(float64))
	nim := 0

	if role == string(util.MAHASISWA) {
		if err := db.WithContext(ctx).Table("mahasiswa").Select("nim").Where("id", id).Scan(&nim).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		condition = fmt.Sprintf("mahasiswa.nim = %d", nim)
	} else {
		if role == string(util.OPERATOR) {
			queryParams.Prodi = idProdi
		}

		if queryParams.Nim != 0 {
			condition = fmt.Sprintf("mahasiswa.nim = %d", nim)
		}

		if queryParams.Prodi != 0 {
			if condition != "" {
				condition += fmt.Sprintf(" AND mahasiswa.id_prodi = %d", queryParams.Prodi)
			} else {
				condition = fmt.Sprintf("mahasiswa.id_prodi = %d", queryParams.Prodi)
			}
		}

		if queryParams.Semester != 0 {
			if condition != "" {
				condition += fmt.Sprintf(" AND id_semester = %d", queryParams.Semester)
			} else {
				condition = fmt.Sprintf("id_semester = %d", queryParams.Semester)
			}
		}
	}

	if err := db.WithContext(ctx).Preload("Mahasiswa.Prodi.Fakultas").Preload("Semester").
		Joins("JOIN mahasiswa ON mahasiswa.id = prestasi.id_mahasiswa").
		Where(condition).
		Offset(util.CountOffset(queryParams.Page, limit)).Limit(limit).
		Find(&result).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{
		Page: queryParams.Page,
		Data: result,
	})
}

func GetPrestasiByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := &response.DetailPrestasi{}

	role := util.GetClaimsFromContext(c)["role"].(string)
	if role == string(util.MAHASISWA) {
		if err := prestasiAuthorization(c, id, db, ctx); err != nil {
			return err
		}
	}

	if err := db.WithContext(ctx).
		Preload("Mahasiswa.Prodi.Fakultas").Preload("Semester").
		Preload("DosenPembimbing.Fakultas").Preload("DosenPembimbing.Prodi").
		Table("prestasi").First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func InsertPrestasiHandler(c echo.Context) error {
	req := &request.Prestasi{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	claims := util.GetClaimsFromContext(c)
	idMahasiswa := int(claims["id"].(float64))

	if err := db.WithContext(ctx).First(new(model.Mahasiswa), "id", idMahasiswa).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "mahasiswa tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	sertifikat, _ := c.FormFile("sertifikat")
	if sertifikat == nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "sertifikat tidak boleh kosong"})
	}

	if err := util.CheckFileIsPDF(sertifikat); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	dSertifikat, err := storage.CreateFile(sertifikat, env.GetPrestasiFolderId())
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Create(req.MapRequest(
		idMahasiswa, util.CreateFileUrl(dSertifikat.Id))).Error; err != nil {
		storage.DeleteFile(dSertifikat.Id)

		return checkPrestasiError(c, err.Error())
	}

	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func EditPrestasiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	req := &request.Prestasi{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if errAuth := prestasiAuthorization(c, id, db, ctx); errAuth != nil {
		return errAuth
	}

	if err := db.WithContext(ctx).Omit("id_mahasiswa", "sertifikat").
		Where("id", id).Updates(req.MapRequest(0, "")).Error; err != nil {
		return checkPrestasiError(c, err.Error())
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeletePrestasiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if errAuth := prestasiAuthorization(c, id, db, ctx); errAuth != nil {
		return errAuth
	}

	sertifikat := ""
	if err := db.WithContext(ctx).Table("prestasi").Select("sertifikat").
		Where("id", id).Scan(&sertifikat).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	query := db.WithContext(ctx).Delete(new(model.Prestasi), "id", id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	sertifikatArr := strings.Split(sertifikat, "/")
	sertifikatId := sertifikatArr[len(sertifikatArr)-2]
	storage.DeleteFile(sertifikatId)

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func EditSertifikatPrestasiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if errAuth := prestasiAuthorization(c, id, db, ctx); errAuth != nil {
		return errAuth
	}

	sertifikat, _ := c.FormFile("sertifikat")
	if sertifikat == nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "sertifikat tidak boleh kosong"})
	}

	if err := util.CheckFileIsPDF(sertifikat); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	dSertifikat, errDrive := storage.CreateFile(sertifikat, env.GetPrestasiFolderId())
	if errDrive != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Table("prestasi").Where("id", id).
		Update("sertifikat", util.CreateFileUrl(dSertifikat.Id)).Error; err != nil {
		storage.DeleteFile(dSertifikat.Id)
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func prestasiAuthorization(c echo.Context, id int, db *gorm.DB, ctx context.Context) error {
	claims := util.GetClaimsFromContext(c)
	idMahasiswa := int(claims["id"].(float64))
	role := claims["role"].(string)

	if role == string(util.ADMIN) || role == string(util.OPERATOR) {
		return nil
	}

	result := 0
	query := db.WithContext(ctx).Table("prestasi").Select("id_mahasiswa").
		Where("id", id).Scan(&result)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "prestasi tidak ditemukan"})
	}

	if result == idMahasiswa {
		return nil
	}

	return util.FailedResponse(http.StatusUnauthorized, nil)
}

func checkPrestasiError(c echo.Context, err string) error {
	message := ""
	if strings.Contains(err, "dosen") {
		message = "dosen pembimbing"
	} else if strings.Contains(err, "semester") {
		message = "semester"
	}

	if message != "" {
		message += " tidak ditemukan"
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": message})
	}

	return util.FailedResponse(http.StatusInternalServerError, nil)
}
