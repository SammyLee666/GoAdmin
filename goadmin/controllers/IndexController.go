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
	"log"
	"strconv"
	"goadmin/goadmin/utils"
)

func Show(c *gin.Context) {
	var oldMsg url.Values
	var errMsg models.ErrMsg
	var menus []models.Menu
	var tree models.TreeDom
	var Li models.LiDom

	db.Mysql.Find(&menus)

	tree = models.TreeDom{Class: "dd-list"}
	option := new([]models.Option)
	*option = append(*option, models.Option{"Value": "0", "Selected": "selected", "Text": "Root"})

	RootCounter := 0
	for _, v := range menus {

		if v.Pid == 0 {
			//Create Root Node
			Li = models.LiDom{
				Class:  "dd-item",
				DataId: strconv.Itoa(v.ID),
				Handle: models.DivHandleDom{
					Class: "dd-handle",
					Fa: models.IDom{
						Class: "fa " + v.Icon,
					},
					Strong: v.Title,
					A: models.ADom{
						Href:  v.Uri,
						Class: "dd-nodrag",
						Title: "&nbsp;&nbsp;" + v.Uri,
					},
					Span: models.SpanDom{
						Class: "pull-right dd-nodrag",
						AList: []models.AListDom{
							{
								Href: "/auth/menu/" + strconv.Itoa(v.ID) + "/edit",
								I: models.IDom{
									Class: "fa fa-edit",
								},
							},
							{
								Href:   "javascript:void(0);",
								DataId: strconv.Itoa(v.ID),
								Class:  "tree_branch_delete",
								I: models.IDom{
									Class: "fa fa-trash",
								},
							},
						},
					},
				},
				//ChildDom: tree,
			}
			tree.List = append(tree.List, Li)

			//Is a Child Node
			handle := models.TreeDom{Class: "dd-list"}
			models.TreeParse(menus, v.ID, &handle, v.Title, option, 1)
			tree.List[RootCounter].ChildDom = handle
			RootCounter ++
		}
	}

	output, err := xml.MarshalIndent(tree, "  ", "    ")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	html := string(output)

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
	//TreeView
	models.TreeView()

	c.HTML(http.StatusOK, "goadmin/layout/index", gin.H{
		"_errors": errMsg,
		"_old":    oldMsg,
		"_toastr": toastr,
		"tree":    html,
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
		model := db.Mysql.Create(&menu)
		fmt.Println(model)

		utils.Toastr(c).Success("提交成功！")
		c.Redirect(http.StatusFound, c.Request.Header.Get("Referer"))
		return
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
