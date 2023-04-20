package handler

import (
	"be-2/src/api/request"
	"be-2/src/api/response"
	"be-2/src/config/database"
	"be-2/src/model"
	"be-2/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllKategoriProgramProgramKMHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.KategoriProgramKm{}

	if err := db.WithContext(ctx).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertKategoriProgramKMHandler(c echo.Context) error {
	req := &request.KategoriProgramKm{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).Create(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func EditKategoriProgramKMHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	req := &request.KategoriProgramKm{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).Where("id", id).Updates(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteKategoriProgramKMHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.KategoriProgramKm), "id", id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
