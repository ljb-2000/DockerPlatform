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
	RepoName string
	Url      string
	User     *User `orm:"rel(fk)"`
}

type Message struct {
	Message string `json:"message"`
}

func RegisterDB() {

	orm.RegisterModel(new(User), new(Bitbucket))

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/dockerplatform?charset=utf8", 30)

	orm.RunSyncdb("default", false, true)

}
