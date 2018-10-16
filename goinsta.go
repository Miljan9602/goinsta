package goinsta

import (
	"crypto/tls"
	"encoding/json"
	"github.com/sethgrid/pester"
	"net"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"time"
)

// Instagram represent the main API handler
//
// Profiles: Represents instragram's user profile.
// Account:  Represents instagram's personal account.
// Search:   Represents instagram's search.
// Timeline: Represents instagram's timeline.
// Activity: Represents instagram's user activity.
// Inbox:    Represents instagram's messages.
//
// See Scheme section in README.md for more information.
//
// We recommend to use Export and Import functions after first Login.
//
// Also you can use SetProxy and UnsetProxy to set and unset proxy.
// Golang also provides the option to set a proxy using HTTP_PROXY env var.
type Instagram struct {
	user string
	pass string
	// device id: android-1923fjnma8123
	dID string
	// uuid: 8493-1233-4312312-5123
	uuid string
	// rankToken
	rankToken string
	// token
	token string
	// phone id
	pid string
	// ads id
	adid string

	ID int64

	// ApiPath for challenge.
	Nonce string

	// Every account will have it's own user-agent.
	userAgent string

	// Provider for exporting/importing account details and cookies.
	Provider Provider

	// Profiles is the user interaction
	Profiles *Profiles
	// Account stores all personal data of the user and his/her options.
	Account *Account
	// Search performs searching of multiple things (users, locations...)
	Search *Search
	// Timeline allows to receive timeline media.
	Timeline *Timeline
	// Activity are instagram notifications.
	Activity *Activity
	// Inbox are instagram message/chat system.
	Inbox *Inbox
	// Feed for search over feeds
	Feed *Feed

	Challenges *Challenges

	c *pester.Client
}

// SetDeviceID sets device id
func (inst *Instagram) SetDeviceID(id string) {
	inst.dID = id
}

// SetUUID sets uuid
func (inst *Instagram) SetUUID(uuid string) {
	inst.uuid = uuid
}

// SetPhoneID sets phone id
func (inst *Instagram) SetPhoneID(id string) {
	inst.pid = id
}

func (insta *Instagram) GetCookies() ([]*http.Cookie, error) {

	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return nil, err
	}
	return insta.c.Jar.Cookies(url), nil
}

// New creates Instagram structure
func (inst *Instagram) SetUser(username string, pk int64, password string) {

	inst.user = username
	inst.ID = pk
	inst.pass = password
	userAgent, deviceID := UserAgentAndDevice(username, password)
	inst.dID = deviceID
	inst.userAgent = userAgent
	inst.uuid = generateUUID()
	inst.pid = generateUUID()
	inst.Provider.Import()
}

func (inst *Instagram) init() {
	inst.Profiles = newProfiles(inst)
	inst.Activity = newActivity(inst)
	inst.Timeline = newTimeline(inst)
	inst.Search = newSearch(inst)
	inst.Inbox = newInbox(inst)
	inst.Feed = newFeed(inst)
	inst.Challenges = newChallenge(inst)
}
func New(p Provider, options *PesterOptions) *Instagram {
	inst := &Instagram{}
	inst.init()
	inst.Provider = p
	inst.Provider.SetInstagram(inst)

	// this call never returns error
	jar, _ := cookiejar.New(nil)

	inst.c = pester.New()
	inst.c.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	inst.c.Jar = jar

	if options == nil {
		options = DefaultOptions()
	}

	inst.c.Concurrency = options.Concurrency
	inst.c.Backoff = options.Backoff
	inst.c.MaxRetries = options.MaxRetries

	return inst
}

// SetProxy sets proxy for connection.
func (inst *Instagram) SetProxy(url string, insecure bool) error {
	uri, err := neturl.Parse(url)
	if err == nil {
		inst.c.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		}
	}
	return err
}

// Set TCP end point,
func (inst *Instagram) SetTCP(address string) {

	localAddr := &net.TCPAddr{
		IP: net.ParseIP(address),
	}

	inst.c.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			LocalAddr: localAddr,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives: true,
	}
}

// UnsetProxy unsets proxy for connection.
func (inst *Instagram) UnsetProxy() {
	inst.c.Transport = nil
}

