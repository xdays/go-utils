package main

import (
	"github.com/go-redis/redis"

	"fmt"
	"os"
	"time"
)

func getClient(hostname string, port int, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", hostname, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}

// this func may used for incremental migration in future
func migrateByType() {

	var cursor uint64
	var err error

	srcClient := getClient("localhost", 6379, 0)
	dstClient := getClient("localhost", 6380, 0)

	for {
		var keys []string
		if keys, cursor, err = srcClient.Scan(cursor, "", 100).Result(); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(2)
		}

		if len(keys) > 0 {
			fmt.Printf("found %d keys\n", len(keys))
			for _, key := range keys {
				t, err := srcClient.Type(key).Result()
				if err != nil {
					panic(err)
				}
				fmt.Printf("start migrating %s key %s\n", t, key)
				if t == "string" {
					v, err := srcClient.Get(key).Result()
					if err != nil {
						fmt.Println(err)
					}
					err = dstClient.Set(key, v, 0*time.Second).Err()
					if err != nil {
						panic(err)
					}
				} else if t == "list" {
					var v []string
					v, err := srcClient.LRange(key, 0, -1).Result()
					if err != nil {
						fmt.Println(err)
					}
					e, err := dstClient.Exists(key).Result()
					if e == 1 {
						fmt.Printf("delete list %s before migrate", key)
						dstClient.Del(key).Result()
					}
					err = dstClient.RPush(key, v).Err()
					if err != nil {
						panic(err)
					}
				} else if t == "hash" {
					var v map[string]string
					v, err := srcClient.HGetAll(key).Result()
					if err != nil {
						fmt.Println(err)
					}
					z := make(map[string]interface{})
					for k, value := range v {
						z[k] = value
					}
					err = dstClient.HMSet(key, z).Err()
					if err != nil {
						panic(err)
					}
				} else if t == "set" {
					var v []string
					v, err := srcClient.SMembers(key).Result()
					if err != nil {
						fmt.Println(err)
					}
					err = dstClient.SAdd(key, v).Err()
					if err != nil {
						panic(err)
					}
				} else if t == "zset" {
					var v []redis.Z
					rangeBy := redis.ZRangeBy{
						Min: "-inf",
						Max: "+inf",
					}
					v, err := srcClient.ZRangeByScoreWithScores(key, rangeBy).Result()
					if err != nil {
						fmt.Println(err)
					}
					for _, z := range v {
						err = dstClient.ZAdd(key, z).Err()
						if err != nil {
							panic(err)
						}
					}
				}
			}
		}

		if cursor == 0 {
			break
		}
	}
}

func migrateByDump() {
	var cursor uint64
	var err error

	srcClient := getClient("localhost", 6379, 0)
	dstClient := getClient("localhost", 6380, 0)

	for {
		var keys []string
		if keys, cursor, err = srcClient.Scan(cursor, "", 100).Result(); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(2)
		}

		if len(keys) > 0 {
			fmt.Printf("found %d keys\n", len(keys))
			for _, key := range keys {
				fmt.Printf("start migrating key %s\n", key)
				r, _ := srcClient.Dump(key).Result()
				dstClient.Del(key).Result()
				dstClient.Restore(key, 0, r)
				t, _ := srcClient.TTL(key).Result()
				if t != -1*time.Second {
					fmt.Printf("set ttl to %s\n", t)
					dstClient.Expire(key, t)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

}

func main() {
	migrateByDump()
}
