package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

import (
	"github.com/dubbogo/hessian2"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/config"
)

type Gender hessian.JavaEnum

func init() {
	config.SetProviderService(new(UserProvider))
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
		"User{Id:%s, Name:%s, Age:%d, Time:%s, Sex:%s, Country:%v, Remarks:%s}",
		u.Id, u.Name, u.Age, u.Time, u.Sex, u.IsChinese, u.Remarks,
	)
}

func (u User) JavaClassName() string {
	return "com.ikurento.user.User"
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
