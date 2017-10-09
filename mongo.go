package picoweb

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

func getSession() (*mgo.Session, error) {
	var err error
	if baseSession == nil {
		err = createBase()
		fmt.Println(err)
	}

	if baseSession == nil {
		fmt.Println("Base Session is still null")
		return nil, err
	}
	if err = baseSession.Ping(); err != nil {
		return nil, err
	}
	return baseSession.Clone(), nil
}

func createBase() error {
	var err error
	baseSession, err = mgo.Dial("mongodb://127.0.0.1")
	return err
}
