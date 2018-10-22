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
	"log"
	"goadmin/goadmin/utils"
)

func Show(c *gin.Context) {
	var oldMsg url.Values
	var errMsg models.ErrMsg
	var menus []models.Menu
	option := new([]models.Option)

	db.Mysql.Find(&menus)
	//a := []int{1}
	//db.Mysql.Exec("update menus set `order` = ? where id = ?;update menus set `order` = ? where id = ?;update menus set `order` = ? where id = ?;" )

	//TreeView
	tree := models.Tree{Menus: menus}
	ol := tree.CreateList()
	for _, v := range menus {
		if v.Pid == 0 {
			handle := tree.CreateHandle(&v, ol)
			tree.TreeView(handle, v.ID, option, 1)
		}
	}
	treeHtml, err := tree.Doc.WriteToString()
	if err != nil {
		log.Println(err.Error())
	}

	session := sessions.Default(c)
	errors := session.Flashes("errors")
	oldForm := session.Flashes("oldForm")
	toastr := session.Flashes("toastr")
	session.Save()
	if len(errors) > 0 {
		errMsg = errors[0].(models.ErrMsg)
	}
	if len(oldForm) > 0 {
		oldMsg = oldForm[0].(url.Values)
	}

	c.HTML(http.StatusOK, "goadmin/layout/index", gin.H{
		"_errors": errMsg,
		"_old":    oldMsg,
		"_toastr": toastr,
		"tree":    treeHtml,
		"select":  option,
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
				log.Printf("session.Save %s", err.Error())
				return
			}
			c.Redirect(http.StatusFound, c.Request.Header.Get("Referer"))
			return
		}

		db.Mysql.Create(&menu)
		utils.Toastr(c).Success("提交成功！")
		c.Redirect(http.StatusFound, c.Request.Header.Get("Referer"))
		return
	} else {
		changeParentId(order)
		c.JSON(http.StatusOK,gin.H{"message":"保存成功 !"})
	}

}

func Edit(c *gin.Context) {
	fmt.Println(1)
}

func Dele(c *gin.Context) {

}

func changeParentId(order string) {
	var childs []*models.Child
	var up []map[string]int
	json.Unmarshal([]byte(order), &childs)

	parse(childs, 0, &up)
	stmt, err := db.Mysql.DB().Prepare("UPDATE menus SET pid = ? WHERE id = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	for _, v := range up {
		stmt.Exec(v["pid"], v["id"])
	}
}

func parse(childs []*models.Child, PID int, up *[]map[string]int) {
	for _, v := range childs {
		//Create
		*up = append(*up, map[string]int{"id": v.ID, "pid": PID})
		//db.Mysql.Debug().Model(models.Menu{ID: v.ID}).Update("pid", PID)
		if len(v.Child) > 0 {
			//解析Child
			parse(v.Child, v.ID, up)
		}
	}
	return
}
