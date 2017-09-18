package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DoPost(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader((data))
	request, err := http.NewRequest("POST", url, body)
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

func DoGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
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
	body := bytes.NewReader((data))
	request, _ := http.NewRequest(method, url, body)

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("tanzhixu", "1QAZ2wsx")

	var resp *http.Response
	resp, _ = http.DefaultClient.Do(request)

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return result, err
}
