package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id        int
	Name      string
	Password  string
	Isadmin   bool
	Bitbucket []*Bitbucket `orm:"reverse(many)"`
}

type Bitbucket struct {
	Id       int
	Project  string
	RepoName string
	Url      string
	User     *User     `orm:"rel(fk)"`
	Pipeline *Pipeline `orm:"reverse(one)"`
}

type Pipeline struct {
	Id         int
	Name       string
	Version    string
	BuildMsg   string
	Dockerfile string     `orm:"size(1000)"`
	Bitbucket  *Bitbucket `orm:"rel(one)"`
}

type Message struct {
	Message string `json:"message"`
}

func RegisterDB() {

	orm.RegisterModel(new(User), new(Bitbucket), new(Pipeline))

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(10.10.35.1:32770)/dockerplatform?charset=utf8", 30)

	orm.RunSyncdb("default", false, true)

}
