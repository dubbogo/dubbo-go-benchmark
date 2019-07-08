package main

import (
	"context"
	"fmt"
	"time"
)

import (
	"github.com/apache/dubbo-go/config"
)

var userProvider = new(UserProvider)

func init() {
	config.SetConsumerService(userProvider)
}

type Gender int

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

// User -------------------------------------------------
type User struct {
	Id        string
	Name      string
	Age       int32
	Time      time.Time
	Sex       Gender
	IsChinese bool
	Remarks   string
}

func (u User) String() string {
	return fmt.Sprintf(
		"User{Id:%s, Name:%s, Age:%d, Time:%s, Sex:%s, Country:%v}",
		u.Id, u.Name, u.Age, u.Time, u.Sex, u.IsChinese,
	)
}

type UserProvider struct {
	GetUser func(ctx context.Context, req []interface{}, rsp *User) error
	Echo    func(req interface{}) (string, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
