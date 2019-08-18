package main

import (
	"flag"
	"fmt"

	"github.com/xdays/go-utils/v2ray-sdk"
)

func main() {
	uid := flag.String("i", "", "create a new user")
	globalStat := flag.Bool("g", false, "get global stats")
	userStat := flag.Bool("u", false, "get user stats")
	add := flag.Bool("a", false, "add a new user")
	delete := flag.Bool("d", false, "delete a new user")
	flag.Parse()

	var u v2ray.User
	var err error
	v := v2ray.Client{
		Host: "127.0.0.1",
		Port: 10085,
	}
	c, err := v.GetConnection()
	if err != nil {
		panic(err)
	}

	if *add || *delete {
		if *uid == "" {
			panic("uid is required")
		}
		u = v2ray.User{
			Level:   0,
			Email:   fmt.Sprintf("%s@x.xx", *uid),
			UUID:    *uid,
			AlterID: 64,
		}
		hsClient := v.GetBondClient(c)
		if *add {
			_, err := v.AddUser(hsClient, "default", u)
			if err != nil {
				panic(err)
			}
		} else if *delete {
			_, err := v.RemoveUser(hsClient, "default", u)
			if err != nil {
				panic(err)
			}
		}

	}
	if *globalStat || *userStat {
		qsClient := v.GetStatClient(c)
		if *globalStat {
			fmt.Println(v.QueryStats(qsClient))

		} else if *userStat {
			if *uid == "" {
				panic("uid is required")
			}
			u = v2ray.User{
				Level:   0,
				Email:   fmt.Sprintf("%s@x.xx", *uid),
				UUID:    *uid,
				AlterID: 64,
			}
			fmt.Println(v.GetStats(qsClient, u))
		}
	}
}
