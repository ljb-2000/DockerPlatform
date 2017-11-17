package controllers

import (
	"encoding/json"
	// "github.com/astaxie/beego"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"strings"
)

type DeploymentController struct {
	baseControllers
}

func (this *DeploymentController) Get() {
	username := this.GetSession("username").(string)

	namespaces := ""
	if CheckAdmin(username) {
		namespaces = this.GetString("namespaces")
	} else {
		namespaces = username
	}

	url := "http://10.10.7.175:8081/deployment/list"

	//json 序列化
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	//post数据
	result, _ := Request("POST", url, jsonStr)

	// json数据错误，去掉(MISSING)字符串
	str := string(result)
	str = strings.Replace(str, "(MISSING)", "", -1)

	// 重新构建json格式数据
	jsonStr2 := []byte(str)

	var deploymentlist []v1beta1.Deployment
	json.Unmarshal(jsonStr2, &deploymentlist)

	this.Data["deploymentlist"] = deploymentlist
	this.Layout = "layout.html"
	this.TplName = "deployment.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
