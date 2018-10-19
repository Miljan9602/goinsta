package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func main()  {

	redisClient := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
		ReadTimeout: time.Minute,
	})

	_, err := redisClient.Ping().Result()

	if err != nil {
		fmt.Println(fmt.Sprintf(redisConnectionFailedPatternError, err))
		return
	}

	keys,err := loadKeys(redisClient)

	if err != nil {
		fmt.Println(fmt.Sprintf(failedToLoadKeysPatternError, err))
		return
	}

	accounts, err := loadAccounts(redisClient, keys)

	if err != nil {
		fmt.Println(fmt.Sprintf(failedToLoadAccountsPatternError, err))
		return
	}

	if size := len(loginAccounts(accounts)); size != 0 {
		fmt.Println(fmt.Sprintf(failedToLoginAccountsPatternError, strconv.Itoa(size)))
	}

	fmt.Println("Start")
}