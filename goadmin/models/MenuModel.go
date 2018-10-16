package models

import (
	"time"
	"encoding/gob"
)

func init() {
	gob.Register(ErrMsg{})
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
)
