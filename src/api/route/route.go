package route

import (
	"be-2/src/api/handler"

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

	fakultas := v1.Group("/fakultas")
	fakultas.GET("", handler.GetAllFakultasHandler)
	fakultas.GET("/:id", handler.GetFakultasByIdHandler)
	fakultas.POST("", handler.InsertFakultasHandler)
	fakultas.PUT("/:id", handler.EditFakultasHandler)
	fakultas.DELETE("/:id", handler.DeleteFakultasHandler)

	prodi := v1.Group("/prodi")
	prodi.GET("", handler.GetAllProdiHandler)
	prodi.GET("/:id", handler.GetProdiByIdHandler)
	prodi.POST("", handler.InsertProdiHandler)
	prodi.PUT("/:id", handler.EditProdiHandler)
	prodi.DELETE("/:id", handler.DeleteProdiHandler)

	semester := v1.Group("/semester")
	semester.GET("", handler.GetAllSemesterHandler)
	semester.POST("", handler.InsertSemesterHandler)
	semester.PUT("/:id", handler.EditSemesterHandler)
	semester.DELETE("/:id", handler.DeleteSemesterHandler)

	return app
}
