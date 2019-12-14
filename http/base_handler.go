package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type baseHandler struct {
}

func (handler *baseHandler) BadRequestError(ctx *gin.Context, message string) {
	handler.respond(ctx, http.StatusBadRequest, message)
}

func (handler *baseHandler) ForbiddenError(ctx *gin.Context, message string)  {
	handler.respond(ctx, http.StatusForbidden, message)
}

func (handler *baseHandler) InternalServerError(ctx *gin.Context, message string)  {
	handler.respond(ctx, http.StatusInternalServerError, message)
}

func (handler *baseHandler) respond(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, &Response{Error: true, Message: message})
}