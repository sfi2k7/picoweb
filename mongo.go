package picoweb

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

func getSession() (*mgo.Session, error) {
	var err error
	if baseSession == nil {
		err = createBase()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		baseSession.SetMode(mgo.Strong, false)
	}

	if baseSession == nil {
		fmt.Println("Base Session is still null")
		return nil, err
	}

	if err = baseSession.Ping(); err != nil {
		return nil, err
	}

	return baseSession.Copy(), nil
}

func createBase() error {
	var err error
	fmt.Println("MONGO URL", mongoURL)
	baseSession, err = mgo.Dial(mongoURL)
	if err != nil {
		fmt.Println("Creating Base", err)
	}
	return err
}
