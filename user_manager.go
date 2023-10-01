package picoweb

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

/*
	Use Mongo as User Manager store
		Ask for URL in UseUserManager() method

	Register:
		username, password, meta map
		Encrypt password - add some salt
	Account Status:
		enabled / disabled
	Login:
		username, password

	Loggoff



*/

var (
	red *redis.Client
)

type user struct {
	id   string
	name string
}

type usermanager struct {
	appname string
}

func (um *usermanager) IsAnonymouse(session string) bool {
	if red == nil {
		return true
	}

	if len(session) == 0 {
		return true
	}

	v, _ := red.Get(context.Background(), um.appname+session).Result()
	return len(v) == 0
}

func (um *usermanager) Login(username, session string) bool {
	if red == nil {
		return false
	}

	_, err := red.Set(context.Background(), um.appname+session, username, time.Minute*10).Result()
	return err == nil
}

func (um *usermanager) Refresh(session string) bool {
	if red == nil {
		return false
	}

	u, err := red.Get(context.Background(), um.appname+session).Result()
	if err != nil {
		return false
	}

	if len(u) == 0 {
		return false
	}

	_, err = red.Set(context.Background(), um.appname+session, u, time.Minute*10).Result()
	return err == nil
}

func (um *usermanager) Logoff(session string) bool {
	if red == nil {
		return false
	}

	_, err := red.HDel(context.Background(), um.appname+"_usermanager", session).Result()
	return err == nil
}

func useUserManager(url, password string) {
	red = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     url,
		Password: password,
		DB:       0,
	})
}
