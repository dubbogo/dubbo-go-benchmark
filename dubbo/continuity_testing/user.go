package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

import (
	"github.com/apache/dubbo-go/config"
	"github.com/dubbogo/hessian2"
)

type Gender hessian.JavaEnum

var userProvider = new(UserProvider)

func init() {
	config.SetConsumerService(userProvider)
}

const (
	MAN hessian.JavaEnum = iota
	WOMAN
)

var genderName = map[hessian.JavaEnum]string{
	MAN:   "MAN",
	WOMAN: "WOMAN",
}

var genderValue = map[string]hessian.JavaEnum{
	"MAN":   MAN,
	"WOMAN": WOMAN,
}

func (g Gender) JavaClassName() string {
	return "com.ikurento.user.Gender"
}

func (g Gender) String() string {
	s, ok := genderName[hessian.JavaEnum(g)]
	if ok {
		return s
	}

	return strconv.Itoa(int(g))
}

func (g Gender) EnumValue(s string) hessian.JavaEnum {
	v, ok := genderValue[s]
	if ok {
		return v
	}

	return hessian.InvalidJavaEnum
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

func (u User) JavaClassName() string {
	return "com.ikurento.user.User"
}

type UserProvider struct {
	GetUser func(ctx context.Context, req []interface{}, rsp *User) error
	Echo    func(req interface{}) (string, error) // Echo represent EchoFilter will be used
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
