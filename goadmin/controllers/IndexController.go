package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"goadmin/goadmin/models"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"github.com/gin-contrib/sessions"
	"net/url"
	"goadmin/db"
	"encoding/xml"
)

func Show(c *gin.Context) {
	var oldMsg url.Values
	var errMsg models.ErrMsg

	v := &models.TreeDom{Class: "dd-list"}
	v.List = append(v.List, models.LiDom{
		Class:  "dd-item",
		DataId: "1",
		Handle: models.DivHandleDom{

		},
	})
	//v.Svs = append(v.Svs, server{"Beijing_VPN", "127.0.0.2"})
	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// os.Stdout.Write([]byte(xml.Header))

	// os.Stdout.Write(output)
	//将字节流转换成string输出
	fmt.Println("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(output))

	session := sessions.Default(c)
	errors := session.Flashes("errors")
	oldForm := session.Flashes("oldForm")
	session.Save()

	if len(errors) > 0 {
		errMsg = errors[0].(models.ErrMsg)
	}
	if len(oldForm) > 0 {
		oldMsg = oldForm[0].(url.Values)
	}
	//TreeView
	models.TreeView()

	c.HTML(http.StatusOK, "goadmin/layout/index", gin.H{
		"_errors": errMsg,
		"_old":    oldMsg,
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
			oldForm := url.Values{}
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
			oldForm = c.Request.PostForm
			session.AddFlash(oldForm, "oldForm")
			session.AddFlash(errMsg, "errors")
			if err := session.Save(); err != nil {
				c.String(400, err.Error())
				return
			}
			c.Redirect(http.StatusFound, c.Request.Header.Get("Referer"))
			return
		}
		model := db.Mysql.Create(&menu)
		fmt.Println(model)
		fmt.Println(menu.ID)
	} else {
		changeParentId(order)
	}

}

func Edit(c *gin.Context) {
	fmt.Println(1)
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
