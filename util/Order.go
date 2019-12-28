package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Order struct{
	KyfwPrt *Kyfw
	RepeatSubmitToken string
	KeyCheckIsChangeToken string
	PassengerInfo *GetPassengerResult
	LeftTicket string
	Secret string
}

type QueryTicketResult struct {
	Data QueryTicketSubResult
	Httpstatus int
	Messages string
	status bool
}

type QueryTicketSubResult struct {
	Flag string
	Map map[string]string
	Result []string
}

type GetPassengerResult struct {
	Data PassengerData
	Httpstatus int
	Messages []string
	Status bool
	ValidateMessages interface{}
	ValidateMessagesShowId string
}

type PassengerData struct {
	Dj_passengers []map[string]string
	ExMsg string
	IsExist bool
	Normal_passengers []map[string]string
	Notify_for_gat string
	Other_isOpenClick []string
	Two_isOpenClick []string
}

func (order *Order) CheckUser()  {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/login/checkUser"

	data := url.Values{}
	data.Set("_json_att","")

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		fmt.Printf("check user检测用户：%v\n",string(b))
	})
}

func (order *Order) SubmitOrderRequest()  {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/leftTicket/submitOrderRequest"

	data := url.Values{}
	data.Set("secretStr",order.Secret)
	data.Set("train_date","2019-12-28")
	data.Set("back_train_date","2019-12-28")
	data.Set("tour_flag","dc")
	data.Set("purpose_codes","ADULT")
	data.Set("query_from_station_name","广州")
	data.Set("query_to_station_name","深圳")

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		fmt.Printf("发起提交订单请求：%v\n",string(b))
	})
}

func (order *Order) QueryTicket(trainNo string) error {
	var secretStrDecode string

	urlVal := "https://kyfw.12306.cn/otn/leftTicket/queryZ?leftTicketDTO.train_date=2019-12-29&leftTicketDTO.from_station=SZQ&leftTicketDTO.to_station=GZQ&purpose_codes=ADULT"
	var client HttpClient
	err := client.Get(urlVal, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("User-Agent", "application/x-www-form-urlencoded")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("Host", "kyfw.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc&fs=%E6%B7%B1%E5%9C%B3,SZQ&ts=%E5%B9%BF%E5%B7%9E,GZQ&date=2019-12-29&flag=N,Y,Y")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		b, _ := ioutil.ReadAll(resp.Body)

		var result QueryTicketResult
		jsonerr := json.Unmarshal(b,&result)
		if jsonerr != nil{
			fmt.Printf("查询车票json解析错误:%v\n",jsonerr)
		}

		var  secretStr string
		for _,v := range result.Data.Result  {
			trainArr := strings.Split(v,"|")
			if trainArr[3] == trainNo{
				secretStr = trainArr[0]
				order.LeftTicket = trainArr[12]
				break
			}
		}


		var decodeErr error
		secretStrDecode,decodeErr = url.QueryUnescape(secretStr)
		order.Secret = secretStrDecode
		if decodeErr != nil{
			fmt.Printf("车次秘钥解析错误：%v\n",decodeErr.Error())
		}

		fmt.Printf("车次秘钥：%v\n",order.LeftTicket)
		//fmt.Printf("查询车票：%v\n",string(b))
	})

	return err
}

func (order *Order) GetPassengerDTOs() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/getPassengerDTOs"

	data := url.Values{}
	data.Set("_json_att","")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken)

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)

		var passengerInfo GetPassengerResult
		jsonerr := json.Unmarshal(b,&passengerInfo)
		if jsonerr != nil{
			fmt.Printf("解析乘客josn出错：%v\n",jsonerr.Error())
		}

		order.PassengerInfo = &passengerInfo

		//ps := order.PassengerInfo.Data.Normal_passengers

		//fmt.Printf("获取乘客：%v\n",passengerInfo.Data)
		fmt.Printf("获取乘客信息：%v\n",string(b))
	})

	return err
}

