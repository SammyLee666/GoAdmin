package db

import (
	"os"
	"gopkg.in/mgo.v2"
)

var (
	Admin    *mgo.Collection
	Sessions *mgo.Collection
)

func init() {
	session, err := mgo.Dial(os.Getenv("MONGO"))
	if err != nil {
		panic(err)
	}

	Admin = session.DB("blog").C("admin")
	Sessions = session.DB("blog").C("sessions")

}
