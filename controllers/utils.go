package controllers

import (
	"DockerPlatform/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func RemoteCommand(cmd string) (status bool, msg string) {
	user := "root"
	pass := "5sjws!JS51l"
	host := "10.10.35.1"
	port := 6123

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(cmd); err != nil {
		return false, stderr.String()
	}
	return true, stdout.String()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MkdirIfNoExist(dir string) {
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0777)
	}
}
func CheckAdmin(username string) bool {
	o := orm.NewOrm()
	user := models.User{Name: username}
	o.Read(&user, "Name")

	return user.Isadmin
}

func CheckLogin(username string, password string) (err error) {
	// 从数据库认证用户是否存在，如果不存在，从bitbucket认证
	o := orm.NewOrm()
	user := models.User{Name: username}
	err = o.Read(&user, "Name")

	if err == orm.ErrNoRows {
		err = AuthBitbucket(username, password)
		if err == nil {
			// 认证成功，写入数据库
			o := orm.NewOrm()
			var user models.User
			user.Name = username
			user.Password = password
			o.Insert(&user)

			return nil
		} else {
			return err
		}
	}
	if password != user.Password {
		return errors.New("密码错误")
	}
	return nil
}

func GetPasswd(username string) (password string) {
	o := orm.NewOrm()
	user := models.User{Name: username}
	o.Read(&user, "Name")
	password = user.Password
	return
}

//检查harbor 项目示是否存在，默认以用户名为harbor项目名
func CheckHarborProExist(username string) bool {
	harborhost := beego.AppConfig.String("harborhost")
	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")

	url := harborhost + "/api/projects"

	result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)

	var harborprojectlist []HarborProjects
	json.Unmarshal(result, &harborprojectlist)

	for _, v := range harborprojectlist {
		if v.Name == username {
			return true
		}
	}
	return false
}

// 创建harbor 项目
func CreateHarborPro(proname string) {
	harborhost := beego.AppConfig.String("harborhost")
	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")

	url := harborhost + "/api/projects"

	d := make(map[string]interface{})

	d["project_name"] = proname
	d["public"] = 1

	project, _ := json.Marshal(d)
	RequestForAuth("POST", url, harboruser, harborpass, project)
}

// 获取harbor 项目id
func GetHarborProId(proname string) (id float64) {
	harborhost := beego.AppConfig.String("harborhost")
	harboruser := beego.AppConfig.String("harboruser")
	harborpass := beego.AppConfig.String("harborpass")

	url := harborhost + "/api/projects"

	result, _ := RequestForAuth("GET", url, harboruser, harborpass, nil)

	var harborprojectlist []HarborProjects
	json.Unmarshal(result, &harborprojectlist)

	for _, v := range harborprojectlist {
		if v.Name == proname {
			return v.Project_id
		}
	}
	return -1
}

// 认证bitbucket
func AuthBitbucket(username string, password string) (err error) {

	bitbucketapiurl := beego.AppConfig.String("bitbucketapiurl")
	url := bitbucketapiurl + "/projects"
	result, _ := RequestForAuth("GET", url, username, password, nil)

	var bitbucketlist interface{}
	json.Unmarshal(result, &bitbucketlist)

	m := bitbucketlist.(map[string]interface{})

	if _, ok := m["errors"]; ok {
		return errors.New("用户名或密码错误")
	}
	return nil
}

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func CheckType(jsonStr string) {
	var f interface{}
	var j = []byte(jsonStr)
	json.Unmarshal(j, &f)
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}

func Request(method, url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte(""), err
	}

	request.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return result, err
}

func RequestForAuth(method, url, user, passwd string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, _ := http.NewRequest(method, url, body)

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(user, passwd)

	var resp *http.Response
	resp, _ = http.DefaultClient.Do(request)

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return result, err
}