func (order *Order) InitDc() error {
	urlVal := "https://kyfw.12306.cn/otn/confirmPassenger/initDc"
	var client HttpClient
	err := client.Get(urlVal, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("User-Agent", "application/x-www-form-urlencoded")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("Host", "kyfw.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		b, _ := ioutil.ReadAll(resp.Body)
		str := string(b)
		tokenIndex := strings.Index(str,"globalRepeatSubmitToken =")
		tokenNextIndex := strings.Index(str," var global_lang")
		tokenStr := strings.Split(str[tokenIndex:tokenNextIndex],"'")
		order.RepeatSubmitToken = tokenStr[1]

		//key_check_isChange
		keyCheckIsChangeName := "key_check_isChange':'"
		keyCheckIsChangeIndex := strings.Index(str,keyCheckIsChangeName)
		keyCheckIsChangeNextIndex := strings.Index(str,"','leftDetails")
		keyCheckIsChangeToken := str[keyCheckIsChangeIndex+(len(keyCheckIsChangeName)):keyCheckIsChangeNextIndex]
		order.KeyCheckIsChangeToken = keyCheckIsChangeToken

		fmt.Printf("生成订单页面的token：%v\n", tokenStr[1])
		fmt.Printf("排序token：%v\n",keyCheckIsChangeToken)
	})

	return err

}

func (order *Order) CheckOrderInfo() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo"

	ps := order.PassengerInfo.Data.Normal_passengers
	var passengeStr string
	for _,v := range ps {
		if v["passenger_name"] == "徐钊峰"{
			passengeStr = v["allEncStr"]
			break
		}
	}

	data := url.Values{}
	data.Set("cancel_flag","2")
	data.Set("bed_level_order_num","000000000000000000000000000000")
	data.Set("passengerTicketStr",""+passengeStr)
	data.Set("oldPassengerStr","")
	data.Set("tour_flag","dc")
	data.Set("randCode","")
	data.Set("whatsSelect","1")
	data.Set("sessionId","")
	data.Set("sig","")
	data.Set("scene","nc_login")
	data.Set("_json_att","nc_login")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken)

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)

		fmt.Printf("检测订单信息：%v\n",string(b))
	})

	return err
}

func (order *Order) GetQueueCount() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/getQueueCount"

	data := url.Values{}
	data.Set("train_date","Sun Dec 29 2019 00:00:00 GMT+0800 (中国标准时间)")
	data.Set("train_no","6i000G62481A")
	data.Set("stationTrainCode","G6248")
	data.Set("seatType","0")
	data.Set("fromStationTelecode","IOQ")
	data.Set("toStationTelecode","IZQ")
	data.Set("leftTicket",order.LeftTicket)
	data.Set("purpose_codes","00")
	data.Set("train_location","Q7")
	data.Set("_json_att","")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken)

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)

		fmt.Printf("查询排队数：%v\n",string(b))
	})

	return err
}

func (order *Order) ConfirmSingleForQueue() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueue"

	ps := order.PassengerInfo.Data.Normal_passengers
	//passengeStr := ps[0]["allEncStr"]

	var passengeStr string
	for _,v := range ps {
		if v["passenger_name"] == "徐钊峰"{
			passengeStr = v["allEncStr"]
			break
		}
	}

	data := url.Values{}
	data.Set("passengerTicketStr",","+passengeStr)
	data.Set("oldPassengerStr","")
	data.Set("randCode","")
	data.Set("purpose_codes","00")
	data.Set("key_check_isChange",order.KeyCheckIsChangeToken)
	data.Set("leftTicketStr",order.LeftTicket)
	data.Set("train_location","Q7")
	data.Set("choose_seats","1F")
	data.Set("seatDetailType","000")
	data.Set("whatsSelect","1")
	data.Set("roomType","00")
	data.Set("dwAll","N")
	data.Set("_json_att","")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken)

	err = client.Post(urlVal,data, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "https://exservice.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
		for k,v := range order.KyfwPrt.Cookies  {
			req.AddCookie(&http.Cookie{Name:k,Value:v})
		}
	}, func(resp *http.Response) {
		codeCookie := GetCookies(resp.Cookies())
		for k,v := range codeCookie  {
			order.KyfwPrt.Cookies[k] = v
		}

		var b []byte
		b, err = ioutil.ReadAll(resp.Body)

		fmt.Printf("提交订单到队列：%v\n",string(b))
	})

	return err
}