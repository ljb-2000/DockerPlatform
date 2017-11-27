package controllers

import (
	"encoding/json"
	// "github.com/astaxie/beego"
	"fmt"
	v1 "k8s.io/client-go/pkg/api/v1"
	// v1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta2 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type PodsController struct {
	baseControllers
}

type Pod struct {
	Name              string
	Namespace         string
	CreationTimestamp metav1.Time
	HostIP            string
	PodIP             string
	Image             string
	Domain            string
	ClusterIP         string
	Ports             []v1.ServicePort
}

func (this *PodsController) GetPodsLogs() {
	name := this.GetString("name")
	namespace := this.GetString("namespace")

	cmd := "/root/local/bin/kubectl logs --tail 1000"
	cmd = fmt.Sprintf("%s %s %s %s", cmd, name, "-n", namespace)

	msg := RemoteCommand(cmd)
	this.Data["json"] = map[string]string{"status": "200", "msg": msg}
	this.ServeJSON()
}

func (this *PodsController) Details() {
	appname := this.GetString("name")
	username := this.GetSession("username").(string)

	namespaces := ""
	if CheckAdmin(username) {
		namespaces = this.GetString("namespaces")
	} else {
		namespaces = username
	}

	pod := new(Pod)

	podmsg := PodMsg(namespaces)
	for _, podvalue := range podmsg {
		if appname == podvalue.ObjectMeta.Labels["app"] {
			pod.Name = podvalue.ObjectMeta.Name
			pod.Namespace = podvalue.ObjectMeta.Namespace
			pod.CreationTimestamp = podvalue.ObjectMeta.CreationTimestamp
			pod.HostIP = podvalue.Status.HostIP
			pod.PodIP = podvalue.Status.PodIP
			pod.Image = podvalue.Spec.Containers[0].Image
		}
	}

	ingmsg := IngressMsg(namespaces)
	for _, ingvalue := range ingmsg {
		if appname == ingvalue.ObjectMeta.Name {
			pod.Domain = ingvalue.Spec.Rules[0].Host
		}
	}

	servicemsg := ServiceMsg(namespaces)
	for _, svcvalue := range servicemsg {
		if appname == svcvalue.ObjectMeta.Name {
			pod.ClusterIP = svcvalue.Spec.ClusterIP
			pod.Ports = svcvalue.Spec.Ports
		}
	}

	// //获取deployment信息
	// depurl := "http://10.10.7.175:8081/deployment/list"
	// depmsg, _ := Request("POST", depurl, jsonStr)
	// // json数据错误，去掉(MISSING)字符串
	// str := string(depmsg)
	// str = strings.Replace(str, "(MISSING)", "", -1)
	// // 重新构建json格式数据
	// jsonStr2 := []byte(str)
	// var deploymentlist []v1beta1.Deployment
	// json.Unmarshal(jsonStr2, &deploymentlist)

	this.Data["appname"] = appname
	this.Data["pod"] = pod
	this.Layout = "layout.html"
	this.TplName = "podsdetails.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}

//获取pod信息
func PodMsg(namespaces string) []v1.Pod {
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	podsurl := "http://10.10.7.175:8081/pods/list"
	podsmsg, _ := Request("POST", podsurl, jsonStr)
	var podslist []v1.Pod
	json.Unmarshal([]byte(podsmsg), &podslist)
	return podslist
}

//获取service信息
func ServiceMsg(namespaces string) []v1.Service {
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	svcurl := "http://10.10.7.175:8081/services/list"
	svcmsg, _ := Request("POST", svcurl, jsonStr)
	var servicelist []v1.Service
	json.Unmarshal(svcmsg, &servicelist)
	return servicelist
}

// 获取pod域名信息
func IngressMsg(namespaces string) []v1beta2.Ingress {
	jsonmap := make(map[string]string)
	jsonmap["namespaces"] = namespaces
	jsonStr, _ := json.Marshal(jsonmap)

	ingurl := "http://10.10.7.175:8081/ingress/list"
	ingmsg, _ := Request("POST", ingurl, jsonStr)
	var ingresslist []v1beta2.Ingress
	json.Unmarshal([]byte(ingmsg), &ingresslist)
	return ingresslist
}

func (this *PodsController) Get() {
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

	//获取pods信息
	podsurl := "http://10.10.7.175:8081/pods/list"
	podsmsg, _ := Request("POST", podsurl, jsonStr)
	var podslist []v1.Pod
	json.Unmarshal([]byte(podsmsg), &podslist)

	this.Data["podslist"] = podslist
	this.Layout = "layout.html"
	this.TplName = "pods.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "html_head.html"
	this.LayoutSections["BodyHead"] = "body_head.html"
	this.LayoutSections["Sidebar"] = "sidebar.html"
}
