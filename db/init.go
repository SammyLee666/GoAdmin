package db

import (
	"os"
	"gopkg.in/mgo.v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Mysql    *gorm.DB
	Sessions *mgo.Collection
)

func init() {
	//session
	session, err := mgo.Dial(os.Getenv("MONGO"))
	if err != nil {
		panic(err)
	}
	Sessions = session.DB("blog").C("sessions")

	//db
	db, err := gorm.Open("mysql", os.Getenv("MYSQL"))
	//defer db.Close()
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	Mysql = db

}
