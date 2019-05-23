package main

import (
	"context"
	"fmt"
	"time"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/config"
)

type Gender int

func init() {
	config.SetProviderService(new(UserProvider))
}

const (
	MAN = iota
	WOMAN
)

var genderStrings = [...]string{
	"MAN",
	"WOMAN",
}

func (g Gender) String() string {
	return genderStrings[g]
}

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Time      time.Time `json:"time"`
	Sex       Gender    `json:"sex"`
	IsChinese bool      `json:"isChinese"`
	Remarks   string    `json:"remarks"`
}

// UserProvider -------------------------------------------------
type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, req []interface{}, rsp *User) error {
	rsp.Id = req[0].(string)
	rsp.Name = "name"
	rsp.Age = 20
	rsp.Sex = Gender(MAN)
	rsp.Time = time.Now()
	rsp.IsChinese = true
	rsp.Remarks = req[1].(string)
	fmt.Printf("GetUser")
	return nil
}

func (u *UserProvider) Service() string {
	return "com.ikurento.user.UserProvider"
}

func (u *UserProvider) Version() string {
	return ""
}
