package profile

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"sumwhere_meet/utils"
)

type Controller struct{}

func (c Controller) Init(g *echo.Group) {
	g.POST("/profile", c.Create)
	g.POST("/profile/image", c.Image)
	g.PATCH("/profile/:id", c.Update)
	g.GET("/profile/:id", c.Get)
	g.GET("/profile/city", c.GetCity)
	g.GET("/profile/district", c.GetDistrict)
}

func (Controller) Create(e echo.Context) error {
	var p Profile
	if err := e.Bind(&p); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := e.Validate(p); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := p.Create(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusCreated, p)
}

func (Controller) Update(e echo.Context) error {

	id, err := strconv.ParseInt(e.Param("id"), 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotFound, utils.ApiErrorNotFound, err)
	}

	p := Profile{}
	if err := p.Get(e.Request().Context(), id); err != nil {
		return utils.ReturnApiFail(e, http.StatusOK, utils.ApiErrorUserNotExists, err)
	}

	var vo Profile

	if err := e.Bind(&vo); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	vo.UserID = id

	if err := e.Validate(&vo); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	p.Age = vo.Age
	p.Job = vo.Job

	if err := p.Update(e.Request().Context()); err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, p)
}

func (Controller) Get(e echo.Context) error {
	id, err := strconv.ParseInt(e.Param("id"), 10, 64)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotFound, utils.ApiErrorNotFound, err)
	}

	p := Profile{}
	if err := p.Get(e.Request().Context(), id); err != nil {
		return utils.ReturnApiFail(e, http.StatusOK, utils.ApiErrorUserNotExists, err)
	}
	return utils.ReturnApiSucc(e, http.StatusOK, p)
}

func (Controller) Image(e echo.Context) error {
	form, err := e.MultipartForm()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	size := len(form.File["image"])
	if size < 3 {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, errors.New("image count not enough"))
	}

	files := form.File["image"]

	err = Image{}.Save(e.Request().Context(), 1, files)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusCreated, true)
}

func (Controller) GetCity(e echo.Context) error {
	area, err := Area{}.GetCity(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, area)
}

func (Controller) GetDistrict(e echo.Context) error {
	area, err := Area{}.GetDistrict(e.Request().Context(), e.QueryParam("city"))
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusOK, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, area)
}
