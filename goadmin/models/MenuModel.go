package models

import (
	"time"
	"encoding/gob"
	"net/url"
	"encoding/xml"
	"strconv"
)

func init() {
	gob.Register(ErrMsg{})
	gob.Register(url.Values{})
}

const Indent = `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`

type (
	//Table
	Menu struct {
		ID        int    `gorm:"primary_key"`
		Pid       int    `gorm:"type:int(5);not null;" form:"parent_id" binding:"min=0,max=100000000"`
		Order     int    `gorm:"type:int(5);not null;"`
		Title     string `gorm:"type:varchar(50);not null;" form:"title" binding:"min=1,max=50"`
		Icon      string `gorm:"type:varchar(50);not null;" form:"icon" binding:"min=1,max=50"`
		Uri       string `gorm:"type:varchar(50);not null;" form:"uri" binding:"min=1,max=50"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	//tree struct
	Child struct {
		ID    int      `json:"id"`
		Child []*Child `json:"children"`
	}

	ErrMsg map[string]interface{}

	Option map[string]string

	TreeDom struct {
		XMLName xml.Name `xml:"ol,omitempty"`
		Class   string   `xml:"class,attr,omitempty"`
		List    []LiDom  `xml:"li"`
	}

	LiDom struct {
		Class    string       `xml:"class,attr,omitempty"`
		DataId   string       `xml:"data-id,attr"`
		Handle   DivHandleDom `xml:"div"`
		ChildDom TreeDom      `xml:"ol,omitempty"`
	}

	IDom struct {
		Class string `xml:"class,attr"`
	}

	ADom struct {
		Href  string `xml:"href,attr"`
		Class string `xml:"class,attr"`
		Title string `xml:",innerxml"`
	}

	AListDom struct {
		Href   string `xml:"href,attr"`
		DataId string `xml:"data-id,attr,omitempty"`
		Class  string `xml:"class,attr,omitempty"`
		I      IDom   `xml:"i"`
	}

	SpanDom struct {
		Class string     `xml:"class,attr"`
		AList []AListDom `xml:"a"`
	}

	DivHandleDom struct {
		Class  string  `xml:"class,attr"`
		Fa     IDom    `xml:"i"`
		Strong string  `xml:"strong"`
		A      ADom    `xml:"a"`
		Span   SpanDom `xml:"span"`
	}
)

func TreeView() {

}

func TreeParse(menus []Menu, currentID int, handle *TreeDom, title string, option *[]Option, depth int) {
	var currentIndent string
	for i := 0; i < depth; i++ {
		currentIndent += Indent
	}
	*option = append(*option, Option{"Value":strconv.Itoa(currentID),"Text":currentIndent + title})
	//Is a Child Node
	children, isChild := CheckIsAChild(menus, currentID)
	if isChild {
		//Create Child Node
		for _, v := range children {
			Li := LiDom{
				Class:  "dd-item",
				DataId: strconv.Itoa(v.ID),
				Handle: DivHandleDom{
					Class: "dd-handle",
					Fa: IDom{
						Class: "fa " + v.Icon,
					},
					Strong: v.Title,
					A: ADom{
						Href:  v.Uri,
						Class: "dd-nodrag",
						Title: "&nbsp;&nbsp;" + v.Uri,
					},
					Span: SpanDom{
						Class: "pull-right dd-nodrag",
						AList: []AListDom{
							{
								Href: "/auth/menu/" + strconv.Itoa(v.ID) + "/edit",
								I: IDom{
									Class: "fa fa-edit",
								},
							},
							{
								Href:   "javascript:void(0);",
								DataId: strconv.Itoa(v.ID),
								Class:  "tree_branch_delete",
								I: IDom{
									Class: "fa fa-trash",
								},
							},
						},
					},
				},
			}
			handle.List = append(handle.List, Li)
			TreeParse(children, v.ID, handle, v.Title, option, depth+1)
		}
	}
}

func CheckIsAChild(menus []Menu, CurrentNodeId int) (Children []Menu, status bool) {
	for _, v := range menus {
		if v.Pid == CurrentNodeId {
			Children = append(Children, v)
		}
	}
	if len(Children) > 0 {
		status = true
	}
	return
}
