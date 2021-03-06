package utils

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/gin-contrib/sessions"
)

func Toastr(c *gin.Context) *toastr {
	return &toastr{context: c}
}

type toastr struct {
	context *gin.Context
	typ     string
	message string
}

func (t *toastr) Success(message string) {
	js := fmt.Sprintf("toastr.%s('%s');", "success", message)
	session := sessions.Default(t.context)
	session.AddFlash(js, "toastr")
	session.Save()
}

func (t *toastr) Info(message string) {
	js := fmt.Sprintf("toastr.%s('%s');", "info", message)
	session := sessions.Default(t.context)
	session.AddFlash(js, "toastr")
	session.Save()
}

func (t *toastr) Warning(message string) {
	js := fmt.Sprintf("toastr.%s('%s');", "warning", message)
	session := sessions.Default(t.context)
	session.AddFlash(js, "toastr")
	session.Save()
}

func (t *toastr) Error(message string) {
	js := fmt.Sprintf("toastr.%s('%s');", "error", message)
	session := sessions.Default(t.context)
	session.AddFlash(js, "toastr")
	session.Save()
}
