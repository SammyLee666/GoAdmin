package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserProfile(c *gin.Context)  {
	c.HTML(http.StatusOK, "goadmin/layout/index", nil)
}