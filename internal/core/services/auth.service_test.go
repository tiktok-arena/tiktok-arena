package services

import (
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"tiktok-arena/internal/api/controllers"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/data/repository"
)

type AuthSuite struct {
	suite.Suite
	app        *fiber.App
	w          *httptest.ResponseRecorder
	controller *controllers.AuthController
	mock       sqlmock.Sqlmock
	db         *sql.DB
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func (as *AuthSuite) TearDownTest() {
	as.db.Close()
}

func (as *AuthSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.Nil(as.T(), err)
	as.db = db

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 as.db,
		PreferSimpleProtocol: true,
	})

	database, err := gorm.Open(dialector)
	assert.Nil(as.T(), err)
	userRepository := repository.NewUserRepository(database)
	authService := NewAuthService(userRepository)
	as.controller = controllers.NewAuthController(authService)
	app := fiber.New(fiber.Config{})
	as.app = app
	as.w = httptest.NewRecorder()
	as.mock = mock
}

func (as *AuthSuite) TestNewUser() {
	newUser := dtos.AuthInput{Name: "test", Password: "test"}
	rows := sqlmock.NewRows([]string{"id"}).AddRow(nil)
	as.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "users" WHERE name = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(newUser.Name).
		WillReturnRows(rows)
	as.mock.ExpectBegin()
	id, _ := uuid.NewUUID()
	rows = sqlmock.NewRows([]string{"id", "name", "password"}).AddRow(id, newUser.Name, newUser.Password)
	as.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("photo_url","name","password") VALUES ($1,$2,$3) RETURNING "id","name","password"`)).
		WithArgs(sqlmock.AnyArg(), newUser.Name, sqlmock.AnyArg()).
		WillReturnRows(rows)
	as.mock.ExpectCommit()
	as.app.Post("/register", as.controller.RegisterUser)

	body, err := json.Marshal(newUser)
	if err != nil {
		assert.Error(as.T(), err)
	}
	reader := strings.NewReader(string(body))
	req := httptest.NewRequest("POST", "http://localhost:8000/register", reader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := as.app.Test(req)
	if err != nil {
		assert.Error(as.T(), err)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Equal(as.T(), resp.StatusCode, fiber.StatusCreated)
	assert.Contains(as.T(), string(bodyBytes), newUser.Name)
	assert.Contains(as.T(), string(bodyBytes), id.String())
	assert.Nil(as.T(), err)
}

func (as *AuthSuite) TestGetUsernameAndPassword() {
	newUser := dtos.AuthInput{Name: "test", Password: "test"}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		assert.Error(as.T(), err)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		assert.Error(as.T(), err)
	}
	rows := sqlmock.NewRows([]string{"id", "name", "password"}).AddRow(id, newUser.Name, string(hashedPassword))
	as.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE name = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(newUser.Name).
		WillReturnRows(rows)

	as.app.Post("/login", as.controller.LoginUser)

	body, err := json.Marshal(newUser)
	if err != nil {
		assert.Error(as.T(), err)
	}
	reader := strings.NewReader(string(body))
	req := httptest.NewRequest("POST", "http://localhost:8000/login", reader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := as.app.Test(req)
	if err != nil {
		assert.Error(as.T(), err)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Equal(as.T(), resp.StatusCode, fiber.StatusOK)
	assert.Contains(as.T(), string(bodyBytes), newUser.Name)
	assert.Nil(as.T(), err)
}

// func (as *AuthSuite) TestWhoAmI() {
// 	newUser := dtos.WhoAmI{Name: "test"}
// 	id, err := uuid.NewUUID()
// 	if err != nil {
// 		assert.Error(as.T(), err)
// 	}
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
// 	if err != nil {
// 		assert.Error(as.T(), err)
// 	}
// 	rows := sqlmock.NewRows([]string{"id", "name", "password"}).AddRow(id, newUser.Name, string(hashedPassword))
// 	as.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "users" WHERE name = $1 ORDER BY "users"."id" LIMIT 1`)).
// 		WithArgs(newUser.Name).
// 		WillReturnRows(rows)
// 	as.mock.ExpectBegin()
// 	rows = sqlmock.NewRows([]string{"id", "name", "password"}).AddRow(id, newUser.Name, string(hashedPassword))
// 	as.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "photo_url" FROM "users" WHERE "users"."id" = $1`)).
// 		WithArgs(id).
// 		WillReturnRows(rows)
// 	as.mock.ExpectCommit()
//
// 	as.app.Post("/whoami", as.controller.WhoAmI)
//
// 	body, err := json.Marshal(newUser)
// 	if err != nil {
// 		assert.Error(as.T(), err)
// 	}
// 	reader := strings.NewReader(string(body))
// 	req := httptest.NewRequest("POST", "http://localhost:8000/whoami", reader)
// 	req.Header.Set("Content-Type", "application/json")
//
// 	resp, err := as.app.Test(req)
// 	if err != nil {
// 		assert.Error(as.T(), err)
// 	}
// 	bodyBytes, _ := io.ReadAll(resp.Body)
// 	assert.Equal(as.T(), resp.StatusCode, fiber.StatusOK)
// 	assert.Contains(as.T(), string(bodyBytes), newUser.Name)
// 	assert.Nil(as.T(), err)
// }
