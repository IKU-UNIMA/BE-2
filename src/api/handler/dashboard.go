package handler

import (
	"be-2/src/api/request"
	"be-2/src/api/response"
	"be-2/src/config/database"
	"be-2/src/model"
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

	db := database.DB
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
		SELECT id_mahasiswa, fakultas.id AS fakultas_id from fakultas
		LEFT JOIN prodi ON prodi.id_fakultas = fakultas.id
		LEFT JOIN mahasiswa ON mahasiswa.id_prodi = prodi.id
		AND prodi.jenjang IN ('S1','D3')
		LEFT JOIN kampus_merdeka ON kampus_merdeka.id_mahasiswa = mahasiswa.id
		%s
		UNION
		SELECT id_mahasiswa, fakultas.id AS fakultas_id from fakultas
		LEFT JOIN prodi ON prodi.id_fakultas = fakultas.id
		LEFT JOIN mahasiswa ON mahasiswa.id_prodi = prodi.id
		LEFT JOIN prestasi ON prestasi.id_mahasiswa = mahasiswa.id
		AND prestasi.tingkat_prestasi IN ('Internasional','Nasional')
		%s
	) a
	GROUP BY fakultas_id ORDER BY fakultas_id
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
			ID:              mhs[i].ID,
			Fakultas:        mhs[i].Fakultas,
			JumlahMahasiswa: mhs[i].Jumlah,
			Jumlah:          jumlahCapaian[i],
			Persentase:      fmt.Sprintf("%.2f", persentase) + "%",
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

func GetDashboardByFakultasHandler(c echo.Context) error {
	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	idFakultas, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	data := &response.DashboardPerProdi{}
	data.Detail = []response.DashboardDetailPerProdi{}

	fakultasConds := ""
	if idFakultas > 0 {
		fakultasConds = fmt.Sprintf("WHERE prodi.id_fakultas = %d", idFakultas)
	}

	fakultas := ""
	if err := db.WithContext(ctx).Raw("SELECT nama FROM fakultas WHERE id = ?", idFakultas).First(&fakultas).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Fakultas = fakultas

	mhs := []struct {
		Jumlah    int
		KodeProdi int
		Prodi     string
		Jenjang   string
	}{}

	mhsQuery := fmt.Sprintf(`
	SELECT COUNT(mahasiswa.id) as jumlah, prodi.kode_prodi, prodi.nama as prodi, prodi.jenjang FROM prodi
	left JOIN mahasiswa ON mahasiswa.id_prodi = prodi.id
	%s GROUP BY prodi.id ORDER BY prodi.id
	`, fakultasConds)

	if err := db.WithContext(ctx).Raw(mhsQuery).Find(&mhs).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalMahasiswa, total int
	if len(mhs) != 0 {
		condition := ""
		if queryParams.Tahun > 2000 {
			condition = fmt.Sprintf("AND YEAR(created_at) = %d", queryParams.Tahun)
		}

		query := fmt.Sprintf(`
		SELECT COUNT(id_mahasiswa) FROM (
			SELECT id_mahasiswa, prodi.id AS prodi_id from prodi
			LEFT JOIN mahasiswa ON mahasiswa.id_prodi = prodi.id
			AND prodi.jenjang IN ('S1','D3')
			LEFT JOIN kampus_merdeka ON kampus_merdeka.id_mahasiswa = mahasiswa.id
			%s
			%s
			UNION
			SELECT id_mahasiswa, prodi.id AS prodi_id from prodi
			LEFT JOIN mahasiswa ON mahasiswa.id_prodi = prodi.id
			LEFT JOIN prestasi ON prestasi.id_mahasiswa = mahasiswa.id
			AND prestasi.tingkat_prestasi IN ('Internasional','Nasional')
			%s
			%s
		) a
		GROUP BY prodi_id ORDER BY prodi_id
	`, condition, fakultasConds, condition, fakultasConds)

		jumlahCapaian := []int{}
		if err := db.WithContext(ctx).Raw(query).Find(&jumlahCapaian).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		for i := 0; i < len(mhs); i++ {
			totalMahasiswa += mhs[i].Jumlah
			total += jumlahCapaian[i]

			var persentase float64
			if mhs[i].Jumlah != 0 {
				persentase = util.RoundFloat((float64(jumlahCapaian[i]) / float64(mhs[i].Jumlah)) * 100)
			}

			prodi := fmt.Sprintf("%d - %s (%s)", mhs[i].KodeProdi, mhs[i].Prodi, mhs[i].Jenjang)
			data.Detail = append(data.Detail, response.DashboardDetailPerProdi{
				Prodi:           prodi,
				JumlahMahasiswa: mhs[i].Jumlah,
				Jumlah:          jumlahCapaian[i],
				Persentase:      fmt.Sprintf("%.2f", persentase) + "%",
			})
		}
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
	db := database.DB
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

func GetPrestasiDashboardByTingkatHandler(c echo.Context) error {
	db := database.DB
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
	db := database.DB
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
	db := database.DB
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

func GetDashboardMahasiswa(c echo.Context) error {
	id := int(util.GetClaimsFromContext(c)["id"].(float64))
	db := database.DB
	ctx := c.Request().Context()
	data := []response.TotalDashboard{}
	if err := db.WithContext(ctx).First(new(model.Mahasiswa), "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "mahasiswa tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	query := func(fitur string) string {
		return fmt.Sprintf(`
		SELECT COUNT(%s.id) AS total FROM %s WHERE %s.id_mahasiswa = %d
		`, fitur, fitur, fitur, id)
	}

	// get kampus merdeka
	if err := db.Raw(query("kampus_merdeka")).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data[0].Nama = "Kampus Merdeka"

	// get prestasi
	prestasi := 0
	if err := db.Raw(query("prestasi")).Find(&prestasi).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data = append(data, response.TotalDashboard{
		Nama:  "Prestasi",
		Total: prestasi,
	})

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

	db := database.DB
	ctx := c.Request().Context()
	conds := fmt.Sprintf("bagian='%s' AND tahun=%d", util.IKU2, req.Tahun)

	if err := db.WithContext(ctx).Where(conds).Save(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
