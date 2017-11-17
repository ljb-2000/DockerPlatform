package controllers

import (
	"DockerPlatform/models"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type BuildController struct {
	baseControllers
}

func (this *BuildController) Post() {
	username := this.GetSession("username").(string)
	name := this.GetString("name")
	url := this.GetString("url")
	url = strings.Replace(url, " ", "", -1)
	url = strings.Replace(url, "\n", "", -1)

	o := orm.NewOrm()
	user := models.User{Name: username}
	_ = o.Read(&user, "Name")
	bitbucket := models.Bitbucket{RepoName: strings.ToLower(name), Url: url, User: &user}
	o.Insert(&bitbucket)

	this.Data["json"] = map[string]string{"status": "200"}
	this.ServeJSON()
}

func (this *BuildController) RepoDelete() {
	repoid := this.GetString("repoid")

	id, _ := strconv.Atoi(repoid)
	o := orm.NewOrm()
	if _, err := o.Delete(&models.Bitbucket{Id: id}); err == nil {
		this.Data["json"] = map[string]string{"status": "200"}
	} else {
		this.Data["json"] = map[string]string{"status": "300"}
	}

	this.ServeJSON()
}

func (this *BuildController) List() {
	username := this.GetString("username")

	if len(username) == 0 {
		username = this.GetSession("username").(string)
	}

	// username := this.GetSession("username").(string)

	o := orm.NewOrm()
	user := models.User{Name: username}
	_ = o.Read(&user, "Name")

	var repos []*models.Bitbucket
	orm.NewOrm().QueryTable("bitbucket").Filter("User", user.Id).RelatedSel().All(&repos)

	this.Data["repos"] = repos
	this.Layout = "layout.html"
	this.TplName = "buildlist.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func (this *BuildController) Get() {
	datadir := beego.AppConfig.String("datadir")
	username := this.GetSession("username").(string)
	id := this.GetString("id")
	logfilename := this.GetString("logfilename", "")

	bitid, _ := strconv.Atoi(id)
	o := orm.NewOrm()
	repos := models.Bitbucket{Id: bitid}
	o.Read(&repos)

	projectname := repos.RepoName

	builddir := filepath.Join(datadir, "pipelogs", username, projectname)

	logfile, _ := ioutil.ReadDir(builddir)

	var logfilelist []string

	for _, v := range logfile {
		logfilelist = append(logfilelist, v.Name())
	}

	if len(logfilename) != 0 {
		logfulldir := builddir + "/" + logfilename
		msg := make(map[string]string)
		f, _ := ioutil.ReadFile(logfulldir)
		json.Unmarshal(f, &msg)
		this.Data["log"] = msg
	} else {
		this.Data["log"] = ""
	}

	this.Data["id"] = repos.Id
	this.Data["repos"] = repos
	this.Data["logfilelist"] = logfilelist
	this.Layout = "layout.html"
	this.TplName = "build.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func (this *BuildController) GetForAdmin() {
	username := this.GetSession("username").(string)

	o := orm.NewOrm()
	qs := o.QueryTable("user")

	var usermsg map[string]interface{}
	usermsg = make(map[string]interface{})

	var list []orm.Params
	qs.Exclude("name", username).Values(&list)
	for _, v := range list {
		name := v["Name"]
		count, _ := orm.NewOrm().QueryTable("bitbucket").Filter("User", v["Id"]).Count()
		usermsg[name.(string)] = count
	}

	this.Data["usermsg"] = usermsg
	this.Layout = "layout.html"
	this.TplName = "buildforadmin.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"

}

func (this *BuildController) PipelineToClone() {
	datadir := beego.AppConfig.String("datadir")
	giturl := this.GetString("giturl")
	username := this.GetSession("username").(string)
	password := GetPasswd(username)

	password = ":" + password + "@"

	arr := strings.Split(giturl, "@")
	giturl = arr[0] + password + arr[1]

	arr1 := strings.Split(giturl, "/")
	arr2 := strings.Split(arr1[len(arr1)-1], ".")
	projectname := arr2[0]

	// 记录克隆日志
	timestring := time.Now().Format("20060102-150405")
	file := timestring

	builddir := filepath.Join(datadir, "pipelogs", username, projectname)
	MkdirIfNoExist(builddir)

	logfile := filepath.Join(builddir, file)
	log := make(map[string]string)

	os.Chdir(datadir)

	// 创建克隆地址
	clonedir := filepath.Join(datadir, "gitcode", username, projectname)

	if _, err := os.Stat(clonedir); err == nil {
		os.Chdir(clonedir)
		err, stdout, stderr := Shellout("git pull")
		if err != nil {
			this.Data["json"] = map[string]string{"status": "300", "err": stderr}
		} else {
			this.Data["json"] = map[string]string{"status": "200", "msg": stdout, "logfile": logfile}
		}
		log["clone"] = string(stdout)
		mjsonstring, _ := json.Marshal(log)
		ioutil.WriteFile(logfile, mjsonstring, 0666)
	} else {
		cmdstring := "git clone " + giturl + " " + clonedir
		err, stdout, stderr := Shellout(cmdstring)
		if err != nil {
			this.Data["json"] = map[string]string{"status": "300", "err": stderr}
		} else {
			this.Data["json"] = map[string]string{"status": "200", "msg": stdout, "logfile": logfile}
		}
		log["clone"] = string(stdout)
		mjsonstring, _ := json.Marshal(log)
		ioutil.WriteFile(logfile, mjsonstring, 0666)
	}

	this.ServeJSON()
}

func (this *BuildController) PipelineToBuild() {
	datadir := beego.AppConfig.String("datadir")
	giturl := this.GetString("giturl")
	logfile := this.GetString("logfile")
	username := this.GetSession("username").(string)
	arr1 := strings.Split(giturl, "/")
	arr2 := strings.Split(arr1[len(arr1)-1], ".")
	projectname := arr2[0]

	// tmpdir := "/tmp" + username
	// os.Mkdir(tmpdir, 0777)

	// clonedir := "/tmp/" + username + "/" + projectname
	clonedir := filepath.Join(datadir, "gitcode", username, projectname)

	// 记录构建日志
	msg := make(map[string]string)
	f, _ := ioutil.ReadFile(logfile)
	json.Unmarshal(f, &msg)

	// pipeline -> build
	os.Chdir(clonedir)
	err, out, _ := Shellout("mvn package")

	if err != nil {
		this.Data["json"] = map[string]string{"status": "300", "err": string(out)}
	} else {
		this.Data["json"] = map[string]string{"status": "200", "msg": string(out), "logfile": logfile}
	}
	msg["build"] = string(out)
	mjsonstring, _ := json.Marshal(msg)
	ioutil.WriteFile(logfile, mjsonstring, 0666)
	this.ServeJSON()
}

type PomResult struct {
	Version    string `xml:"version"`
	ArtifactId string `xml:"artifactId"`
}

func GetPomMsg(pomfile string) *PomResult {
	file, err := os.Open(pomfile) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	v := PomResult{}
	err = xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	return &v
}

func (this *BuildController) PipelineToPush() {
	datadir := beego.AppConfig.String("datadir")
	giturl := this.GetString("giturl")
	logfile := this.GetString("logfile")
	username := this.GetSession("username").(string)
	arr1 := strings.Split(giturl, "/")
	arr2 := strings.Split(arr1[len(arr1)-1], ".")
	projectname := arr2[0]

	// tmpdir := "/tmp" + username
	// os.Mkdir(tmpdir, 0777)

	// 记录push日志
	msg := make(map[string]string)
	f, _ := ioutil.ReadFile(logfile)
	json.Unmarshal(f, &msg)

	// clonedir := "/tmp/" + username + "/" + projectname
	clonedir := filepath.Join(datadir, "gitcode", username, projectname)

	pomfile := clonedir + "/pom.xml"

	os.Chdir(clonedir)

	pommsg := GetPomMsg(pomfile)
	version := pommsg.Version
	artifactid := pommsg.ArtifactId
	app := "alpine-" + artifactid

	// 默认以用户名为harbor项目名
	harborpro := username
	if !CheckHarborProExist(harborpro) {
		CreateHarborPro(harborpro)
	}

	image := "harbor.gqichina.com/" + harborpro + "/" + app + ":" + version + "-" + time.Now().Format("20060102-150405")
	imagelatst := "harbor.gqichina.com/" + harborpro + "/" + app + ":latest\n"
	cmdtobuild := "docker build -t " + image + " --build-arg APP_VERSION=" + version + " ./\n"
	cmdtotag := "docker tag " + image + " " + imagelatst
	cmdtopushimage := "docker push " + image + "\n"
	cmdtopushimagelatst := "docker push " + imagelatst

	in := bytes.NewBuffer(nil)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd := exec.Command("/usr/bin/bash")
	cmd.Stdin = in
	cmd.Stdout = &out
	cmd.Stderr = &err
	go func() {
		in.WriteString(cmdtobuild)
		in.WriteString(cmdtotag)
		in.WriteString(cmdtopushimage)
		in.WriteString(cmdtopushimagelatst)
	}()
	cmd.Run()

	if len(err.String()) != 0 {
		this.Data["json"] = map[string]string{"status": "300", "err": err.String()}
		msg["push"] = string(err.String())
		mjsonstring, _ := json.Marshal(msg)
		ioutil.WriteFile(logfile, mjsonstring, 0666)
	} else {
		this.Data["json"] = map[string]string{"status": "200", "msg": out.String()}
		msg["push"] = string(out.String())
		mjsonstring, _ := json.Marshal(msg)
		ioutil.WriteFile(logfile, mjsonstring, 0666)
	}

	this.ServeJSON()
}
