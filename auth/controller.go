package auth

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"sumwhere_meet/factory"
	"sumwhere_meet/utils"
)

type Controller struct {
}

func (c Controller) Init(g *echo.Group) {
	g.POST("/login", c.Login)
	g.POST("/signup/email", c.SignUp)
	g.POST("/login/facebook", c.FaceBook)
	g.POST("/login/kakao", c.Kakao)
	g.POST("/mobile", c.MobileAuth)
}

// login
// @Summary default 로그인
// @Description 로그인
// @ID get-string-by-int
// @Param email body auth.User true "email,password만 사용"
// @Tags signin&signup
// @Accept  json
// @Produce  json
// @Success 200 {object} auth.Token
// @Router /login [post]
func (c Controller) Login(e echo.Context) error {
	var v struct {
		Email    string `json:"email" valid:"required"`
		Password string `json:"password" valid:"required"`
	}

	if err := e.Bind(&v); err != nil {
		return utils.ReturnApiFail(e, http.StatusUnauthorized, utils.ApiErrorMissParameter, err)
	}

	if err := e.Validate(&v); err != nil {
		return utils.ReturnApiFail(e, http.StatusUnauthorized, utils.ApiErrorIllegalRequest, err)
	}

	u := User{
		Email:    v.Email,
		Password: v.Password,
	}

	user, err := u.Search(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotFound, utils.ApiErrorNotFound, err)
	}

	if !utils.ComparePasswords(u.Password, []byte(u.Password)) {
		return utils.ReturnApiFail(e, http.StatusUnauthorized, utils.ApiErrorPassword, err)
	}

	t, err := user.GenerateToken()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorTokenInvaild, err)
	}
	return utils.ReturnApiSucc(e, http.StatusCreated, map[string]string{"token": t})
}

// signup
// @Summary 회원가입
// @Description 로그인
// @Param email body auth.User true "유저모델"
// @Tags signin&signup
// @Accept  json
// @Produce  json
// @Success 200 {object} auth.Token
// @Router /signup/email [post]
func (c Controller) SignUp(e echo.Context) error {

	var u User
	err := e.Bind(&u)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	err = e.Validate(&u)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	u.Password = utils.HashAndSalt([]byte(u.Password))

	err = u.Insert(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	t, err := u.GenerateToken()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusInternalServerError, utils.ApiErrorSystem, err)
	}

	return utils.ReturnApiSucc(e, http.StatusCreated, map[string]string{"token": t})
}

// signup
// @Summary 페이스북으로 로그인
// @Description 페이스북으로 로그인
// @Param token body auth.Token true "페이스북 토큰"
// @Tags signin&signup
// @Accept  json
// @Produce  json
// @Success 200 {object} auth.Token
// @Router /login/facebook [post]
func (c Controller) FaceBook(e echo.Context) error {
	var v Token
	err := e.Bind(&v)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	url := "https://graph.facebook.com/v3.3/me?fields=id,email,name&access_token="
	res, err := http.Get(url + v.Token)
	defer res.Body.Close()
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	var f FaceBookUser
	err = json.Unmarshal(data, &f)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := e.Validate(f); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	return c.CreateOrSearch(e, &f, v.Token)

}

// login by kakao
// @Summary 카카오 로그인
// @Description 카카오로 로그인
// @Param token body auth.Token true "카카오 토큰"
// @Tags signin&signup
// @Accept  json
// @Produce  json
// @Success 200 {object} auth.Token
// @Router /login/kakao [post]
func (c Controller) Kakao(e echo.Context) error {

	var v Token
	err := e.Bind(&v)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorKakaoAuth, err)
	}
	req.Header.Add("Authorization", "bearer "+v.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorKakaoAuth, err)
	}

	defer res.Body.Close()
	user, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusNotAcceptable, utils.ApiErrorKakaoAuth, err)
	}

	var k KakaoUser
	err = json.Unmarshal(user, &k)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	//TODO: email이 없을 경우 어떻게 진행할것인가.
	if !k.KakaoAccount.HasEmail {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, errors.New("No have kakao email"))
	}

	return c.CreateOrSearch(e, &k, v.Token)
}

func (Controller) CreateOrSearch(e echo.Context, s SocialAction, token string) error {
	generater := func(u *User) error {
		t, err := u.GenerateToken()
		if err != nil {
			return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorSystem, err)
		}
		return utils.ReturnApiSucc(e, http.StatusCreated, map[string]string{"token": t})
	}

	group := s.IsExist(e.Request().Context())
	if len(group) == 1 {
		return generater(&group[0].User)
	}

	u := s.ToUser()
	err := u.Insert(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	social := s.ToSocial(u.ID, token)
	err = social.Create(e.Request().Context())
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}
	return generater(u)
}

func (Controller) MobileAuth(e echo.Context) error {
	var v struct {
		Phone string `json:"phone" valid:"numeric,required"`
	}

	if err := e.Bind(&v); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	if err := e.Validate(&v); err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorParameter, err)
	}

	factory.Firebase(e.Request().Context())

	return nil
}
