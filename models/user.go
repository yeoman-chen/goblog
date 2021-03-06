package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//用户模型
type User struct {
	Id         int
	UserName   string    `orm:"unique;size(15)"`
	Password   string    `orm:"size(32)"`
	Email      string    `orm:"size(50)"`
	LastLogin  time.Time `orm:"auto_now_add;type(datetime)"`
	LoginCount int
	LastIp     string `orm:"size(32)"`
	Authkey    string `orm:"size(10)"`
	Active     int8
}

func (m *User) TableName() string {
	return TableName("user")
}

func (m *User) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(field ...string) error {
	if err := orm.NewOrm().Read(m, field...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(field ...string) error {
	if _, err := orm.NewOrm().Update(m, field...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}
