package middleware

import (
	"be-2/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GrantAdminUmum(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.ADMIN) ||
			claims["bagian"].(string) != "UMUM" {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantMahasiswa(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.MAHASISWA) {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}
