package main

import (
	"easy/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ImageResult struct {
	Image string
	ResultMessage string
	ResultCode string
}

var JSESSIONID string
var imageCookie map[string]string

func main()  {
	answer := getAnswer()
	fmt.Println(answer)
	fmt.Println(imageCookie)
	method  := "POST"
	urlVal := "https://kyfw.12306.cn/passport/web/login"
	//data := "username=234324&password=23423423&appid=otn&answer="+answer

	data := url.Values{}
	data.Set("username", "rnben")
	data.Set("password", "rnben")
	data.Set("appid", "otn")
	data.Set("answer", answer)

	client := &http.Client{}
	req, createErr := http.NewRequest(method, urlVal,  strings.NewReader(data.Encode()))
	fmt.Printf("创建失败:%v\n",createErr)

	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/login/init")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://kyfw.12306.cn/otn/resources/login.html")
	req.Header.Set("Host", "kyfw.12306.cn")
	req.Header.Set("User-Agent", "application/x-www-form-urlencoded")

	//可以添加多个cookie
	cookie1 := &http.Cookie{Name:"JSESSIONID",Value:JSESSIONID}
	req.AddCookie(cookie1)

	resp, err := client.Do(req)

	resp.Cookies()

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func getAnswer() string {
	jeCookie := getJSESSIONIDCookie()
	var imageResult ImageResult
	image := getBase64Image(jeCookie["JSESSIONID"])
	JSESSIONID = jeCookie["JSESSIONID"]
	json.Unmarshal(image,&imageResult)
	imageName,_ := saveImage(imageResult.Image)
	answer := util.GetAnswer(imageName)

	return answer
}

func getBase64Image(cookie string) []byte {
	method  := ""
	urlVal := "https://kyfw.12306.cn/passport/captcha/captcha-image64?login_site=E&module=login&rand=sjrand&1577093928867&"
	data := ""

	_,err := http.Get(urlVal)

	client := &http.Client{}
	var req *http.Request

	if data == "" {
		urlArr := strings.Split(urlVal,"?")
		if len(urlArr)  == 2 {
			urlVal = urlArr[0] + "?" + getParseParam(urlArr[1])
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	}else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}

	//可以添加多个cookie
	cookie1 := &http.Cookie{Name:"JSESSIONID",Value:cookie}
	req.AddCookie(cookie1)

	resp, err := client.Do(req)

	imageCookie = getCookie(resp.Cookies())

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	return b
}

func getJSESSIONIDCookie() map[string]string  {
	response,err := http.Get("https://kyfw.12306.cn/otn")
	if err != nil{
		fmt.Printf("错误 %v\n",err.Error())
	}

	ck := response.Cookies()
	defer response.Body.Close()

	return getCookie(ck)
}

func getCookie(ck []*http.Cookie) map[string]string {
	cookies := make(map[string]string)
	for _,v := range ck {
		vs := strings.Split(v.String(),";")
		for _,vsub := range vs {
			sub := strings.Split(string(vsub),"=")
			cookies[sub[0]] = sub[1]
		}
	}

	return cookies
}

func saveImage(base64ImgString string) (string,error)  {
	dir, _ := os.Getwd()
	imageName := time.Now().Format("2006-10-12-12-23-34")+string(rand.Intn(1000))+".png"
	imagePath := dir+"/upload/image/"+imageName
	ddd, _ := base64.StdEncoding.DecodeString(base64ImgString) //成图片文件并把文件写入到buffer
	err := ioutil.WriteFile(imagePath, ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
	return  imageName,err
}

func getParseParam(param string) string  {
	return url.PathEscape(param)
}

func httpHandle(method, urlVal,data string)  {
	client := &http.Client{}
	var req *http.Request

	if data == "" {
		urlArr := strings.Split(urlVal,"?")
		if len(urlArr)  == 2 {
			urlVal = urlArr[0] + "?" + getParseParam(urlArr[1])
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	}else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}

	//可以添加多个cookie
	cookie1 := &http.Cookie{Name:"BIGipServerindex=1071186186.43286.0000",Value:"1071186186.43286.0000"}
	//cookie1 := &http.Cookie{Name:"JSESSIONID",Value:"123123123"}
	req.AddCookie(cookie1)

	resp, err := client.Do(req)

	resp.Cookies()

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}
