package profile

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/testfixtures.v2"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sumwhere_meet/database"
	"sumwhere_meet/middlewares"
	"sumwhere_meet/utils"
	"testing"
)

var (
	server   *echo.Echo
	handler  func(handlerFunc echo.HandlerFunc, c echo.Context) error
	db       *xorm.Engine
	testdb   *sql.DB
	fixtures *testfixtures.Context
)

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}

func TestMain(m *testing.M) {

	server = echo.New()
	os.Setenv("DATABASE_NAME", "test")
	db, err := database.NewDatabase()
	if err != nil {
		fmt.Println(err)
	}
	databasehandler := middlewares.ContextDB("database", db)
	_ = db.Sync2(new(Image), new(Profile), new(Area), new(InterestCategory), new(InterestDetail))
	handler = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return databasehandler(handlerFunc)(c)
	}
	server.Validator = &Validator{}

	testdb, err := sql.Open("mysql", "root:1q2w3e4r@tcp(localhost:3306)/test?charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}

	fixtures, err = testfixtures.NewFolder(testdb, &testfixtures.MySQL{}, "fixtures")
	if err != nil {
		log.Fatal(err)
	}
	testfixtures.SkipDatabaseNameCheck(true)

	if err := fixtures.DetectTestDatabase(); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		fmt.Println(err)
	}
}

func TestController_Create(t *testing.T) {
	signupProfile := `{"userID":1,"age":20,"phone":"01051416906","job":"강사"}`
	req := httptest.NewRequest(echo.POST, "/profile", strings.NewReader(signupProfile))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	require.NoError(t, handler(Controller{}.Create, ctx))
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body)

	var v struct {
		Result  Profile        `json:"result"`
		Success bool           `json:"success"`
		Error   utils.ApiError `json:"error"`
	}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	assert.Equal(t, true, v.Success, v.Error)
}

func TestController_Update(t *testing.T) {
	updateProfile := `{"age":10,"job":"하이"}`
	req := httptest.NewRequest(echo.PATCH, "/profile/1", strings.NewReader(updateProfile))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	require.NoError(t, handler(Controller{}.Update, ctx))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body)
}

func TestController_Exist(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/profile/1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	require.NoError(t, handler(Controller{}.Get, ctx))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body)
}

func TestArea(t *testing.T) {
	prepareTestDatabase()
}
