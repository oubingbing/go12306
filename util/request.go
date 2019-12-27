package util

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HttpClient struct {

}

type beforeRequestHandle func(req *http.Request)

type afterRequestHandle func(resp *http.Response)

/**
 * get请求
 */
func (h HttpClient) Get(url string ,beforeHandle beforeRequestHandle,afterHandle afterRequestHandle) error {
	client := &http.Client{}
	var req *http.Request

	urlArr := strings.Split(url,"?")
	if len(urlArr)  == 2 {
		url = urlArr[0] + "?" + getParseParam(urlArr[1])
	}
	req, _ = http.NewRequest("GET", url, nil)

	if beforeHandle != nil{
		beforeHandle(req)
	}

	resp, err := client.Do(req)

	if afterHandle != nil {
		afterHandle(resp)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return  err
}

/**
 * Post请求
 */
func (h HttpClient) Post(urlVal string,data url.Values,beforeHandle beforeRequestHandle,afterHandle afterRequestHandle) error {
	method  := "POST"
	client := &http.Client{}
	req, createErr := http.NewRequest(method, urlVal,  strings.NewReader(data.Encode()))
	if createErr != nil {
		fmt.Printf("创建失败:%v\n",createErr)
	}

	if beforeHandle != nil{
		beforeHandle(req)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(req)

	if afterHandle != nil{
		afterHandle(resp)
	}

	defer resp.Body.Close()

	return err
}

func getParseParam(param string) string  {
	return url.PathEscape(param)
}
