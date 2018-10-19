package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
)

type Provider interface {
	Export() error
	Import() error
	SetInstagram(insta *Instagram)
}

type RedisProvider struct {
	inst *Instagram

	client *redis.Client

	tableName string
}

func (redis *RedisProvider) Export() error {
	cookies, err := redis.inst.GetCookies()

	if err != nil {
		return err
	}

	config := ConfigFile{
		ID:        redis.inst.ID,
		User:      redis.inst.user,
		DeviceID:  redis.inst.dID,
		Nonce:     redis.inst.Nonce,
		UUID:      redis.inst.uuid,
		RankToken: redis.inst.rankToken,
		Token:     redis.inst.token,
		UserAgent: redis.inst.userAgent,
		PhoneID:   redis.inst.pid,
		Cookies:   cookies,
	}

	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	field := strconv.FormatInt(redis.inst.ID, 10)
	_,err = redis.client.HSet(redis.tableName,field, string(bytes[:])).Result()

	return err
}

func (redis *RedisProvider) Import() error {

	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return err
	}
	field := strconv.FormatInt(redis.GetInstagram().ID, 10)
	res, err := redis.client.HGet(redis.tableName, field).Result()

	if err != nil {
		return err
	}

	if len(res) == 0{
		return errors.New(fmt.Sprintf("ID: %s Does not exist in redis", field))
	}

	bytes := []byte(res)
	config := ConfigFile{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}


	inst := redis.GetInstagram()
	inst.user = config.User
	inst.dID = config.DeviceID
	inst.userAgent = config.UserAgent
	inst.uuid = config.UUID
	inst.rankToken = config.RankToken
	inst.token = config.Token
	inst.pid = config.PhoneID
	inst.Nonce = config.Nonce

	inst.c.Jar, err = cookiejar.New(nil)
	if err != nil {
		return err
	}
	inst.c.Jar.SetCookies(url, config.Cookies)
	inst.init()
	inst.Account = &Account{inst: inst, ID: config.ID}
	inst.Account.Sync()

	return err
}

func (redis *RedisProvider) SetInstagram(insta *Instagram) {
	redis.inst = insta
}

func (redis *RedisProvider) GetInstagram() *Instagram {
	return redis.inst
}

func (redis *RedisProvider) GetClient() *redis.Client {
	return redis.client
}

func NewRedisProvider(cl *redis.Client, tableName string) (*RedisProvider, error) {

	_, err := cl.Ping().Result()

	if err != nil {
		return nil,err
	}

	return &RedisProvider{client: cl, tableName: tableName}, nil
}

