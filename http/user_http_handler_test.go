package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/adigunhammedolalekan/sms-forwarder/mocks"
	"github.com/adigunhammedolalekan/sms-forwarder/types"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHttpHandler_CreateUserHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUser := &types.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		Email:    "user@test.io",
		Password: "password",
		Token:    "token",
	}
	userStore := mocks.NewMockUserStore(controller)
	userStore.EXPECT().CreateUser(mockUser.Email, mockUser.Password).Return(mockUser, nil)

	handler := NewUserHttpHandler(userStore)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/user/new", bytes.NewBufferString(`{"email": "user@test.io", "password": "password"}`))

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.CreateUserHandler(ctx)
	if w.Code != http.StatusOK {
		t.Fatalf("expected httpOK. Got %d", w.Code)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	var response struct {
		Error bool `json:"error"`
		Message string `json:"message"`
		Data struct{
			Email string `json:"email"`
			Password string `json:"password"`
			Token string `json:"token"`
		}
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}
	if response.Error {
		t.Fatal("not expecting error")
	}
	if response.Data.Token == "" {
		t.Fatal("empty token returned")
	}
}

func TestUserHttpHandler_CreateUserHandlerBadRequest(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userStore := mocks.NewMockUserStore(controller)
	handler := NewUserHttpHandler(userStore)

	in := &bytes.Buffer{}
	req := httptest.NewRequest("POST", "/", in)
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	handler.CreateUserHandler(ctx)

	if want, got := http.StatusBadRequest, w.Code; want != got {
		t.Fatalf("error: want code %d; got %d", want, got)
	}
}

func TestUserHttpHandler_CreateUserHandlerInternalServerError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUser := &types.User{
		Email:    "user@test.io",
		Password: "password",
	}

	userStore := mocks.NewMockUserStore(controller)
	userStore.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("connection closed"))
	handler := NewUserHttpHandler(userStore)

	in := &bytes.Buffer{}
	json.NewEncoder(in).Encode(mockUser)
	req := httptest.NewRequest("POST", "/", in)
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	handler.CreateUserHandler(ctx)

	if want, got := http.StatusInternalServerError, w.Code; want != got {
		t.Fatalf("error: want code %d; got %d", want, got)
	}
}

func TestUserHttpHandler_AuthenticateUserHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	type mockUserRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	user := &mockUserRequest{
		Email:    "user@test.co",
		Password: "password",
	}
	mockUser := &types.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		Email:    "user@test.co",
		Password: "password",
		Token:    "token",
	}
	in := &bytes.Buffer{}
	json.NewEncoder(in).Encode(user)

	userStore := mocks.NewMockUserStore(controller)
	userStore.EXPECT().AuthenticateUser(user.Email, user.Password).Return(mockUser, nil)

	handler := NewUserHttpHandler(userStore)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/user/authenticate", in)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.AuthenticateUserHandler(ctx)
	if want, got := http.StatusOK, w.Code; want != got {
		t.Fatalf("error: want %d, got %d", want, got)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	var response struct {
		Error bool `json:"error"`
		Message string `json:"message"`
		Data struct{
			Email string `json:"email"`
			Password string `json:"password"`
			Token string `json:"token"`
		}
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}
	if response.Error {
		t.Fatal("not expecting error")
	}
	if response.Data.Token == "" {
		t.Fatal("empty token returned")
	}
}

func TestUserHttpHandler_AuthenticateUserHandlerInternalServerError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUser := &types.User{
		Email:    "user@test.co",
		Password: "password",
	}

	userStore := mocks.NewMockUserStore(controller)
	userStore.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("connection closed"))
	handler := NewUserHttpHandler(userStore)

	in := &bytes.Buffer{}
	json.NewEncoder(in).Encode(mockUser)
	req := httptest.NewRequest("POST", "/", in)
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	handler.AuthenticateUserHandler(ctx)

	if want, got := http.StatusInternalServerError, w.Code; want != got {
		t.Fatalf("error: want code %d; got %d", want, got)
	}
}