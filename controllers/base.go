package controllers

import (
	"github.com/astaxie/beego"
)

type baseControllers struct {
	beego.Controller
}

func (this *baseControllers) Prepare() {
	username := this.GetSession("username").(string)

	this.Data["User"] = username

	if CheckAdmin(username) {
		this.Data["IsAdmin"] = true
	} else {
		this.Data["IsAdmin"] = false
	}
}
