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

type Uamtk struct {
	Apptk string
	ResultMessage string
	ResultCode int
	Newapptk string
}

type Login struct {
	ResultMessage string
	ResultCode int
	Uamtk string
}

var JSESSIONID string
var theAnswer string
var imageCookie map[string]string
var codeCookie map[string]string
var loginCookie map[string]string
var uamtkCookie map[string]string
var uamtkTiket string

func main()  {
	//uamtk()
	//login()
	//answer := util.GetAnswer("test.png")
	//fmt.Println(answer)
	uamauthclient()
}

func loginTest()  {

}

func checkCode()  {
	answer := getAnswer()
	method := "GET"
	data := ""
	urlVal := "https://kyfw.12306.cn/passport/captcha/captcha-check?callback=jQuery1910028362015323499357_1577349946476&rand=sjrand&login_site=E&_=1577349946480&answer="+answer

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
	cookie1 := &http.Cookie{Name:"_passport_session",Value:imageCookie["_passport_session"]}
	req.AddCookie(cookie1)
	cookie2 := &http.Cookie{Name:"_passport_ct",Value:imageCookie["_passport_ct"]}
	req.AddCookie(cookie2)

	resp, err := client.Do(req)

	codeCookie := getCookie(resp.Cookies())
	fmt.Println(codeCookie)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(b))
}

