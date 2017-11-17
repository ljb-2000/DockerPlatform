package controllers

import (
	"DockerPlatform/models"
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	v1 "k8s.io/client-go/pkg/api/v1"
	"path/filepath"
	"strconv"
)

type DashboardController struct {
	baseControllers
}

func (this *DashboardController) Get() {
	username := this.GetSession("username").(string)
	if CheckAdmin(username) {
		this.Data["nothavenamespace"] = ""
	} else {
		url := "http://10.10.7.175:8081/namespaces/list"

		result, _ := Request("GET", url, nil)

		var namespacelist []v1.Namespace
		json.Unmarshal([]byte(result), &namespacelist)

		namelist := []string{}
		for _, v := range namespacelist {
			namelist = append(namelist, v.GetName())
		}

		if stringInSlice(username, namelist) {
			this.Data["nothavenamespace"] = ""
		} else {
			var str bytes.Buffer
			str.WriteString("<p>暂时无Kubenetes环境，请初始化环境</p><button class=\"")
			str.WriteString("btn btn-success btn-xs init\"")
			str.WriteString(">初始化</button>")
			this.Data["nothavenamespace"] = str.String()
		}
	}

	this.Layout = "layout.html"
	this.TplName = "dashboard.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func (this *DashboardController) GetBuildMsg() {
	datadir := beego.AppConfig.String("datadir")
	username := this.GetSession("username").(string)

	o := orm.NewOrm()
	user := models.User{Name: username}
	_ = o.Read(&user, "Name")

	var repos []*models.Bitbucket
	orm.NewOrm().QueryTable("bitbucket").Filter("User", user.Id).RelatedSel().All(&repos)

	var namelist []string
	var valuelist []int

	for _, v := range repos {
		namelist = append(namelist, v.RepoName)
		builddir := filepath.Join(datadir, "pipelogs", username, v.RepoName)
		// builddir := "/tmp/logs/" + username + "/" + v.RepoName
		b, _ := ioutil.ReadDir(builddir)
		valuelist = append(valuelist, len(b))
	}

	data := make(map[string]interface{})
	data["status"] = "200"
	data["namelist"] = namelist
	data["valuelist"] = valuelist

	this.Data["json"] = data
	this.ServeJSON()
}

func (this *DashboardController) GetPodsMsg() {
	username := this.GetSession("username").(string)

	namespaces := ""
	if CheckAdmin(username) {
		namespaces = this.GetString("namespaces")
	} else {
		namespaces = username
	}

	url := "http://10.10.7.175:8081/pods/list"

	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	result, _ := Request("POST", url, jsonStr)

	var podslist []v1.Pod
	json.Unmarshal([]byte(result), &podslist)

	running := 0
	pending := 0

	for _, v := range podslist {
		if v.Status.Phase == "Running" {
			running = running + 1
		}
		if v.Status.Phase == "Pending" {
			pending = pending + 1
		}
	}

	red := map[string]interface{}{"normal": map[string]string{"color": "red"}}
	green := map[string]interface{}{"normal": map[string]string{"color": "green"}}

	msg := []map[string]interface{}{}
	msg = append(msg,
		map[string]interface{}{"name": "Running", "value": strconv.Itoa(running), "itemStyle": green},
		map[string]interface{}{"name": "Pending", "value": strconv.Itoa(pending), "itemStyle": red})

	data := make(map[string]interface{})
	data["status"] = "200"
	data["msg"] = msg

	this.Data["json"] = data
	this.ServeJSON()

}
