package goinsta

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
)

type Challenges struct {
	inst *Instagram
}

func newChallenge(inst *Instagram) *Challenges {
	challenges := &Challenges{
		inst: inst,
	}
	return challenges
}

func (ch *Challenges) GetInfo(pk int64, nonce string) (*SelectVerifyMethod, error) {

	// Save data for later use.
	ch.inst.Nonce = nonce
	ch.inst.ID = pk

	ApiPath := "challenge/" + strconv.FormatInt(pk, 10) + "/" + nonce

	body, err := ch.inst.sendRequest(&reqOptions{
		Endpoint: ApiPath,
		Query: map[string]string{
			"device_id": ch.inst.dID,
		},
	})

	if err != nil {
		return nil, err
	}

	resp := SelectVerifyMethod{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (ch *Challenges) SelectVerifyMethod(choice string) (*InstagramVerification, error) {

	nonce := ch.inst.Nonce
	pk := ch.inst.ID

	ApiPath := "challenge/" + strconv.FormatInt(pk, 10) + "/" + nonce

	data := ch.inst.prepareDataQuery(
		map[string]interface{}{
			"choice":    choice,
			"device_id": ch.inst.dID,
		},
	)

	body, err := ch.inst.sendRequest(
		&reqOptions{
			Endpoint: ApiPath,
			Query:    data,
			IsPost:   true,
		},
	)

	fmt.Println("Response",string(body[:]))

	resp := InstagramVerification{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (ch *Challenges) VerifyCode(code string) error {

	nonce := ch.inst.Nonce
	pk := ch.inst.ID

	ApiPath := "challenge/" + strconv.FormatInt(pk, 10) + "/" + nonce

	data := ch.inst.prepareDataQuery(
		map[string]interface{}{
			"security_code": code,
			"device_id":     ch.inst.dID,
		},
	)

	body, err := ch.inst.sendRequest(
		&reqOptions{
			Endpoint: ApiPath,
			Query:    data,
			IsPost:   true,
		},
	)

	if err != nil {
		return err
	}

	if !gjson.Get(string(body[:]), "logged_in_user").Bool() {
		return errors.New("logged_in_user does not exist")
	}

	res := accountResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	ch.inst.Account = &res.Account
	ch.inst.Account.inst = ch.inst
	ch.inst.rankToken = strconv.FormatInt(ch.inst.Account.ID, 10) + "_" + ch.inst.uuid
	ch.inst.zrToken()
	return nil
}
