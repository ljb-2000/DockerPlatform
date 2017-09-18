package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"strings"
)

type DeploymentController struct {
	beego.Controller
}

func (c *DeploymentController) Get() {
	namespaces := c.GetString("namespaces", "")
	url := "http://10.10.7.175:8081/deployment/list"

	//json 序列化
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	//post数据
	result, _ := DoPost(url, jsonStr)

	// json数据错误，去掉(MISSING)字符串
	str := string(result)
	str = strings.Replace(str, "(MISSING)", "", -1)

	// 重新构建json格式数据
	jsonStr2 := []byte(str)

	var deploymentlist []v1beta1.Deployment
	json.Unmarshal(jsonStr2, &deploymentlist)

	c.Data["deploymentlist"] = deploymentlist
	c.Layout = "layout.html"
	c.TplName = "deployment.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}
