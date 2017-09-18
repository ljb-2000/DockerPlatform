package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	// "reflect"
	// "io/ioutil"
	// v1 "k8s.io/client-go/pkg/api/v1"
	// "net/http"
)

type RepoController struct {
	beego.Controller
}

type Repository struct {
	Name string
}

func (c *RepoController) Get() {
	projectid := c.GetString("projectid", "")
	url := "http://10.10.7.101/api/repositories?project_id=" + projectid

	result, _ := Request("GET", url, nil)

	var f []interface{}
	json.Unmarshal(result, &f)

	s := make([]string, len(f))
	for i, v := range f {
		s[i] = fmt.Sprint(v)
	}
	fmt.Println(s)

	// repos := []*Repository{}
	// for _, v := range f {
	// 	var re1 = new(Repository)
	// 	re1.Name = v.(string)
	// 	repos = append(repos, re1)
	// }
	// fmt.Println(repos)
	// for _, v := range repos {
	// 	fmt.Println(&v)
	// }

	c.Data["repositories"] = f
	c.Layout = "layout.html"
	c.TplName = "repos.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}
