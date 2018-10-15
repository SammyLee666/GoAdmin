package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"goadmin/goadmin/models"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
)

func Show(c *gin.Context) {
	c.HTML(http.StatusOK, "goadmin/layout/index", nil)
}

func Post(c *gin.Context) {
	order := c.PostForm("_order")
	if len(order) == 0 {
		//初次运行创建表
		//db.Mysql.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&models.Menu{})

		var menu models.Menu
		err := c.ShouldBind(&menu);
		if err != nil {
			ve := err.(validator.ValidationErrors)
			for _, e := range ve {
				fmt.Println(e.Field)
				if e.Field == "Pid" {
					switch e.Tag {
					case "required":
						c.JSON(http.StatusBadRequest, gin.H{"error": "父级菜单必选"})
						return
					default:
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
				}
				if e.Field == "Title" {
					switch e.Tag {
					case "min":
						c.JSON(http.StatusBadRequest, gin.H{"error": "标题必须填写"})
						return
					case "max":
						c.JSON(http.StatusBadRequest, gin.H{"error": "标题长度不超过50"})
						return
					default:
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
				}
				if e.Field == "Icon" {
					switch e.Tag {
					case "min":
						c.JSON(http.StatusBadRequest, gin.H{"error": "图标必须选择"})
						return
					case "max":
						c.JSON(http.StatusBadRequest, gin.H{"error": "选择图标出现错误"})
						return
					default:
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
				}
				if e.Field == "Uri" {
					switch e.Tag {
					case "min":
						c.JSON(http.StatusBadRequest, gin.H{"error": "路径必须填写"})
						return
					case "max":
						c.JSON(http.StatusBadRequest, gin.H{"error": "路径长度不超过50"})
						return
					default:
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func ListOfErrors(e error) []map[string]string {
	ve := e.(validator.ValidationErrors)
	InvalidFields := make([]map[string]string, 0)

	for _, e := range ve {
		errors := map[string]string{}
		// field := reflect.TypeOf(e.NameNamespace)
		errors[e.Name] = e.Tag
		InvalidFields = append(InvalidFields, errors)
	}

	return InvalidFields
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
