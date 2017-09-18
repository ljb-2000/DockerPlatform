package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type PodsController struct {
	beego.Controller
}

func (c *PodsController) Get() {
	namespaces := c.GetString("namespaces", "")
	url := "http://10.10.7.175:8081/pods/list"

	//json 序列化
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	//post数据
	result, _ := DoPost(url, jsonStr)

	var podslist []v1.Pod
	json.Unmarshal([]byte(result), &podslist)

	c.Data["podslist"] = podslist

	c.Layout = "layout.html"
	c.TplName = "pods.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}
