package http

import (
	s "github.com/adigunhammedolalekan/sms-forwarder/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHttpHandler struct {
	store s.UserStore
	*baseHandler
}

func NewUserHttpHandler(store s.UserStore) *UserHttpHandler {
	return &UserHttpHandler{store: store}
}

func (handler *UserHttpHandler) CreateUserHandler(ctx *gin.Context) {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		handler.BadRequestError(ctx, "bad request: malformed json data")
		return
	}

	user, err := handler.store.CreateUser(req.Email, req.Password)
	if err != nil {
		handler.InternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &Response{Error: false, Message: "success", Data: user})
}

func (handler *UserHttpHandler) AuthenticateUserHandler(ctx *gin.Context) {
	type request struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		handler.BadRequestError(ctx, "bad request: malformed json data")
		return
	}

	user, err := handler.store.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		handler.InternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &Response{Error: false, Message: "success", Data: user})
}