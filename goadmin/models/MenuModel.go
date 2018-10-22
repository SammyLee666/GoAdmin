package models

import (
	"time"
	"encoding/gob"
	"net/url"
	"github.com/beevik/etree"
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

	Tree struct {
		Doc   *etree.Document
		Menus []Menu
	}
)

func (t *Tree) TreeView(handle *etree.Element, currentID int, option *[]Option, depth int) {
	children, isChild := t.CheckIsAChild(currentID)
	if isChild {
		ol := handle.CreateElement("ol")
		ol.CreateAttr("class", "dd-list")
		for _, v := range children {
			createHandle := t.CreateHandle(&v, ol)
			t.TreeView(createHandle, v.ID, option, depth)
		}
	}

}

func (t *Tree) CreateList() *etree.Element {
	t.Doc = etree.NewDocument()
	t.Doc.Indent(2)
	ol := t.Doc.CreateElement("ol")
	ol.CreateAttr("class", "dd-list")
	return ol
}

func (t *Tree) CreateHandle(v *Menu, ol *etree.Element) *etree.Element {
	ddItem := ol.CreateElement("li")
	ddItem.CreateAttr("class", "dd-item")
	ddItem.CreateAttr("data-id", strconv.Itoa(v.ID))

	ddHandle := ddItem.CreateElement("div")
	ddHandle.CreateAttr("class", "dd-handle")

	fa := ddHandle.CreateElement("i")
	fa.SetText("")
	fa.CreateAttr("class", "fa "+v.Icon)

	strong := ddHandle.CreateElement("strong")
	strong.SetText(v.Title)

	ddNodrag := ddHandle.CreateElement("a")
	ddNodrag.CreateAttr("class", "dd-nodrag")
	ddNodrag.CreateAttr("href", v.Uri)
	ddNodrag.SetText("    " + v.Uri)

	pullRight := ddHandle.CreateElement("span")
	pullRight.CreateAttr("class", "pull-right dd-nodrag")

	rightEdit := pullRight.CreateElement("a")
	rightEdit.CreateAttr("href", "/auth/menu/"+strconv.Itoa(v.ID)+"/edit")

	rightEditI := rightEdit.CreateElement("i")
	rightEditI.SetText("")
	rightEditI.CreateAttr("class", "fa fa-edit")

	rightDelete := pullRight.CreateElement("a")
	rightDelete.CreateAttr("class", "tree_branch_delete")
	rightDelete.CreateAttr("data-id", strconv.Itoa(v.ID))
	rightDelete.CreateAttr("href", "javascript:void(0);")

	rightDeleteI := rightDelete.CreateElement("i")
	rightDeleteI.SetText("")
	rightDeleteI.CreateAttr("class", "fa fa-trash")
	return ddItem
}

func (t *Tree) CheckIsAChild(CurrentNodeId int) (Children []Menu, status bool) {
	for _, v := range t.Menus {
		if v.Pid == CurrentNodeId {
			Children = append(Children, v)
		}
	}
	if len(Children) > 0 {
		status = true
	}
	return
}
