package route

import (
	"be-2/src/api/handler"

	customMiddleware "be-2/src/api/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer() *echo.Echo {
	app := echo.New()
	app.Use(middleware.CORS())

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

	return app
}
