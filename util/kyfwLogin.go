package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type KyfwCookies map[string]string

type Kyfw struct {
	Username string
	Password string
	Answer string
	Cookies KyfwCookies
	Tk string
}

type ImageResponse struct {
	Image string
	ResultMessage string
	ResultCode string
}

type UamtkResult struct {
	Apptk string
	ResultMessage string
	ResultCode int
	Newapptk string
}

const deviceString = "X05VHDCVI2ThoQ2S1147iuZqsMDKNo1QusC8orrnprztgmmteMoFdXNgyRSCuGJ4m0TsYn2Tpv4vXiKcDWJ2GC1gLs4zCvP_13eiaDLzRI-CBnYHGb9dIfVYFzQsGDiLoamEqOPOc29DOV1BHTBokDKuBFKqAlcA"

/**
 * 初始化登录信息
 */
func (kyfw *Kyfw) InitLogin() error {
	url := "https://kyfw.12306.cn/otn/login/init"
	var client HttpClient
	err := client.Get(url, nil, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			kyfw.Cookies[k] = v
		}
	})

	return  err
}

/**
 * Uamauthclient
 */
func (kyfw *Kyfw) Uamauthclient() error  {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/uamauthclient"

	data := url.Values{}
	data.Set("tk",kyfw.Tk)

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://exservice.12306.cn/excater/index.html")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range kyfw.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			kyfw.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		fmt.Printf("Uamauthclient：%v\n",string(b))
	})

	return err
}

/**
 * Uamtk
 */
func (kyfw *Kyfw) Uamtk() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/passport/web/auth/uamtk"

	data := url.Values{}
	data.Set("appid", "otn")

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/passport?redirect=/otn/login/userLogin")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range kyfw.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			kyfw.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		fmt.Printf("Uamtk：%v\n",string(b))

		var uamtk UamtkResult
		err := json.Unmarshal(b,&uamtk)
		if err != nil{
			fmt.Printf("错误：%v\n",err)
		}
		kyfw.Tk = uamtk.Newapptk
	})

	return err
}

/**
 * 登录
 */
func (kyfw *Kyfw) Login() error {
	var urlVal string
	var client HttpClient
	var err error

	kyfw.Username = "13425144866"
	kyfw.Password = ""
	urlVal = "https://kyfw.12306.cn/passport/web/login"

	data := url.Values{}
	data.Set("username", kyfw.Username)
	data.Set("password", kyfw.Password)
	data.Set("appid", "otn")
	data.Set("answer", kyfw.Answer)

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("User-Agent", "application/x-www-form-urlencoded")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")

		cookie5 := &http.Cookie{Name:"RAIL_DEVICEID",Value:deviceString}
		req.AddCookie(cookie5)

		for k,v := range kyfw.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}

	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			kyfw.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		fmt.Printf("登录：%v\n",string(b))
	})

	return err
}

/**
 * 校验验证码
 */
func (kyfw *Kyfw) CheckCode() error {
	url := "https://kyfw.12306.cn/passport/captcha/captcha-check?callback=jQuery1910028362015323499357_1577349946476&rand=sjrand&login_site=E&_=1577349946480&answer="+kyfw.Answer
	var client HttpClient
	err := client.Get(url, func(req *http.Request) {
		for k,v := range kyfw.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			kyfw.Cookies[k] = v
		}

		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("验证码校验：%v\n",string(b))
	})

	return  err
}

/**
 * 获取验证码
 */
func (kyfw *Kyfw) GetAnswer(image string) {
	kyfw.Answer = GetAnswer(image)
}

/**
 * 获取base64图片
 */
func (kyfw *Kyfw) GetBase64Image() (string,error) {
	url := "https://kyfw.12306.cn/passport/captcha/captcha-image64?login_site=E&module=login&rand=sjrand&1577093928867"
	var client HttpClient
	var err error
	var base65Image string

	err = client.Get(url, func(req *http.Request) {
		for k,v := range kyfw.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			kyfw.Cookies[k] = v
		}

		var imageResponse ImageResponse
		var imageData []byte
		imageData, err = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(imageData,&imageResponse)
		base65Image = imageResponse.Image
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	return  base65Image,err
}

/**
 * 保存图片
 */
func SaveImage(base64ImgString string) (string,error)  {
	dir, _ := os.Getwd()
	imageName := time.Now().Format("2006-10-12-12-23-34")+string(rand.Intn(1000))+".png"
	imagePath := dir+"/upload/image/"+imageName
	ddd, _ := base64.StdEncoding.DecodeString(base64ImgString) //成图片文件并把文件写入到buffer
	err := ioutil.WriteFile(imagePath, ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
	return  imageName,err
}

/**
 * 将cookie转化成map
 */
func GetCookies(ck []*http.Cookie) map[string]string {
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