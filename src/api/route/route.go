package route

import (
	"be-2/src/api/handler"
	"be-2/src/util/validation"

	customMiddleware "be-2/src/api/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer() *echo.Echo {
	app := echo.New()
	app.Use(middleware.CORS())

	app.Validator = &validation.CustomValidator{Validator: validator.New()}

	app.GET("", func(c echo.Context) error {
		return c.JSON(200, "Welcome to IKU 2 API")
	})

	v1 := app.Group("/api/v1")

	fakultas := v1.Group("/fakultas", customMiddleware.Authentication)
	fakultas.GET("", handler.GetAllFakultasHandler)
	fakultas.GET("/:id", handler.GetFakultasByIdHandler)
	fakultas.POST("", handler.InsertFakultasHandler, customMiddleware.GrantAdminUmum)
	fakultas.PUT("/:id", handler.EditFakultasHandler, customMiddleware.GrantAdminUmum)
	fakultas.DELETE("/:id", handler.DeleteFakultasHandler, customMiddleware.GrantAdminUmum)

	prodi := v1.Group("/prodi", customMiddleware.Authentication)
	prodi.GET("", handler.GetAllProdiHandler)
	prodi.GET("/:id", handler.GetProdiByIdHandler)
	prodi.POST("", handler.InsertProdiHandler, customMiddleware.GrantAdminUmum)
	prodi.PUT("/:id", handler.EditProdiHandler, customMiddleware.GrantAdminUmum)
	prodi.DELETE("/:id", handler.DeleteProdiHandler, customMiddleware.GrantAdminUmum)

	semester := v1.Group("/semester", customMiddleware.Authentication)
	semester.GET("", handler.GetAllSemesterHandler)
	semester.POST("", handler.InsertSemesterHandler, customMiddleware.GrantAdminUmum)
	semester.PUT("/:id", handler.EditSemesterHandler, customMiddleware.GrantAdminUmum)
	semester.DELETE("/:id", handler.DeleteSemesterHandler, customMiddleware.GrantAdminUmum)

	akun := v1.Group("/akun")
	akun.POST("/login", handler.LoginHandler)
	akun.PATCH("/password/change", handler.ChangePasswordHandler, customMiddleware.Authentication)
	akun.PATCH("/password/reset/:id", handler.ResetPasswordHandler, customMiddleware.Authentication, customMiddleware.GrantAdminUmum)

	admin := v1.Group("/admin", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	admin.GET("", handler.GetAllAdminHandler)
	admin.GET("/:id", handler.GetAdminByIdHandler)
	admin.POST("", handler.InsertAdminHandler)
	admin.PUT("/:id", handler.EditAdminHandler)
	admin.DELETE("/:id", handler.DeleteAdminHandler)

	rektor := v1.Group("/rektor", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	rektor.GET("", handler.GetAllRektorHandler)
	rektor.GET("/:id", handler.GetRektorByIdHandler)
	rektor.POST("", handler.InsertRektorHandler)
	rektor.PUT("/:id", handler.EditRektorHandler)
	rektor.DELETE("/:id", handler.DeleteRektorHandler)

	operator := v1.Group("/operator", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	operator.GET("", handler.GetAllOperatorHandler)
	operator.GET("/:id", handler.GetOperatorByIdHandler)
	operator.POST("", handler.InsertOperatorHandler)
	operator.PUT("/:id", handler.EditOperatorHandler)
	operator.DELETE("/:id", handler.DeleteOperatorHandler)

	mahasiswa := v1.Group("/mahasiswa", customMiddleware.Authentication)
	mahasiswa.GET("", handler.GetAllMahasiswaHandler)
	mahasiswa.GET("/:id", handler.GetMahasiswaByIdHandler)
	mahasiswa.POST("", handler.InsertMahasiswaHandler, customMiddleware.GrantAdminUmum)
	mahasiswa.PUT("/:id", handler.EditMahasiswaHandler, customMiddleware.GrantAdminUmum)
	mahasiswa.DELETE("/:id", handler.DeleteMahasiswaHandler, customMiddleware.GrantAdminUmum)

	dosen := v1.Group("/dosen", customMiddleware.Authentication)
	dosen.GET("", handler.GetAllDosenHandler)
	dosen.GET("/:id", handler.GetDosenByIdHandler)
	dosen.POST("", handler.InsertDosenHandler, customMiddleware.GrantAdminUmum)
	dosen.PUT("/:id", handler.EditDosenHandler, customMiddleware.GrantAdminUmum)
	dosen.DELETE("/:id", handler.DeleteDosenHandler, customMiddleware.GrantAdminUmum)

	prestasi := v1.Group("/prestasi", customMiddleware.Authentication)
	prestasi.GET("", handler.GetAllPrestasiHandler)
	prestasi.GET("/:id", handler.GetPrestasiByIdHandler)
	prestasi.POST("", handler.InsertPrestasiHandler, customMiddleware.GrantMahasiswa)
	prestasi.PUT("/:id", handler.EditPrestasiHandler)
	prestasi.DELETE("/:id", handler.DeletePrestasiHandler)
	prestasi.PATCH("/:id/sertifikat", handler.EditSertifikatPrestasiHandler)

	kategoriProgramKM := v1.Group("/kategori-program", customMiddleware.Authentication)
	kategoriProgramKM.GET("", handler.GetAllKategoriProgramProgramKMHandler)
	kategoriProgramKM.POST("", handler.InsertKategoriProgramKMHandler, customMiddleware.GrantAdminIKU2)
	kategoriProgramKM.PUT("/:id", handler.EditKategoriProgramKMHandler, customMiddleware.GrantAdminIKU2)
	kategoriProgramKM.DELETE("/:id", handler.DeleteKategoriProgramKMHandler, customMiddleware.GrantAdminIKU2)

	kampusMerdeka := v1.Group("/kampus-merdeka", customMiddleware.Authentication)
	kampusMerdeka.GET("", handler.GetAllKMHandler)
	kampusMerdeka.GET("/:id", handler.GetKMByIdHandler)
	kampusMerdeka.POST("", handler.InsertKMHandler, customMiddleware.GrantMahasiswa)
	kampusMerdeka.PUT("/:id", handler.EditKMHandler, customMiddleware.GrantAdminIKU2OperatorAndMahasiswa)
	kampusMerdeka.DELETE("/:id", handler.DeleteKMHandler, customMiddleware.GrantAdminIKU2OperatorAndMahasiswa)
	kampusMerdeka.PATCH("/:id/surat-tugas", handler.EditSuratTugasHandler, customMiddleware.GrantAdminIKU2OperatorAndMahasiswa)
	kampusMerdeka.PATCH("/:id/berita-acara", handler.EditBeritaAcaraHandler, customMiddleware.GrantAdminIKU2OperatorAndMahasiswa)

	dashboard := v1.Group("/dashboard")
	dashboard.GET("/kampus-merdeka/kategori", handler.GetKMDashboardByKategoriHandler)
	dashboard.GET("/kampus-merdeka/fakultas", handler.GetKMDashboardByFakultasHandler)
	dashboard.GET("/:fitur/detail", handler.GetDetailDashboardHandler)
	dashboard.GET("/prestasi/tingkat", handler.GetPrestasiDashboardByTingkatHandler)
	dashboard.GET("/prestasi/fakultas", handler.GetPrestasiDashboardByFakultasHandler)
	dashboard.GET("/total", handler.GetTotalDashboardHandler)

	return app
}
