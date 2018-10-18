package models

import (
	"time"
	"encoding/gob"
	"net/url"
	"fmt"
	"encoding/xml"
)

func init() {
	gob.Register(ErrMsg{})
	gob.Register(url.Values{})
}

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

	TreeDom struct {
		XMLName xml.Name `xml:"ol"`
		Class   string   `xml:"class,attr"`
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
	//var menus []Menu
	//var treeString string
	//db.Mysql.Find(&menus)
	//treeParse()
}

func treeParse(childs []*Child, PID int) {
	for _, v := range childs {
		//Create
		fmt.Println(PID, v.ID, v.Child)

		if len(v.Child) > 0 {
			//解析Child
			treeParse(v.Child, v.ID)
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
