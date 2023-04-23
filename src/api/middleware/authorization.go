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
			claims["bagian"].(string) != util.UMUM {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU2(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.ADMIN) ||
			claims["bagian"].(string) != util.IKU2 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantMahasiswa(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.MAHASISWA) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU2OperatorAndMahasiswa(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		role := claims["role"].(string)
		bagian := claims["bagian"].(string)
		if role != string(util.MAHASISWA) && role != string(util.ADMIN) &&
			role != string(util.OPERATOR) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		if role == string(util.ADMIN) && bagian != util.IKU2 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU2OperatorAndRektor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		role := claims["role"].(string)
		bagian := claims["bagian"].(string)
		if role != string(util.REKTOR) && role != string(util.ADMIN) &&
			role != string(util.OPERATOR) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		if role == string(util.ADMIN) && bagian != util.IKU2 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU2AndRektor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		role := claims["role"].(string)
		bagian := claims["bagian"].(string)
		if role != string(util.REKTOR) && role != string(util.ADMIN) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		if role == string(util.ADMIN) && bagian != util.IKU2 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}
