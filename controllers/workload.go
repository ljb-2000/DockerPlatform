package controllers

import (
	"encoding/json"
	// "github.com/astaxie/beego"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"strings"
)

type WorkloadController struct {
	baseControllers
}

func (this *WorkloadController) Get() {
	username := this.GetSession("username").(string)
	namespaces := ""
	if CheckAdmin(username) {
		namespaces = this.GetString("namespaces")
	} else {
		namespaces = username
	}

	/// Get Deployment
	deploymenturl := "http://10.10.7.175:8081/deployment/list"
	//json 序列化
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)
	//post数据
	deploymenresult, _ := Request("POST", deploymenturl, jsonStr)
	// json数据错误，去掉(MISSING)字符串
	str := string(deploymenresult)
	str = strings.Replace(str, "(MISSING)", "", -1)
	// 重新构建json格式数据
	jsonStr2 := []byte(str)
	var deploymentlist []v1beta1.Deployment
	json.Unmarshal(jsonStr2, &deploymentlist)

	/// Get services
	serviceurl := "http://10.10.7.175:8081/services/list"
	serviceresult, _ := Request("POST", serviceurl, jsonStr)
	var servicelist []v1.Service
	json.Unmarshal(serviceresult, &servicelist)

	///Get Pods
	podsurl := "http://10.10.7.175:8081/pods/list"
	//post数据
	podresult, _ := Request("POST", podsurl, jsonStr)
	var podslist []v1.Pod
	json.Unmarshal([]byte(podresult), &podslist)

	this.Data["deploymentlist"] = deploymentlist
	this.Data["servicelist"] = servicelist
	this.Data["podslist"] = podslist
	this.Layout = "layout.html"
	this.TplName = "workload.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
