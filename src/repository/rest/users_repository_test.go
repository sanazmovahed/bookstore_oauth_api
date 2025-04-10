package rest

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

var mock bool

func TestMain(m *testing.M) {
	// fmt.Println("about to start test cases...")
	// rest.StartMockupServer()
	// os.Exit(m.Run())
	//  Register custom flags
	flag.BoolVar(&mock, "mock", false, "Use mock server")

	// Important: Parse all flags (including test-internal ones like -test.testlogfile)
	flag.Parse()

	if mock {
		fmt.Println("Starting mock server...")
		// Initialize your mock server or mock setup here
	}

	fmt.Println("about to start test cases...")
	os.Exit(m.Run())
}

func TestLoginUserFromTimeoutAPI(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore.com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com", "password":"password_"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password_")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restclient response when trying to login user", err.Status)

}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore.com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com", "password":"password_"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credential", "status":404, "error":"not_found"}`,

		// RespBody:     `{"message":"invalid login credential", "status":"404", "error":"not_found"}`, //on-purpose fault in type of status
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password_")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restclient response when trying to login user", err.Status)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore.com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com", "password":"password_"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credential", "status":404, "error":"not_found"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password_")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.NotFound, err.Status)
	assert.EqualValues(t, "invalid login credential", err.Status)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore.com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com", "password":"password_"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "10","first_name": "ilkin","last_name": "salmani","email": "i.salmani95@gmail.com","date_created": "2025-04-06 05:38:31","status": "active"}`,

		//on-purpose fault in type of id int to string
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password_")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal users login response", err.Status)
}
func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore.com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com", "password":"password_"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 10,"first_name": "ilkin","last_name": "salmani","email": "i.salmani95@gmail.com","date_created": "2025-04-06 05:38:31","status": "active"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password_")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 10, user.Id)
	assert.EqualValues(t, "ilkin", user.FirstName)
	assert.EqualValues(t, "salmani", user.LastName)
	assert.EqualValues(t, "i.salmani95@gmail.com", user.Email)

}
