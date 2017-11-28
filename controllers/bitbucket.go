package controllers

import (
	"DockerPlatform/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"path/filepath"
	"strconv"
)

type BitbucketController struct {
	baseControllers
}

func (this *BitbucketController) Get() {
	username := this.GetSession("username").(string)

	bitbucketapiurl := beego.AppConfig.String("bitbucketapiurl")
	bitbucketuser := username
	bitbucketpass := GetPasswd(bitbucketuser)

	url := bitbucketapiurl + "/projects?limit=100"
	result, _ := RequestForAuth("GET", url, bitbucketuser, bitbucketpass, nil)

	var bitbucketlist interface{}
	json.Unmarshal(result, &bitbucketlist)

	this.Data["bitbucketlist"] = bitbucketlist
	this.Layout = "layout.html"
	this.TplName = "bitbucket.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

//获取gitlab 中commits信息
func (this *BitbucketController) GetCommits() {
	id := this.GetString("id")

	bitbucketapiurl := beego.AppConfig.String("bitbucketapiurl")
	bitbucketuser := this.GetSession("username").(string)
	bitbucketpass := GetPasswd(bitbucketuser)

	//根据id获取bitbucket信息
	// var bb bitbucket
	// repos:
	// {5 CASH cashloan-fe http://tanzhixu@bitbucket.gqihome.com/scm/cash/cashloan-fe.git 0xc420248870 <nil>}
	bitid, _ := strconv.Atoi(id)
	o := orm.NewOrm()
	bb := models.Bitbucket{Id: bitid}
	o.Read(&bb)

	// 获取代码提交信息
	// hhttp://bitbucket.gqihome.com/rest/api/1.0/projects/CASH/repos/cashloan-fe/commits
	// Response:
	// {
	//     "values": [
	//         {
	//             "id": "344d799e2f83ada9b4f4caf9501f6bc8682c0fca",
	//             "displayId": "344d799e2f8",
	//             "author": {
	//                 "name": "wangyongqing",
	//                 "emailAddress": "wangyongqing@gqhmt.com",
	//                 "id": 468,
	//                 "displayName": "王 永庆",
	//                 "active": true,
	//                 "slug": "wangyongqing",
	//                 "type": "NORMAL",
	//                 "links": {
	//                     "self": [
	//                         {
	//                             "href": "http://bitbucket.gqihome.com/users/wangyongqing"
	//                         }
	//                     ]
	//                 }
	//             },
	//             "authorTimestamp": 1511488508000,
	//             "message": "Merge pull request #364 in CASH/cashloan-fe from feature/release-08 to master\n\n* commit 'a15d2d7b6e9fb67d0ee61e3fe4f9cf829a23b73c': (108 commits)\n  [maven-release-plugin] prepare for next development iteration\n  [maven-release-plugin] prepare release 1.2.2\n  docker\n  docker\n  [maven-release-plugin] prepare for next development iteration\n  [maven-release-plugin] prepare release 1.2.1\n  修改图片\n  [maven-release-plugin] prepare for next development iteration\n  [maven-release-plugin] prepare release 1.2.0\n  还款失败\n  样式修改回来\n  还款失败\n  修改样式，测试APP\n  还款失败\n  [maven-release-plugin] prepare for next development iteration\n  [maven-release-plugin] prepare release 1.1.9\n  与APP交互，之前交接有误，修改文案\n  还款失败\n  还款失败\n  还款\n  ...",
	//             "parents": [
	//                 {
	//                     "id": "c7f39a3caad068bb12a748228b12f3fd9f13e611",
	//                     "displayId": "c7f39a3caad"
	//                 },
	//                 {
	//                     "id": "a15d2d7b6e9fb67d0ee61e3fe4f9cf829a23b73c",
	//                     "displayId": "a15d2d7b6e9"
	//                 }
	//             ]
	//         },
	//         ...
	//     ],
	//     "size": 25,
	//     "isLastPage": false,
	//     "start": 0,
	//     "limit": 25,
	//     "nextPageStart": 25
	// }

	projectname := bb.Project
	reponame := bb.RepoName
	url := bitbucketapiurl + "/" + filepath.Join("projects", projectname, "repos", reponame, "commits")
	result, _ := RequestForAuth("GET", url, bitbucketuser, bitbucketpass, nil)
	var commits interface{}
	json.Unmarshal(result, &commits)

	this.Data["commits"] = commits
	this.Layout = "layout.html"
	this.TplName = "commits.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

func (this *BitbucketController) BitbucketRepos() {
	bitbucketapiurl := beego.AppConfig.String("bitbucketapiurl")
	bitbucketuser := this.GetSession("username").(string)
	bitbucketpass := GetPasswd(bitbucketuser)

	project := this.GetString("project")
	url := bitbucketapiurl + "/projects/" + project + "/repos"
	result, _ := RequestForAuth("GET", url, bitbucketuser, bitbucketpass, nil)

	var bitbucketrepos interface{}
	json.Unmarshal(result, &bitbucketrepos)

	this.Data["bitbucketrepos"] = bitbucketrepos
	this.Data["projectname"] = project
	this.Layout = "layout.html"
	this.TplName = "bitbucketrepos.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
