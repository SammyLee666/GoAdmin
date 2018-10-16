package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"goadmin/goadmin/models"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"github.com/gin-contrib/sessions"
	"log"
)

func Show(c *gin.Context) {
	var errMsg models.ErrMsg
	session := sessions.Default(c)
	errors := session.Flashes("errors")
	session.Save()
	if len(errors) > 0 {
		errMsg = errors[0].(models.ErrMsg)
	}

	c.HTML(http.StatusOK, "goadmin/layout/index", gin.H{
		"errors": errMsg,
	})
}

func Post(c *gin.Context) {
	order := c.PostForm("_order")
	if len(order) == 0 {
		//初次运行创建表
		//db.Mysql.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&models.Menu{})

		var menu models.Menu
		err := c.ShouldBind(&menu);
		if err != nil {
			session := sessions.Default(c)
			ve := err.(validator.ValidationErrors)
			errMsg := models.ErrMsg{}
			for _, e := range ve {
				if e.Field == "Pid" {
					switch e.Tag {
					case "min":
						errMsg["parent_id"] = "父级菜单ID出错"
					case "max":
						errMsg["parent_id"] = "父级菜单ID出错"
					default:
						errMsg["parent_id"] = err.Error()
					}
				}
				if e.Field == "Title" {
					switch e.Tag {
					case "min":
						errMsg["title"] = "标题必须填写"
					case "max":
						errMsg["title"] = "标题长度不超过50"
					default:
						errMsg["title"] = err.Error()
					}
				}
				if e.Field == "Icon" {
					switch e.Tag {
					case "min":
						errMsg["icon"] = "图标必须选择"
					case "max":
						errMsg["icon"] = "选择图标出现错误"
					default:
						errMsg["icon"] = err.Error()
					}
				}
				if e.Field == "Uri" {
					switch e.Tag {
					case "min":
						errMsg["uri"] = "路径必须填写"
					case "max":
						errMsg["uri"] = "路径长度不超过50"
					default:
						errMsg["uri"] = err.Error()
					}
				}
			}
			log.Println(errMsg)
			session.AddFlash(errMsg, "errors")
			if err := session.Save(); err != nil {
				c.String(400, err.Error())
				return
			}
			c.Redirect(http.StatusFound, c.Request.Header.Get("Referer"))
			return
		}
	} else {
		changeParentId(order)
	}

}

func Edit(c *gin.Context) {

}

func Dele(c *gin.Context) {

}

func changeParentId(order string) {
	var childs []*models.Child
	json.Unmarshal([]byte(order), &childs)
	parse(childs, 0)
}

func parse(childs []*models.Child, PID int) {
	for _, v := range childs {
		//Create
		fmt.Println(PID, v.ID, v.Child)

		if len(v.Child) > 0 {
			//解析Child
			parse(v.Child, v.ID)
		}
	}
}
