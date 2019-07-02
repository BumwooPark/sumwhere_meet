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

// Create Profile
// @Summary 프로필 생성
// @Description create user profile
// @ID get-string-by-int
// @Tags profile
// @Accept  json
// @Produce  json
// @Param profile body profile.Profile true "profile"
// @Success 201 {object} profile.Profile
// @Header 200 {string} Token "qwerty"
// @Router /profile [post]
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

// Update Profile
// @Summary 프로필 업데이트
// @Description update user profile
// @ID get-string-by-int
// @Tags profile
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} profile.Profile
// @Header 200 {string} Token "qwerty"
// @Router /profile/{id} [patch]
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

// Get Profile
// @Summary 프로필 가져오기
// @Description get user profile
// @Tags profile
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} profile.Profile
// @Header 200 {string} Token "qwerty"
// @Router /profile/{id} [get]
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

// Set profile Image
// @Summary 프로필 이미지 추가
// @Description set user profile image
// @ID profile
// @Tags profile
// @Accept multipart/form-data
// @Produce  json
// @Param id path int true "User ID"
// @Param file formData file true "image must 3"
// @Failure 400 {object} utils.ApiError
// @Success 200 {boolean} true
// @Header 200 {string} Token "qwerty"
// @Router /profile/image/{id} [post]
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

// Get City Models
// @Summary 전국 시 가져오기
// @Description 전국 시 모델 가져오기
// @Tags profile
// @Accept json
// @Produce  json
// @Success 200 {array} profile.Area
// @Header 200 {string} Token "qwerty"
// @Router /profile/city [get]
func (Controller) GetCity(e echo.Context) error {
	area, err := Area{}.GetCity(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, area)
}

// Get City Models
// @Summary 전국 시 모델 가져오기
// @Description 전국 시 모델 가져오기
// @Accept json
// @Tags profile
// @Produce  json
// @Param city query string true "원하는 도시"
// @Success 200 {array} profile.Area
// @Header 200 {string} Token "qwerty"
// @Router /profile/district [get]
func (Controller) GetDistrict(e echo.Context) error {
	area, err := Area{}.GetDistrict(e.Request().Context(), e.QueryParam("city"))
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusOK, utils.ApiErrorDB, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, area)
}
