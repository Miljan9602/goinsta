package main

import "../.."

/**
 */
type Request interface {
	Save() error
	Execute() error
	SetClient(inst *goinsta.Instagram)
}


type FollowingRequest struct {

}


func (req *FollowingRequest) Save() error {
	return nil;
}

func (req *FollowingRequest) Execute() error {
	return nil;
}

func (req *FollowingRequest) SetClient(inst *goinsta.Instagram) {

}
