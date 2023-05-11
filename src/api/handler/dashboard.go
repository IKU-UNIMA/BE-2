package handler

import (
	"be-2/src/api/request"
	"be-2/src/api/response"
	"be-2/src/config/database"
	"be-2/src/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dashboardQueryParam struct {
	Fakultas int `query:"fakultas"`
	Prodi    int `query:"prodi"`
	Tahun    int `query:"tahun"`
}

func GetDashboardHandler(c echo.Context) error {
	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	condition := ""
	if queryParams.Tahun > 2000 {
		condition = fmt.Sprintf("AND YEAR(created_at) = %d", queryParams.Tahun)
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.Dashboard{}

	var target float64
	targetQuery := fmt.Sprintf(`
	SELECT target FROM target
	WHERE bagian = 'IKU 2' AND tahun = %d
	`, queryParams.Tahun)
	if err := db.WithContext(ctx).Raw(targetQuery).Find(&target).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Target = fmt.Sprintf("%.1f", util.RoundFloat(target))

	mhs := []struct {
		ID       int
		Fakultas string
		Jumlah   int
	}{}

	mhsQuery := `
	SELECT fakultas.id, fakultas.nama AS fakultas, COUNT(mahasiswa.id) AS jumlah FROM fakultas
	left JOIN prodi ON prodi.id_fakultas = fakultas.id
	left join mahasiswa ON mahasiswa.id_prodi = prodi.id
	GROUP BY fakultas.id ORDER BY fakultas.id
	`

	if err := db.WithContext(ctx).Raw(mhsQuery).Find(&mhs).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	query := fmt.Sprintf(`
	SELECT COUNT(id_mahasiswa) AS jumlah FROM (
		SELECT id_mahasiswa, id_fakultas from prodi
		LEFT JOIN mahasiswa ON mahasiswa.id_prodi = prodi.id
		AND prodi.jenjang IN ('S1','D3')
		LEFT JOIN kampus_merdeka ON kampus_merdeka.id_mahasiswa = mahasiswa.id
		%s
		UNION
		SELECT id_mahasiswa, id_fakultas from prodi
		LEFT JOIN mahasiswa ON mahasiswa.id_prodi = prodi.id
		LEFT JOIN prestasi ON prestasi.id_mahasiswa = mahasiswa.id
		AND prestasi.tingkat_prestasi IN ('Internasional','Nasional')
		%s
	) a
	GROUP BY id_fakultas ORDER BY id_fakultas
	`, condition, condition)

	jumlahCapaian := []int{}
	if err := db.WithContext(ctx).Raw(query).Find(&jumlahCapaian).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalMahasiswa, total int
	for i := 0; i < len(mhs); i++ {
		totalMahasiswa += mhs[i].Jumlah
		total += jumlahCapaian[i]

		var persentase float64
		if mhs[i].Jumlah != 0 {
			persentase = util.RoundFloat((float64(jumlahCapaian[i]) / float64(mhs[i].Jumlah)) * 100)
		}

		data.Detail = append(data.Detail, response.DashboardDetailPerFakultas{
			ID:         mhs[i].ID,
			Fakultas:   mhs[i].Fakultas,
			Jumlah:     jumlahCapaian[i],
			Persentase: fmt.Sprintf("%.2f", persentase) + "%",
		})
	}

	data.Total = total
	data.TotalMahasiswa = totalMahasiswa

	var pencapaian float64
	if totalMahasiswa != 0 {
		pencapaian = util.RoundFloat((float64(total) / float64(totalMahasiswa)) * 100)
	}

	data.Pencapaian = fmt.Sprintf("%.2f", pencapaian) + "%"

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetKMDashboardByKategoriHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.KategoriDashboard{}

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
	if queryParams.Tahun > 2000 {
		condition = fmt.Sprintf("AND YEAR(created_at) = %d", queryParams.Tahun)
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

func GetPrestasiDashboardByTingkatHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.KategoriDashboard{}

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

func GetTotalDashboardHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.TotalDashboard{}

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
	total := 0
	if err := db.WithContext(ctx).Raw(kmQuery).Find(&total).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data = append(data, response.TotalDashboard{
		Nama:  "Kampus Merdeka",
		Total: total,
	})

	// find total prestasi
	if err := db.WithContext(ctx).Raw(prestasiQuery).Find(&total).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data = append(data, response.TotalDashboard{
		Nama:  "Prestasi",
		Total: total,
	})

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardUmumHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardUmum{}
	fakultasQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	prodiQuery := `SELECT COUNT(id) AS prodi FROM prodi`
	dosenQuery := `SELECT COUNT(id) AS dosen FROM dosen`
	mahasiswaQuery := `SELECT COUNT(id) AS mahasiswa FROM mahasiswa`

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

func InsertTargetHandler(c echo.Context) error {
	req := &request.Target{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	conds := fmt.Sprintf("bagian='%s' AND tahun=%d", util.IKU2, req.Tahun)

	if err := db.WithContext(ctx).Where(conds).Save(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
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