func login() []byte {
	checkCode()

	answer := theAnswer
	method  := "POST"
	urlVal := "https://kyfw.12306.cn/passport/web/login"
	//data := "username=234324&password=23423423&appid=otn&answer="+answer

	data := url.Values{}
	data.Set("username", "13425144866")
	data.Set("password", "bingbing925455")
	data.Set("appid", "otn")
	data.Set("answer", answer)

	client := &http.Client{}
	req, createErr := http.NewRequest(method, urlVal,  strings.NewReader(data.Encode()))
	if createErr != nil {
		fmt.Printf("创建失败:%v\n",createErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")

	//可以添加多个cookie
	cookie1 := &http.Cookie{Name:"JSESSIONID",Value:JSESSIONID}
	req.AddCookie(cookie1)

	cookie2 := &http.Cookie{Name:"_passport_session",Value:imageCookie["_passport_session"]}
	req.AddCookie(cookie2)
	cookie3 := &http.Cookie{Name:"_passport_ct",Value:imageCookie["_passport_ct"]}
	req.AddCookie(cookie3)
	cookie4 := &http.Cookie{Name:"BIGipServerpassport",Value:codeCookie["BIGipServerpassport"]}
	req.AddCookie(cookie4)
	cookie5 := &http.Cookie{Name:"RAIL_DEVICEID",Value:"e4koiUxjwReuTWXCjalPjynEcNk0vk-kZQpbtKAWpIQHAxLxarlZLj27EYmWgQ0QfhEFU7nf43mkfvLEJbVGEpU4igEAmDp_vTIiojk83fSMJ4pWKlDwFy8kjpalg8dLvOIg7pJnAxSiUgibjRMqRcQiEIxsq_x8"}
	req.AddCookie(cookie5)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	loginCookie = getCookie(resp.Cookies())

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.Cookies())
	fmt.Println(string(b))

	return b
}

func uamtk() []byte {
	uamtk := login()
	var loginUa Login
	err := json.Unmarshal(uamtk,&loginUa)
	if err != nil{
		fmt.Printf("错误：%v\n",err)
	}
	uamtkTiket = loginUa.Uamtk

	fmt.Printf("票据：%v\n",uamtkTiket)

	method  := "POST"
	urlVal := "https://kyfw.12306.cn/passport/web/auth/uamtk"
	//data := "username=234324&password=23423423&appid=otn&answer="+answer

	data := url.Values{}
	data.Set("appid", "otn")
	data.Set("uamtk", uamtkTiket)

	client := &http.Client{}
	req, createErr := http.NewRequest(method, urlVal,  strings.NewReader(data.Encode()))
	if createErr != nil {
		fmt.Printf("创建失败:%v\n",createErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://exservice.12306.cn")
	req.Header.Set("Referer", "https://exservice.12306.cn/excater/index.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")

	//可以添加多个cookie
	cookie1 := &http.Cookie{Name:"JSESSIONID",Value:JSESSIONID}
	req.AddCookie(cookie1)

	cookie2 := &http.Cookie{Name:"_passport_session",Value:imageCookie["_passport_session"]}
	req.AddCookie(cookie2)
	cookie3 := &http.Cookie{Name:"_passport_ct",Value:loginCookie["_passport_ct"]}
	req.AddCookie(cookie3)
	cookie4 := &http.Cookie{Name:"BIGipServerpassport",Value:loginCookie["BIGipServerpassport"]}
	req.AddCookie(cookie4)
	cookie5 := &http.Cookie{Name:"RAIL_DEVICEID",Value:"e4koiUxjwReuTWXCjalPjynEcNk0vk-kZQpbtKAWpIQHAxLxarlZLj27EYmWgQ0QfhEFU7nf43mkfvLEJbVGEpU4igEAmDp_vTIiojk83fSMJ4pWKlDwFy8kjpalg8dLvOIg7pJnAxSiUgibjRMqRcQiEIxsq_x8"}
	req.AddCookie(cookie5)

	cookie6 := &http.Cookie{Name:"uamtk",Value:uamtkCookie["uamtk"]}
	req.AddCookie(cookie6)

	resp, err := client.Do(req)

	fmt.Println(resp.Cookies())

	uamtkCookie = getCookie(resp.Cookies())
	fmt.Printf("uamtk:%v\n",uamtkCookie)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("第三步：%v\n",string(b))

	return b
}

func uamauthclient()  {
	uamtkResult := uamtk()
	var uamtk Uamtk
	err := json.Unmarshal(uamtkResult,&uamtk)
	if err != nil{
		fmt.Printf("错误：%v\n",err)
	}

	fmt.Println("分界线")
	fmt.Println("uamtkCookie:",uamtk.Newapptk)

	fmt.Println(uamtk.Newapptk)

	method  := "POST"
	urlVal := "https://exservice.12306.cn/excater/uamauthclient?tk="+uamtk.Newapptk
	//data := "username=234324&password=23423423&appid=otn&answer="+answer

	fmt.Printf("tk是啥:%v\n",uamtk.Newapptk)

	data := url.Values{}
	data.Set("tk",uamtk.Newapptk)
	data.Set("uamtk", uamtkTiket)
	//data.Set("uamtk",uamtkTiket)

	client := &http.Client{}
	req, createErr := http.NewRequest(method, urlVal,  strings.NewReader(data.Encode()))
	if createErr != nil {
		fmt.Printf("创建失败:%v\n",createErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://exservice.12306.cn")
	req.Header.Set("Referer", "https://exservice.12306.cn/excater/index.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")

	//可以添加多个cookie

	cookie2 := &http.Cookie{Name:"_passport_session",Value:imageCookie["_passport_session"]}
	req.AddCookie(cookie2)
	cookie3 := &http.Cookie{Name:"_passport_ct",Value:loginCookie["_passport_ct"]}
	req.AddCookie(cookie3)
	cookie4 := &http.Cookie{Name:"BIGipServerpassport",Value:loginCookie["BIGipServerpassport"]}
	req.AddCookie(cookie4)
	cookie5 := &http.Cookie{Name:"RAIL_DEVICEID",Value:"e4koiUxjwReuTWXCjalPjynEcNk0vk-kZQpbtKAWpIQHAxLxarlZLj27EYmWgQ0QfhEFU7nf43mkfvLEJbVGEpU4igEAmDp_vTIiojk83fSMJ4pWKlDwFy8kjpalg8dLvOIg7pJnAxSiUgibjRMqRcQiEIxsq_x8"}
	req.AddCookie(cookie5)


	cookie6 := &http.Cookie{Name:"tk",Value:uamtk.Newapptk}
	req.AddCookie(cookie6)
	cookie7 := &http.Cookie{Name:"uamtk",Value:uamtkTiket}
	req.AddCookie(cookie7)

	fmt.Printf("uamtk：%v\n",uamtkCookie["uamtk"])

	resp, err := client.Do(req)

	fmt.Printf("最后一步：%v\n",resp.Cookies())

	uamtkCookie = getCookie(resp.Cookies())

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("结果")
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

	theAnswer = answer

	return answer
}

func getBase64Image(cookie string) []byte {
	method  := ""
	urlVal := "https://kyfw.12306.cn/passport/captcha/captcha-image64?login_site=E&module=login&rand=sjrand&1577093928867"
	data := ""

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
