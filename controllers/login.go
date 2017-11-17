package controllers

import (
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Login() {
	username := this.GetString("username", "")
	password := this.GetString("password", "")

	err := CheckLogin(username, password)

	if err == nil {
		this.SetSession("username", username)
		this.Redirect("/dashboard.html", 302)
	} else {
		this.TplName = "login.html"
	}
}

func (this *LoginController) Logout() {
	username := this.GetString("username", "")

	this.DelSession(username)
	this.Redirect("/", 302)
}
