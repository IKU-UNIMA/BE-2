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

type kmQueryParam struct {
	Prodi           int `query:"prodi"`
	Semester        int `query:"semester"`
	KategoriProgram int `query:"kategori"`
	Nim             int `query:"nim"`
	Page            int `query:"page"`
}

func GetAllKMHandler(c echo.Context) error {
	queryParams := &kmQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := []response.KampusMerdeka{}
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
			condition = fmt.Sprintf("mahasiswa.nim = %d", queryParams.Nim)
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

		if queryParams.KategoriProgram != 0 {
			if condition != "" {
				condition += fmt.Sprintf(" AND id_kategori_program = %d", queryParams.KategoriProgram)
			} else {
				condition = fmt.Sprintf("id_kategori_program = %d", queryParams.KategoriProgram)
			}
		}
	}

	if err := db.WithContext(ctx).Preload("Mahasiswa.Prodi.Fakultas").
		Preload("Semester").Preload("KategoriProgram").
		Joins("JOIN mahasiswa ON mahasiswa.id = kampus_merdeka.id_mahasiswa").
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

func GetKMByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := &response.DetailKampusMerdeka{}

	role := util.GetClaimsFromContext(c)["role"].(string)
	if role == string(util.MAHASISWA) {
		if !kmAuthorization(c, id, db, ctx) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}
	}

	if err := db.WithContext(ctx).
		Preload("Mahasiswa.Prodi.Fakultas").
		Preload("Semester").Preload("KategoriProgram").
		Preload("DosenPembimbing.Fakultas").Preload("DosenPembimbing.Prodi").
		Table("kampus_merdeka").First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func InsertKMHandler(c echo.Context) error {
	req := &request.KampusMerdeka{}
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

	suratTugas, _ := c.FormFile("surat_tugas")
	if suratTugas == nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "surat tugas tidak boleh kosong"})
	}

	if err := util.CheckFileIsPDF(suratTugas); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	dSuratTugas, err := storage.CreateFile(suratTugas, env.GetSuratTugasFolderId())
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data, err := req.MapRequest(idMahasiswa, util.CreateFileUrl(dSuratTugas.Id))
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		storage.DeleteFile(dSuratTugas.Id)

		return checkKMError(c, err.Error())
	}

	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func EditKMHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	req := &request.KampusMerdeka{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if !kmAuthorization(c, id, db, ctx) {
		return util.FailedResponse(http.StatusUnauthorized, nil)
	}

	data, errMapping := req.MapRequest(0, "")
	if errMapping != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}

	if err := db.WithContext(ctx).Omit("id_mahasiswa", "surat_tugas", "berita_acara").
		Where("id", id).Updates(data).Error; err != nil {
		return checkKMError(c, err.Error())
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteKMHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if !kmAuthorization(c, id, db, ctx) {
		return util.FailedResponse(http.StatusUnauthorized, nil)
	}

	query := db.WithContext(ctx).Delete(new(model.KampusMerdeka), "id", id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func EditSuratTugasHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if !kmAuthorization(c, id, db, ctx) {
		return util.FailedResponse(http.StatusUnauthorized, nil)
	}

	suratTugas, _ := c.FormFile("surat_tugas")
	if suratTugas == nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "surat tugas tidak boleh kosong"})
	}

	if err := util.CheckFileIsPDF(suratTugas); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	dSuratTugas, errDrive := storage.CreateFile(suratTugas, env.GetSuratTugasFolderId())
	if errDrive != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Table("kampus_merdeka").Where("id", id).
		Update("surat_tugas", util.CreateFileUrl(dSuratTugas.Id)).Error; err != nil {
		storage.DeleteFile(dSuratTugas.Id)
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func EditBeritaAcaraHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if !kmAuthorization(c, id, db, ctx) {
		return util.FailedResponse(http.StatusUnauthorized, nil)
	}

	beritaAcara, _ := c.FormFile("berita_acara")
	if beritaAcara == nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "berita acara tidak boleh kosong"})
	}

	if err := util.CheckFileIsPDF(beritaAcara); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	dBeritaAcara, errDrive := storage.CreateFile(beritaAcara, env.GetBeritaAcaraFolderId())
	if errDrive != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Where("id", id).Update("berita_acara", util.CreateFileUrl(dBeritaAcara.Id)).Error; err != nil {
		storage.DeleteFile(dBeritaAcara.Id)
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func kmAuthorization(c echo.Context, id int, db *gorm.DB, ctx context.Context) bool {
	claims := util.GetClaimsFromContext(c)
	idMahasiswa := int(claims["id"].(float64))
	role := claims["role"].(string)

	if role == string(util.ADMIN) || role == string(util.OPERATOR) {
		return true
	}

	result := 0
	if err := db.WithContext(ctx).Table("kampus_merdeka").Select("id_mahasiswa").
		Where("id", id).Scan(&result).Error; err != nil {
		return false
	}

	return result == idMahasiswa
}

func checkKMError(c echo.Context, err string) error {
	message := ""
	if strings.Contains(err, "dosen") {
		message = "dosen pembimbing"
	} else if strings.Contains(err, "semester") {
		message = "semester"
	} else if strings.Contains(err, "kategori_program") {
		message = "kategori program"
	}

	if message != "" {
		message += " tidak ditemukan"
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": message})
	}

	return util.FailedResponse(http.StatusInternalServerError, nil)
}
