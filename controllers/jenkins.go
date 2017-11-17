package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/bndr/gojenkins"
	"strconv"
	"strings"
	"time"
)

type JenkinsSet struct {
	jk *gojenkins.Jenkins
}

func (j *JenkinsSet) ConnJenkins() *gojenkins.Jenkins {
	jenkinsurl := beego.AppConfig.String("jenkinsurl")
	jenkinsuser := beego.AppConfig.String("jenkinsuser")
	jenkinspass := beego.AppConfig.String("jenkinspass")
	if len(jenkinsuser) != 0 {
		j.jk, _ = gojenkins.CreateJenkins(nil, jenkinsurl, jenkinsuser, jenkinspass).Init()
		return j.jk
	} else {
		j.jk, _ = gojenkins.CreateJenkins(nil, jenkinsurl).Init()
		return j.jk
	}
}

type JenkinsController struct {
	beego.Controller
	jk JenkinsSet
}

func (c *JenkinsController) GetAllJobs() {
	jenkins := c.jk.ConnJenkins()

	jobs, _ := jenkins.GetAllJobs()

	c.Data["jkjoslist"] = jobs
	c.Layout = "layout.html"
	c.TplName = "jenkins.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}

type BuildResult struct {
	Result    string    `json:result`
	Number    int64     `json:number`
	URL       string    `json:url`
	Timestamp time.Time `json:"timestamp"`
}

func (c *JenkinsController) GetAllBuildIds() {
	jobname := c.GetString("jobname", "")

	jenkins := c.jk.ConnJenkins()
	buildids, _ := jenkins.GetAllBuildIds(jobname)

	var buildlist []BuildResult

	for _, v := range buildids {
		b, _ := jenkins.GetBuild(jobname, v.Number)
		buildlist = append(buildlist, BuildResult{b.GetResult(), v.Number, v.URL, b.GetTimestamp()})
	}

	fmt.Println(buildlist)
	c.Data["buildids"] = buildlist
	c.Data["jobname"] = jobname
	c.Layout = "layout.html"
	c.TplName = "jkbuildids.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}

func (c *JenkinsController) GetConsoleOutput() {
	jobname := c.GetString("jobname", "")
	jobid := c.GetString("jobid", "")

	jobidtoint64, _ := strconv.ParseInt(jobid, 10, 64)

	jenkins := c.jk.ConnJenkins()
	buildresp, _ := jenkins.GetBuild(jobname, jobidtoint64)

	buildstr := buildresp.GetConsoleOutput()
	buildstr = strings.Replace(buildstr, "\n", "<br>", -1)
	buildstr = strings.Replace(buildstr, "\t", "&emsp;&emsp;", -1)

	c.Data["buildresp"] = buildstr
	c.Layout = "layout.html"
	c.TplName = "jkbuildresp.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = "html_head.html"
	c.LayoutSections["BodyHead"] = "body_head.html"
	c.LayoutSections["Sidebar"] = "sidebar.html"
}

func (c *JenkinsController) Build() {
	jobname := c.GetString("jobname")

	jenkins := c.jk.ConnJenkins()
	a, _ := jenkins.BuildJob(jobname)
	fmt.Println(a)

	c.Data["json"] = map[string]string{"status": "200"}
	c.ServeJSON()
}
