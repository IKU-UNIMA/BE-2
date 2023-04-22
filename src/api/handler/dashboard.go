package handler

import (
	"be-2/src/api/response"
	"be-2/src/config/database"
	"be-2/src/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type dashboardQueryParam struct {
	Fakultas int `query:"fakultas"`
	Prodi    int `query:"prodi"`
	Tahun    int `query:"tahun"`
	Semester int `query:"semester"`
}

func GetKMDashboardByKategoriHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Dashboard{}

	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	condition := ""
	if queryParams.Tahun > 2000 {
		condition = fmt.Sprintf(" WHERE YEAR(created_at) = %d", queryParams.Tahun)
	}

	prodiJoin := ""
	if queryParams.Prodi != 0 {
		prodiJoin = "JOIN mahasiswa ON mahasiswa.id = kampus_merdeka.id_mahasiswa"

		if condition != "" {
			condition += fmt.Sprintf(" AND mahasiswa.id_prodi = %d", queryParams.Prodi)
		} else {
			condition = fmt.Sprintf(" mahasiswa.id_prodi = %d", queryParams.Prodi)
		}
	}

	query := fmt.Sprintf(`
	SELECT nama, COUNT(kampus_merdeka.id) AS jumlah FROM kampus_merdeka
	JOIN kategori_program_km ON kategori_program_km.id = kampus_merdeka.id_kategori_program
	%s %s GROUP BY id_kategori_program;
	`, prodiJoin, condition)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDetailDashboardHandler(c echo.Context) error {
	fitur := checkDashboardFitur(c.Param("fitur"))
	if fitur == "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "fitur tidak didukung"})
	}

	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	condition := ""
	if queryParams.Semester > 0 {
		condition = fmt.Sprintf("id_semester = %d", queryParams.Semester)
		queryParams.Tahun = 0
	}

	if queryParams.Tahun > 2000 {
		if condition != "" {
			condition = fmt.Sprintf(" AND YEAR(created_at) = %d", queryParams.Tahun)
		} else {
			condition = fmt.Sprintf("YEAR(created_at) = %d", queryParams.Tahun)
		}
	}

	if queryParams.Fakultas > 0 {
		if condition != "" {
			condition += fmt.Sprintf(" AND fakultas.id = %d", queryParams.Fakultas)
		} else {
			condition = fmt.Sprintf("fakultas.id = %d", queryParams.Fakultas)
		}
	}

	if condition != "" {
		condition += " WHERE "
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.DetailDashboard{}
	query := fmt.Sprintf(`
	SELECT prodi.id, prodi.kode_prodi, prodi.nama, prodi.jenjang, fakultas.id, fakultas.nama, semester.id, semester.nama, COUNT(%s.id) AS jumlah FROM %s
	JOIN semester on semester.id = %s.id_semester
	JOIN mahasiswa ON mahasiswa.id = %s.id_mahasiswa
	JOIN prodi ON prodi.id = mahasiswa.id_prodi
	JOIN fakultas ON fakultas.id = prodi.id_fakultas
	%s GROUP BY prodi.id;
	`, fitur, fitur, fitur, fitur, condition)

	detailDashboard := response.DetailDashboard{}
	rows, err := db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&detailDashboard.Prodi.ID,
			&detailDashboard.Prodi.KodeProdi,
			&detailDashboard.Prodi.Nama,
			&detailDashboard.Prodi.Jenjang,
			&detailDashboard.Fakultas.ID,
			&detailDashboard.Fakultas.Nama,
			&detailDashboard.Semester.Id,
			&detailDashboard.Semester.Nama,
			&detailDashboard.Jumlah,
		)
		data = append(data, detailDashboard)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetKMDashboardByFakultasHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.DetailDashboard{}

	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	condition := ""
	if queryParams.Tahun > 2000 {
		condition = fmt.Sprintf(" WHERE YEAR(created_at) = %d", queryParams.Tahun)
	}

	query := fmt.Sprintf(`
	SELECT fakultas.nama, COUNT(kampus_merdeka.id) AS jumlah FROM kampus_merdeka
	JOIN mahasiswa on mahasiswa.id = kampus_merdeka.id_mahasiswa
	JOIN prodi on prodi.id = mahasiswa.id_prodi
	JOIN fakultas on fakultas.id = prodi.id_fakultas
	%s GROUP BY fakultas.id;
	`, condition)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetPrestasiDashboardByTingkatHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Dashboard{}

	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	condition := ""
	if queryParams.Tahun > 2000 {
		condition = fmt.Sprintf(" WHERE YEAR(created_at) = %d", queryParams.Tahun)
	}

	prodiJoin := ""
	if queryParams.Prodi != 0 {
		prodiJoin = "JOIN mahasiswa ON mahasiswa.id = prestasi.id_mahasiswa"

		if condition != "" {
			condition += fmt.Sprintf(" AND mahasiswa.id_prodi = %d", queryParams.Prodi)
		} else {
			condition = fmt.Sprintf(" mahasiswa.id_prodi = %d", queryParams.Prodi)
		}
	}

	query := fmt.Sprintf(`
	SELECT tingkat_prestasi as nama, COUNT(id) AS jumlah FROM prestasi
	%s %s GROUP BY tingkat_prestasi;
	`, prodiJoin, condition)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetPrestasiDashboardByFakultasHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Dashboard{}
	tahun, _ := strconv.Atoi(c.Param("tahun"))
	condition := ""
	if tahun > 2000 {
		condition = fmt.Sprintf(" WHERE YEAR(created_at) = %d", tahun)
	}

	query := fmt.Sprintf(`
	SELECT fakultas.nama, COUNT(prestasi.id) AS jumlah FROM prestasi
	JOIN mahasiswa on mahasiswa.id = prestasi.id_mahasiswa
	JOIN prodi on prodi.id = mahasiswa.id_prodi
	JOIN fakultas on fakultas.id = prodi.id_fakultas
	%s GROUP BY fakultas.id;
	`, condition)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetTotalDashboardHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.TotalDashboard{}

	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	condition := ""
	if queryParams.Tahun > 2000 {
		condition = fmt.Sprintf(" WHERE YEAR(created_at) = %d", queryParams.Tahun)
	}

	var prodiKM, prodiPrestasi string
	if queryParams.Prodi != 0 {
		prodiKM = "JOIN mahasiswa ON mahasiswa.id = kampus_merdeka.id_mahasiswa"
		prodiPrestasi = "JOIN mahasiswa ON mahasiswa.id = prestasi.id_mahasiswa"

		if condition != "" {
			condition += fmt.Sprintf(" AND mahasiswa.id_prodi = %d", queryParams.Prodi)
		} else {
			condition = fmt.Sprintf("mahasiswa.id_prodi = %d", queryParams.Prodi)
		}

	}

	kmQuery := fmt.Sprintf(`SELECT COUNT(id) AS total_kampus_merdeka FROM kampus_merdeka %s %s`, prodiKM, condition)
	prestasiQuery := fmt.Sprintf(`SELECT COUNT(id) AS total_prestasi FROM prestasi %s %s`, prodiPrestasi, condition)

	// find total kampus merdeka
	if err := db.WithContext(ctx).Raw(kmQuery).First(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// find total prestasi
	if err := db.WithContext(ctx).Raw(prestasiQuery).First(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardUmumHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardUmum{}
	fakultasQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	prodiQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	dosenQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	mahasiswaQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`

	if err := db.WithContext(ctx).Raw(fakultasQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Raw(prodiQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Raw(dosenQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Raw(mahasiswaQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func checkDashboardFitur(fitur string) string {
	switch fitur {
	case "kampus-merdeka":
		return "kampus_merdeka"
	case "prestasi":
		return fitur
	}

	return ""
}