func (inst *Instagram) readMsisdnHeader() error {
	data, err := json.Marshal(
		map[string]string{
			"device_id": inst.uuid,
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint:   urlMsisdnHeader,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(b2s(data)),
		},
	)
	return err
}

func (inst *Instagram) contactPrefill() error {
	data, err := json.Marshal(
		map[string]string{
			"phone_id":   inst.pid,
			"_csrftoken": inst.token,
			"usage":      "prefill",
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint:   urlContactPrefill,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(b2s(data)),
		},
	)
	return err
}

func (inst *Instagram) zrToken() error {
	_, err := inst.sendRequest(
		&reqOptions{
			Endpoint:   urlZrToken,
			IsPost:     false,
			Connection: "keep-alive",
			Query: map[string]string{
				"device_id":        inst.dID,
				"token_hash":       "",
				"custom_device_id": inst.uuid,
				"fetch_reason":     "token_expired",
			},
		},
	)
	return err
}

func (inst *Instagram) sendAdID() error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"adid": inst.adid,
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint:   urlLogAttribution,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(data),
		},
	)
	return err
}

func (inst *Instagram) Login() error {
	return inst.FullLogin()
}

// Login performs instagram login.
//
// Password will be deleted after login
func (inst *Instagram) FullLogin() error {

	err := inst.readMsisdnHeader()
	if err != nil {
		return err
	}

	err = inst.syncFeatures()
	if err != nil {
		return err
	}

	err = inst.zrToken()
	if err != nil {
		return err
	}

	err = inst.sendAdID()
	if err != nil {
		return err
	}

	err = inst.contactPrefill()
	if err != nil {
		return err
	}

	if _, err := inst.Profiles.ByName(inst.user); err == nil {
		return nil
	}

	result, err := json.Marshal(
		map[string]interface{}{
			"guid":                inst.uuid,
			"login_attempt_count": 0,
			"_csrftoken":          inst.token,
			"device_id":           inst.dID,
			"adid":                inst.adid,
			"phone_id":            inst.pid,
			"username":            inst.user,
			"password":            inst.pass,
			"google_tokens":       "[]",
		},
	)
	if err != nil {
		return err
	}
	body, err := inst.sendRequest(
		&reqOptions{
			Endpoint: urlLogin,
			Query:    generateSignature(b2s(result)),
			IsPost:   true,
			Login:    true,
		},
	)
	if err != nil {
		return err
	}
	inst.pass = ""

	// getting account data
	res := accountResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	inst.Account = &res.Account
	inst.Account.inst = inst
	inst.ID = inst.Account.ID
	inst.rankToken = strconv.FormatInt(inst.Account.ID, 10) + "_" + inst.uuid
	inst.zrToken()

	return err
}

// Logout closes current session
func (inst *Instagram) Logout() error {
	_, err := inst.sendSimpleRequest(urlLogout)
	inst.c.Jar = nil
	inst.c = nil
	return err
}

func (inst *Instagram) syncFeatures() error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"id":          inst.uuid,
			"experiments": goInstaExperiments,
		},
	)
	if err != nil {
		return err
	}

	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint: urlQeSync,
			Query:    generateSignature(data),
			IsPost:   true,
			Login:    true,
		},
	)
	return err
}

func (inst *Instagram) megaphoneLog() error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"id":        inst.Account.ID,
			"type":      "feed_aysf",
			"action":    "seen",
			"reason":    "",
			"device_id": inst.dID,
			"uuid":      generateMD5Hash(string(time.Now().Unix())),
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint: urlMegaphoneLog,
			Query:    generateSignature(data),
			IsPost:   true,
			Login:    true,
		},
	)
	return err
}

func (inst *Instagram) expose() error {
	data, err := inst.prepareData(
		map[string]interface{}{
			"id":         inst.Account.ID,
			"experiment": "ig_android_profile_contextual_feed",
		},
	)
	if err != nil {
		return err
	}

	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint: urlExpose,
			Query:    generateSignature(data),
			IsPost:   true,
		},
	)

	return err
}

// GetMedia returns media specified by id.
//
// The argument can be int64 or string
//
// See example: examples/media/like.go
func (inst *Instagram) GetMedia(o interface{}) (*FeedMedia, error) {
	media := &FeedMedia{
		inst:   inst,
		NextID: o,
	}
	return media, media.Sync()
}
