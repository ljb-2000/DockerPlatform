package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
)

type RepoController struct {
	baseControllers
}

type Repository struct {
	Name          string
	DownloadCount int
}

func (this *RepoController) Get() {
	username := this.GetSession("username").(string)
	harborhost := beego.AppConfig.String("harborhost")
	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")

	if CheckAdmin(username) {
		projectid := this.GetString("projectid", "")
		url := harborhost + "/api/repositories?project_id=" + projectid
		result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)
		var f []interface{}
		json.Unmarshal(result, &f)
		repos := make(map[string]interface{})
		for _, v := range f {
			repos[v.(string)] = GetDownloadCount(v.(string))
		}
		this.Data["projectid"] = projectid
		this.Data["repositories"] = repos
	} else {
		projectid := GetHarborProId(username)
		if projectid > 0 {
			url := harborhost + "/api/repositories?project_id=" + strconv.FormatFloat(projectid, 'f', -1, 64)
			result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)
			var f []interface{}
			json.Unmarshal(result, &f)
			repos := make(map[string]interface{})
			for _, v := range f {
				repos[v.(string)] = GetDownloadCount(v.(string))
			}
			this.Data["projectid"] = projectid
			this.Data["repositories"] = repos
		}

	}

	this.Layout = "layout.html"
	this.TplName = "repos.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func (this *RepoController) GetAll() {
	harborhost := beego.AppConfig.String("harborhost")
	url := harborhost + "/api/projects"
	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")

	result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)

	var harborprojectlist []HarborProjects
	json.Unmarshal(result, &harborprojectlist)

	this.Data["harborprojectlist"] = harborprojectlist
	this.Layout = "layout.html"
	this.TplName = "image.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func GetDownloadCount(repo string) float64 {
	harborhost := beego.AppConfig.String("harborhost")

	url := harborhost + "/api/repositories/top?count=1000"
	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")

	result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)

	var f []interface{}
	json.Unmarshal(result, &f)
	for _, v := range f {
		md, _ := v.(map[string]interface{})
		if md["name"] == repo {
			return md["count"].(float64)
		}
	}
	return 0
}
