package authHandler_test

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todo/internal/entity"
	authHandler "todo/internal/handler/auth"
	"todo/internal/service"
	mock_service "todo/internal/service/mocks"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerLogin(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, user entity.UserLoginReq)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           entity.UserLoginReq
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "test@mail.ru", "password": "qwerty"}`,
			inputUser: entity.UserLoginReq{
				Email:    "test@mail.ru",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuth, user entity.UserLoginReq) {
				s.EXPECT().Login(&user).Return(entity.TokensRes{
					Access:  "tokenaccess",
					Refresh: "tokenrefresh",
				}, nil, http.StatusOK)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"access":"tokenaccess","refresh":"tokenrefresh"}`,
		},
		{
			name:      "Invalid fields",
			inputBody: `{"email": "test.ru", "password": "qwerty"}`,
			mockBehavior: func(s *mock_service.MockAuth, user entity.UserLoginReq) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "\"Key: 'UserLoginReq.Email' Error:Field validation for 'Email' failed on the 'email' tag\"",
		},
		{
			name:      "Failure",
			inputBody: `{"email": "test@mail.ru", "password": "qwerty"}`,
			inputUser: entity.UserLoginReq{
				Email:    "test@mail.ru",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuth, user entity.UserLoginReq) {
				s.EXPECT().Login(&user).Return(entity.TokensRes{}, errors.New("login error"), http.StatusInternalServerError)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"login error"}`,
		},

		{
			name:      "User not exists",
			inputBody: `{"email": "test@mail.ru", "password": "qwerty"}`,
			inputUser: entity.UserLoginReq{
				Email:    "test@mail.ru",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuth, user entity.UserLoginReq) {
				s.EXPECT().Login(&user).Return(entity.TokensRes{}, errors.New("user with email not exists"), http.StatusBadRequest)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"user with email not exists"}`,
		},
		{
			name:      "Not valid login or passward",
			inputBody: `{"email": "test@mail.ru", "password": "qwerty"}`,
			inputUser: entity.UserLoginReq{
				Email:    "test@mail.ru",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuth, user entity.UserLoginReq) {
				s.EXPECT().Login(&user).Return(entity.TokensRes{}, errors.New("not valid login or password"), http.StatusBadRequest)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not valid login or password"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuth(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			//log
			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			//valdiator
			validator := validator.New()
			services := &service.Service{Auth: auth}
			hand := authHandler.NewAuthHandler(log, services, validator)

			router := http.NewServeMux()
			router.HandleFunc("POST /auth/login", hand.Login())

			//test Req
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(testCase.inputBody))

			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody+"\n", w.Body.String())
		})
	}
}
