package main

import (
	"../.."
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sethgrid/pester"
	"strconv"
)

// Load all keys from database.
func loadKeys(rc *redis.Client) ([]string, error) {

	res, err := rc.SMembers(idSetsAccount).Result()

	if err != nil {
		return nil, err
	}

	return res, nil
}

// Load all accounts.
func loadAccounts(rc *redis.Client, keys []string) ([]*goinsta.Instagram, error) {

	accounts := make([]*goinsta.Instagram, 0)



	for _, element := range keys {

		res, err := rc.HGetAll(fmt.Sprintf(basicInfoAccountPattern, element)).Result()

		if err != nil {
			return nil, err
		}

		if len(res[username]) < 1 || len(res[password]) < 1 || len(res[key]) < 1{
			return nil, errors.New(fmt.Sprintf(missingPropertyBasicInfoAccountPatternError, element))
		}


		cli, _ := goinsta.NewRedisProvider(rc, hashAccount)


		insta := goinsta.New(cli, &goinsta.PesterOptions{
			Concurrency : pester.DefaultClient.Concurrency,
			MaxRetries : 3,
			Backoff :  pester.DefaultClient.Backoff,
		})

		key, err := strconv.ParseInt(res[key], 10, 64)
		if err != nil {
			return nil, err
		}

		insta.SetUser(res[username], key, res[password])

		accounts = append(accounts, insta)
	}

	return accounts, nil
}


func loginAccounts(accounts []*goinsta.Instagram) []*goinsta.Instagram{

	notLoggedAccounts := make([]*goinsta.Instagram, 0)

	for _,acc := range accounts {

		if err := acc.Login(); err != nil {
			fmt.Println("Error: ",err)
			notLoggedAccounts = append(notLoggedAccounts, acc)
		}
	}

	return notLoggedAccounts
}
