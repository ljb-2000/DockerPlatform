package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	v1 "k8s.io/client-go/pkg/api/v1"
	"net/http"
)

type NamespacesController struct {
	beego.Controller
}

func (c *NamespacesController) Get() {
	resp, err := http.Get("http://10.10.7.175:8081/namespaces/list")
	if err != nil {
		fmt.Println(err)
	}

	result, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var namespacelist []v1.Namespace
	json.Unmarshal([]byte(result), &namespacelist)

	c.Data["namespaceslist"] = namespacelist

	c.Layout = "layout.html"
	c.TplName = "namespaces.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}
