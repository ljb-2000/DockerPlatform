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

func (this *BuildController) BuildConfigForGet() {
	id := this.GetString("id")

	// get bitbucket msg
	idtoint, _ := strconv.Atoi(id)
	o := orm.NewOrm()

	//获取版本库信息
	repos := &models.Bitbucket{Id: idtoint}
	o.Read(repos)

	//获取 构建信息
	var pipeline models.Pipeline
	o.QueryTable("pipeline").Filter("Bitbucket", idtoint).One(&pipeline)

	name := ""
	version := ""
	buildmsg := ""
	dockerfilemsg := ""
	if len(pipeline.Name) != 0 {
		name = pipeline.Name
	} else {
		name = ""
	}
	if len(pipeline.Version) != 0 {
		version = pipeline.Version
	} else {
		version = ""
	}
	if len(pipeline.BuildMsg) != 0 {
		buildmsg = pipeline.BuildMsg
	} else {
		buildmsg = "mvn package"
	}
	if len(pipeline.Dockerfile) != 0 {
		dockerfilemsg = pipeline.Dockerfile
	} else {
		dockerfilemsg = ""
	}

	this.Data["id"] = id
	this.Data["repos"] = repos
	this.Data["name"] = name
	this.Data["version"] = version
	this.Data["buildmsg"] = buildmsg
	this.Data["dockerfilemsg"] = dockerfilemsg
	this.Layout = "layout.html"
	this.TplName = "buildconfig.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"

}

func (this *BuildController) BuildConfigForPost() {
	id := this.GetString("id")
	idtoint, _ := strconv.Atoi(id)

	name := this.GetString("name")
	version := this.GetString("version")
	buildcmd := this.GetString("buildcmd")
	dockerfilemsg := this.GetString("dockerfilemsg")

	o := orm.NewOrm()
	bitbucket := models.Bitbucket{Id: idtoint}
	_ = o.Read(&bitbucket)

	pipeline := models.Pipeline{
		Name:       name,
		Version:    version,
		BuildMsg:   buildcmd,
		Dockerfile: dockerfilemsg,
		Bitbucket:  &bitbucket}
	_, err := o.Insert(&pipeline)
	if err != nil {
		var pipeline models.Pipeline
		o.QueryTable("pipeline").Filter("Bitbucket", idtoint).One(&pipeline)
		pipeline.Name = name
		pipeline.Version = version
		pipeline.BuildMsg = buildcmd
		pipeline.Dockerfile = dockerfilemsg
		o.Update(&pipeline)

	}

	this.Redirect("/build.html?id="+id, 302)

}

//加入构建列表
func (this *BuildController) Post() {
	username := this.GetSession("username").(string)
	projectname := this.GetString("projectname")
	name := this.GetString("name")
	url := this.GetString("url")
	url = strings.Replace(url, " ", "", -1)
	url = strings.Replace(url, "\n", "", -1)

	o := orm.NewOrm()
	user := models.User{Name: username}
	_ = o.Read(&user, "Name")

	bitbucket := models.Bitbucket{
		Project:  projectname,
		RepoName: strings.ToLower(name),
		Url:      url,
		User:     &user}

	o.Insert(&bitbucket)

	this.Data["json"] = map[string]string{"status": "200"}
	this.ServeJSON()
}

//从构建列表删除
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

// 显示构建列表内容
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
	orm.NewOrm().QueryTable("bitbucket").Filter("User", user.Id).All(&repos)

	this.Data["repos"] = repos
	this.Layout = "layout.html"
	this.TplName = "buildlist.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

