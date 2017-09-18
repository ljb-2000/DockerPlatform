package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	// v1 "k8s.io/client-go/pkg/api/v1"
	"net/http"
)

type ImageController struct {
	beego.Controller
}

func (c *ImageController) Get() {
	url := "http://10.10.7.101/api/projects"

	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("tanzhixu", "1QAZ2wsx")

	var resp *http.Response
	resp, _ = http.DefaultClient.Do(request)

	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(result))

	var harborprojectlist []HarborProjects
	json.Unmarshal(result, &harborprojectlist)

	// fmt.Println(harborprojectlist)

	c.Data["harborprojectlist"] = harborprojectlist
	c.Layout = "layout.html"
	c.TplName = "image.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}
