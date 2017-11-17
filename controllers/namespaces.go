package controllers

import (
	"encoding/json"
	// "github.com/astaxie/beego"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type NamespacesController struct {
	baseControllers
}

func (this *NamespacesController) Get() {
	url := "http://10.10.7.175:8081/namespaces/list"

	result, _ := Request("GET", url, nil)

	var namespacelist []v1.Namespace
	json.Unmarshal([]byte(result), &namespacelist)

	this.Data["namespaceslist"] = namespacelist

	this.Layout = "layout.html"
	this.TplName = "namespaces.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
