package controllers

import (
	"encoding/json"
	"github.com/bndr/gojenkins"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"strconv"
	"strings"
)

type KubernetesController struct {
	baseControllers
}

func (this *KubernetesController) INIT() {
	username := this.GetSession("username").(string)

	jenkinsurl := "http://d.ci.gqichina.com"

	jk, _ := gojenkins.CreateJenkins(nil, jenkinsurl).Init()

	var parameters map[string]string
	parameters = make(map[string]string)
	parameters["NAMESPACES"] = username

	jk.BuildJob("INIT", parameters)

	this.Data["json"] = map[string]string{"status": "200"}
	this.ServeJSON()
}

func (this *KubernetesController) GetKubenetesMsg() {
	username := this.GetSession("username").(string)
	namespaces := ""
	if CheckAdmin(username) {
		namespaces = this.GetString("namespaces")
	} else {
		namespaces = username
	}

	//json 序列化
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	url := "http://10.10.7.175:8081/deployment/list"
	deploymentresult, _ := Request("POST", url, jsonStr)
	str := string(deploymentresult)
	str = strings.Replace(str, "(MISSING)", "", -1)
	jsonStr2 := []byte(str)
	var deploymentlist []v1beta1.Deployment
	json.Unmarshal(jsonStr2, &deploymentlist)

	url = "http://10.10.7.175:8081/services/list"
	servicesresult, _ := Request("POST", url, jsonStr)
	var servicelist []v1.Service
	json.Unmarshal(servicesresult, &servicelist)

	url = "http://10.10.7.175:8081/pods/list"
	podsresult, _ := Request("POST", url, jsonStr)
	var podslist []v1.Pod
	json.Unmarshal(podsresult, &podslist)

	msg := []map[string]string{}
	msg = append(msg,
		map[string]string{"name": "Deployments", "value": strconv.Itoa(len(deploymentlist))},
		map[string]string{"name": "Services", "value": strconv.Itoa(len(servicelist))},
		map[string]string{"name": "Pods", "value": strconv.Itoa(len(podslist))})

	data := make(map[string]interface{})
	data["status"] = "200"
	data["msg"] = msg

	this.Data["json"] = data
	this.ServeJSON()
}
