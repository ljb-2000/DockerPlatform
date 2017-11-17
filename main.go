package main

import (
	"DockerPlatform/controllers"
	"DockerPlatform/models"
	_ "DockerPlatform/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

func CheckGitUrl(in string) (out string) {
	b := strings.Split(in, " ")
	for _, v := range b {
		if !strings.Contains(v, "ssh") {
			return v
		}
	}
	return ""
}

func timeParse(in string) (out string) {
	formate := "2006-01-02 15:04:05"

	in = strings.Split(in, ".")[0]
	s := strings.Replace(in, "T", " ", -1)

	t, _ := time.Parse(formate, s)
	out = t.Local().Format(formate)
	return
}

func init() {
	models.RegisterDB()
}

func main() {
	dockerdatadir := beego.AppConfig.String("datadir")
	controllers.MkdirIfNoExist(dockerdatadir)

	o := orm.NewOrm()
	var user models.User
	user.Name = "admin"
	user.Password = "admin"
	user.Isadmin = true
	o.ReadOrCreate(&user, "Name")

	beego.SetStaticPath("/static", "static")
	beego.AddFuncMap("checkgiturl", CheckGitUrl)
	beego.AddFuncMap("timeparse", timeParse)
	beego.Run()

}