// 进入单独构建页面 ip:port/build.html?id=1
func (this *BuildController) Get() {
	datadir := beego.AppConfig.String("datadir")
	username := this.GetSession("username").(string)
	id := this.GetString("id")
	logfilename := this.GetString("logfilename", "")

	bitid, _ := strconv.Atoi(id)
	o := orm.NewOrm()
	repos := models.Bitbucket{Id: bitid}
	o.Read(&repos)

	// 获取代码分支信息
	// http://bitbucket.gqihome.com/rest/api/1.0/projects/CASH/repos/cashloan-fe/branches
	// Response:
	// {
	//     "size": 20,
	//     "limit": 25,
	//     "isLastPage": true,
	//     "values": [
	//         {
	//             "id": "refs/heads/bugfix/hotfix-0802",
	//             "displayId": "bugfix/hotfix-0802",
	//             "type": "BRANCH",
	//             "latestCommit": "4252ca188434678dbef2979d63a1d7ee37b4a29d",
	//             "latestChangeset": "4252ca188434678dbef2979d63a1d7ee37b4a29d",
	//             "isDefault": false
	//         },
	//         ...
	//     ],
	//     "start": 0
	// }
	bitbucketapiurl := beego.AppConfig.String("bitbucketapiurl")
	bitbucketuser := this.GetSession("username").(string)
	bitbucketpass := GetPasswd(bitbucketuser)

	projectname := repos.Project
	reponame := repos.RepoName
	url := bitbucketapiurl + "/" + filepath.Join("projects", projectname, "repos", reponame, "branches")
	result, _ := RequestForAuth("GET", url, bitbucketuser, bitbucketpass, nil)
	var branches interface{}
	json.Unmarshal(result, &branches)

	// 判断pipeline 表是否有该id的数据
	// @warning: true build.html页面显示警告信息
	// @warning: false 不显示
	var pipeline models.Pipeline
	var warning bool
	err := o.QueryTable("pipeline").Filter("Bitbucket", bitid).One(&pipeline)
	if err != nil {
		warning = true
	} else {
		warning = false
	}

	// 获取构建日志名，并追加到临时列表
	// reponame := repos.RepoName
	builddir := filepath.Join(datadir, "pipelogs", username, reponame)
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
	this.Data["branches"] = branches
	this.Data["logfilelist"] = logfilelist
	this.Data["warning"] = warning
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
	branches := this.GetString("branches")
	branches = strings.TrimSpace(branches)
	logfile := this.GetString("logfile")
	id := this.GetString("id")
	username := this.GetSession("username").(string)

	//获取应用名
	arr1 := strings.Split(giturl, "/")
	arr2 := strings.Split(arr1[len(arr1)-1], ".")
	projectname := arr2[0]

	// 记录构建日志
	msg := make(map[string]string)
	f, _ := ioutil.ReadFile(logfile)
	json.Unmarshal(f, &msg)

	// 切换分支 & 构建代码
	ckeckoutcmd := "git checkout " + branches
	clonedir := filepath.Join(datadir, "gitcode", username, projectname)
	os.Chdir(clonedir)
	idtoint, _ := strconv.Atoi(id)
	o := orm.NewOrm()
	var pipeline models.Pipeline
	o.QueryTable("pipeline").Filter("Bitbucket", idtoint).One(&pipeline)
	buildcmd := pipeline.BuildMsg
	buildcmd = strings.TrimSpace(buildcmd)
	buildcmd = strings.Replace(buildcmd, "\n", "&&", -1)
	cmd := ckeckoutcmd + "&&" + "git pull" + "&&" + buildcmd

	if strings.Contains(cmd, "rm") {
		this.Data["json"] = map[string]string{"status": "300", "err": string("请勿输入危险命令: " + buildcmd)}
		this.ServeJSON()
	}

	err, out, _ := Shellout(cmd)

	//返回json信息
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
	id := this.GetString("id")
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

	// 获取应用名称和版本
	idtoint, _ := strconv.Atoi(id)
	o := orm.NewOrm()
	var pipeline models.Pipeline
	o.QueryTable("pipeline").Filter("Bitbucket", idtoint).One(&pipeline)
	name := pipeline.Name
	version := pipeline.Version

	// 定义克隆地址
	clonedir := filepath.Join(datadir, "gitcode", username, projectname)

	// 如果构建数据库中没有版本，则从pom文件中获取
	if len(version) == 0 {
		pomfile := clonedir + "/pom.xml"
		pommsg := GetPomMsg(pomfile)
		version = pommsg.Version
	}

	// 默认以用户名为harbor项目名
	harborpro := username
	if !CheckHarborProExist(harborpro) {
		CreateHarborPro(harborpro)
	}

	// 进入代码目录，根据dockerfile进行打包
	os.Chdir(clonedir)

	//如果doackerfile不存在，返回错误信息
	dokerfile := filepath.Join(clonedir, "Dockerfile")
	if _, err := os.Stat(dokerfile); err != nil {
		this.Data["json"] = map[string]string{"status": "300", "err": string(err.Error())}
		this.ServeJSON()
	}

	// 打包并上传harbor仓库
	image := "harbor.gqichina.com/" + harborpro + "/" + name + ":" + version + "-" + time.Now().Format("20060102-150405")
	imagelatst := "harbor.gqichina.com/" + harborpro + "/" + name + ":latest\n"
	cmdtobuild := "docker build -t " + image + " --build-arg APP_VERSION=" + version + " ./\n"
	cmdtotag := "docker tag " + image + " " + imagelatst
	cmdtopushimage := "docker push " + image + "\n"
	cmdtopushimagelatst := "docker push " + imagelatst + "\n"
	cmdrmi := "docker rmi " + image + "&&" + "docker rmi " + imagelatst

	// pomfile := clonedir + "/pom.xml"
	// os.Chdir(clonedir)

	// pommsg := GetPomMsg(pomfile)
	// version := pommsg.Version
	// artifactid := pommsg.ArtifactId
	// app := "alpine-" + artifactid

	// 默认以用户名为harbor项目名
	// harborpro := username
	// if !CheckHarborProExist(harborpro) {
	// 	CreateHarborPro(harborpro)
	// }

	// image := "harbor.gqichina.com/" + harborpro + "/" + app + ":" + version + "-" + time.Now().Format("20060102-150405")
	// imagelatst := "harbor.gqichina.com/" + harborpro + "/" + app + ":latest\n"
	// cmdtobuild := "docker build -t " + image + " --build-arg APP_VERSION=" + version + " ./\n"
	// cmdtotag := "docker tag " + image + " " + imagelatst
	// cmdtopushimage := "docker push " + image + "\n"
	// cmdtopushimagelatst := "docker push " + imagelatst

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
		in.WriteString(cmdrmi)
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
