package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"fmt"
)

type (
	child struct {
		ID    int      `json:"id"`
		Child []*child `json:"children"`
	}
)

func Show(c *gin.Context) {
	c.HTML(http.StatusOK, "goadmin/layout/index", nil)
}

func Post(c *gin.Context) {
	var childs []*child
	order := c.PostForm("_order")
	json.Unmarshal([]byte(order), &childs)

	parse(childs,0)

}

func Edit(c *gin.Context) {

}

func Dele(c *gin.Context) {

}

func parse(childs []*child, PID int) {
	for _, v := range childs {
		//Create
		fmt.Println(PID,v.ID, v.Child)

		if len(v.Child) > 0 {
			//解析Child
			parse(v.Child, v.ID)
		}
	}
}
