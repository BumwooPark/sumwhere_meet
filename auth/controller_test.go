package auth

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"showper_server/middlewares"
	"showper_server/utils"
	"strings"
	"testing"
)

var (
	server  *echo.Echo
	handler func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

func init() {
	server = echo.New()

	fb, err := middlewares.NewFireBaseApp()
	if err != nil {
		panic(err)
	}

	fbmi := middlewares.ContextFireBase(fb)

	handler = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return fbmi(handlerFunc)(c)
	}
	server.Validator = &Validator{}
}

func TestController_SignUp(t *testing.T) {
	signupJson := `{"email":"qbne9dfhdfg2@naver.com","nickname":"hadsgd","password":"1q2w3e4r"}`
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(signupJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	assert.NoError(t, handler(Controller{}.SignUp, ctx))

	var v struct {
		Result  User           `json:"result"`
		Success bool           `json:"success"`
		Error   utils.ApiError `json:"error"`
	}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	assert.Equal(t, true, v.Success, v.Error)
}

func TestController_MobileAuth(t *testing.T) {
	number := `{"phone":"01051416906"}`
	req := httptest.NewRequest(echo.POST, "/mobile", strings.NewReader(number))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	require.NoError(t, handler(Controller{}.MobileAuth, ctx))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body)

}
