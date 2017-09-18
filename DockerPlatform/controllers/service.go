package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type ServiceController struct {
	beego.Controller
}

func (c *ServiceController) Get() {
	url := "http://10.10.7.175:8081/services/list"

	result, _ := DoGet(url)

	var servicelist []v1.Service
	json.Unmarshal(result, &servicelist)

	c.Data["servicelist"] = servicelist
	c.Layout = "layout.html"
	c.TplName = "service.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}
