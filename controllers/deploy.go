package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"strings"
)

type DeployController struct {
	baseControllers
}

func (this *DeployController) Post() {
	repos := this.GetString("repos")
	tags := this.GetString("tags")
	harborhost := beego.AppConfig.String("harborhost")

	domain := strings.Split(harborhost, "//")
	image := domain[1] + "/" + repos + ":" + tags

	name := strings.Split(repos, "/")[1]

	jsonmap := make(map[string]string)
	jsonmap["name"] = name
	jsonmap["namespaces"] = "tanzhixu"
	jsonmap["image"] = image
	jsonmap["replicas"] = "1"
	jsonStr, _ := json.Marshal(jsonmap)

	url := "http://10.10.7.175:8081/deployment/update"
	result, _ := Request("POST", url, jsonStr)

	var f interface{}
	json.Unmarshal(result, &f)
	data := f.(map[string]interface{})
	data["status"] = "200"

	this.Data["json"] = data
	this.ServeJSON()
}
