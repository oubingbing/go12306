package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Order struct{
	KyfwPrt *Kyfw
	RepeatSubmitToken string 			//下单需要的token，在下单页面可以获得
	KeyCheckIsChangeToken string 		//下单需要的token，在下单页面可以获得
	PassengerInfo *GetPassengerResult 	//购票乘客信息
	LeftTicket string 					//获取排队数以及提交订单时候需要用到
	Secret string 						//购买车票的秘钥
	//Station map[string][]string 		//所有的车站信息
	TicketForm *QueryTicketForm

	/**
	 * [DZl%2B02BYd56T7BBLRZNTwvjVPBfuwVgJN4zKm%2F7SXgcpCDnA8VpyqibNvC2mKSqetOtkA7JYhEA%2F%0AYun17t5iIbrxbJq5WfMuLZGpM10Uxsxnm9M4RkQQ2%2FQbH9lWMdeRwGimFeCLmSeDGIDU8qvbNZXy%0A6xtT%2BPv4
	 *	sxz9wxUVj5AXxmcVEYtaT%2BbKnTY4v%2FTtrK%2FRlSSB3NGy%2FIGHZR5Z33YZt6BCJWouQG3h%0AAw950UyCk8TFgYKYLMkXjM5QgDtRrvhfIKrKHVwrC6I%2B%2BdOdZbth1CA3%2FB%2FH 预订 6c000D718502 D7185 IZQ MDQ IZQ
	 *	MDQ 06:32 08:43 02:11 Y B703YNGISlA86w%2BW29GAi3cX7tzlhnDH 20200110 3 Q6 01 04 1 0       无    有    O0O0 OO 0 0         0]
	*/
	TargetTrainInfo []string			//购买车次的信息
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

type SubmitOrderRequestResult struct {
	Data string
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

type QueryTicketForm struct {
	TrainNo string
	TrainDate string
	FromStation []string //[gzn 广州南 IZQ guangzhounan gzn 5]
	ToStation []string	 //[gzn 广州南 IZQ guangzhounan gzn 5]
	PurposeCodes string
	PassengerName string //订票人
}

type OrderCommonResponse struct {
	ValidateMessagesShowId string
	Status bool
	Httpstatus int
	Messages []string
	ValidateMessages interface{}
	Data map[string]interface{}
}

/**
 * 检测用户状态
 */
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

/**
 * 发起下单请求
 */
func (order *Order) SubmitOrderRequest() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/leftTicket/submitOrderRequest"

	data := url.Values{}
	data.Set("secretStr",order.Secret)
	data.Set("train_date",order.TicketForm.TrainDate)
	data.Set("back_train_date","")
	data.Set("tour_flag","dc")
	data.Set("purpose_codes",order.TicketForm.PurposeCodes)
	data.Set("query_from_station_name",order.TicketForm.FromStation[1])
	data.Set("query_to_station_name",order.TicketForm.ToStation[1])

	httpErr := client.Post(urlVal,data, func(req *http.Request) {
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
		var submitOrderRequestResult SubmitOrderRequestResult
		b, err = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(b,&submitOrderRequestResult)
		if err != nil{
			fmt.Println(submitOrderRequestResult.Httpstatus)
			fmt.Printf("发起提交订单请求失败：%v\n",err.Error())
		}else{
			fmt.Printf("发起提交订单请求成功：%v\n",string(b))
		}

	})

	if httpErr != nil{
		err = httpErr
	}

	return err
}

/**
 * 查询车票信息
 */
func (order *Order) QueryTicket() error {
	var secretStrDecode string
	var err error
	isExistTrain := false

	urlVal := "https://kyfw.12306.cn/otn/leftTicket/queryZ"+order.TicketForm.String()
	var client HttpClient
	httpErr := client.Get(urlVal, func(req *http.Request) {
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
		err = json.Unmarshal(b,&result)
		if err != nil{
			fmt.Printf("查询车票json解析错误:%v\n",err)
		}

		var secretStr string
		for _,v := range result.Data.Result  {
			trainArr := strings.Split(v,"|")
			if trainArr[3] == order.TicketForm.TrainNo{
				secretStr = trainArr[0]
				order.LeftTicket = trainArr[12]
				order.TargetTrainInfo = trainArr
				isExistTrain = true
				//fmt.Printf("找到了:%v\n",trainArr)
				break
			}
		}

		secretStrDecode,err = url.QueryUnescape(secretStr)
		order.Secret = secretStrDecode
		if err != nil{
			fmt.Printf("车次秘钥解析错误：%v\n",err.Error())
		}

		fmt.Printf("车次秘钥：%v\n",order.LeftTicket)
	})

	if !isExistTrain {
		return errors.New("当前车次不存在")
	}

	if httpErr != nil{
		err = httpErr
	}

	return err
}

/**
 * 获取乘客信息
 */
func (order *Order) GetPassengerDTOs() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/getPassengerDTOs"

	data := url.Values{}
	data.Set("_json_att","")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken) //可在该页面中获取https://kyfw.12306.cn/otn/confirmPassenger/initDc

	httpErr := client.Post(urlVal,data, func(req *http.Request) {
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
		err = json.Unmarshal(b,&passengerInfo)
		if err != nil{
			fmt.Printf("解析乘客josn出错，请重新尝试，错误信息：%v\n",err.Error())
		}else{
			if passengerInfo.Data.ExMsg != "" {
				err = errors.New(passengerInfo.Data.ExMsg)
			}else{
				order.PassengerInfo = &passengerInfo
				fmt.Println("获取乘客信息成功")
			}
		}

	})

	if httpErr != nil{
		err = httpErr
	}

	return err
}

/**
 * 初始化订单页面
 */
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
		if keyCheckIsChangeIndex > 100 {
			keyCheckIsChangeNextIndex := strings.Index(str,"','leftDetails")
			keyCheckIsChangeToken := str[keyCheckIsChangeIndex+(len(keyCheckIsChangeName)):keyCheckIsChangeNextIndex]
			order.KeyCheckIsChangeToken = keyCheckIsChangeToken

			fmt.Printf("生成订单页面的token：%v\n", tokenStr[1])
			fmt.Printf("排序token：%v\n",keyCheckIsChangeToken)
		}
	})

	return err

}

/**
 * 检测订单信息
 */
func (order *Order) CheckOrderInfo() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo"

	ps := order.PassengerInfo.Data.Normal_passengers
	var passengeData map[string]string
	isExistPassenger := false
	for _,v := range ps {
		if v["passenger_name"] == order.TicketForm.PassengerName{
			passengeData = v
			isExistPassenger = true
			break
		}
	}

	if !isExistPassenger {
		isExistErr := "乘客【"+order.TicketForm.PassengerName+"】不存你的账号中，请先添加该乘客然后重试"
		return errors.New(isExistErr)
	}

	data := url.Values{}
	data.Set("cancel_flag","2")
	data.Set("bed_level_order_num","000000000000000000000000000000")
	data.Set("passengerTicketStr","O,0,1,"+passengeData["passenger_name"]+",1,"+passengeData["passenger_id_no"]+","+passengeData["mobile_no"]+",N,"+passengeData["allEncStr"])
	data.Set("oldPassengerStr",passengeData["passenger_name"]+",1,"+passengeData["passenger_id_no"]+",1_")
	data.Set("tour_flag","dc")
	data.Set("randCode","")
	data.Set("whatsSelect","1")
	data.Set("sessionId","")
	data.Set("sig","")
	data.Set("scene","nc_login")
	data.Set("_json_att","nc_login")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken)

	httpErr := client.Post(urlVal,data, func(req *http.Request) {
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

		var commonResponse OrderCommonResponse
		err = json.Unmarshal(b,&commonResponse)
		if err != nil{
			fmt.Printf("检测订单失败：%v\n",err.Error())
		}else{
			fmt.Printf("检测订单成功：%v\n",string(b))
		}

	})

	if httpErr != nil {
		err = httpErr
	}

	return err
}

/**
 * 获取排队数
 */
func (order *Order) GetQueueCount() error {
	var urlVal string
	var client HttpClient
	var err error

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/getQueueCount"

	t,_ := time.Parse("2006-01-02 15:04:05",order.TicketForm.TrainDate+" 00:00:00")
	trainDateString := t.Weekday().String()+" "+t.Month().String()+" "+fmt.Sprintf("%v",t.Day())+" "+fmt.Sprintf("%v",t.Year())+" 00:00:00 GMT+0800 (中国标准时间)"

	data := url.Values{}
	data.Set("train_date",trainDateString)
	data.Set("train_no",order.TargetTrainInfo[2])
	data.Set("stationTrainCode",order.TargetTrainInfo[3])
	data.Set("seatType","WZ")
	data.Set("fromStationTelecode",order.TicketForm.FromStation[2])
	data.Set("toStationTelecode",order.TicketForm.ToStation[2])
	data.Set("leftTicket",order.LeftTicket)
	data.Set("purpose_codes","00")
	data.Set("train_location",order.TargetTrainInfo[15])
	data.Set("_json_att","")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken)

	fmt.Printf("提交的数据：%v\n",data)

	httpErr := client.Post(urlVal,data, func(req *http.Request) {
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

	if httpErr != nil{
		err = httpErr
	}

	return err
}

/**
 * 将订单提交到12306订票队列
 */
func (order *Order) ConfirmSingleForQueue() error {
	var urlVal string
	var client HttpClient
	var err error

	fmt.Println()

	urlVal = "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueue"

	ps := order.PassengerInfo.Data.Normal_passengers

	var passengeData map[string]string
	for _,v := range ps {
		if v["passenger_name"] == order.TicketForm.PassengerName{
			passengeData = v
			break
		}
	}

	data := url.Values{}
	data.Set("passengerTicketStr","O,0,1,"+passengeData["passenger_name"]+",1,"+passengeData["passenger_id_no"]+","+passengeData["mobile_no"]+",N,"+passengeData["allEncStr"])
	data.Set("oldPassengerStr",passengeData["passenger_name"]+",1,"+passengeData["passenger_id_no"]+",1_")
	data.Set("randCode","")
	data.Set("purpose_codes","00")
	data.Set("key_check_isChange",order.KeyCheckIsChangeToken)
	data.Set("leftTicketStr",order.LeftTicket)
	data.Set("train_location",order.TargetTrainInfo[15])
	data.Set("choose_seats","1F")
	data.Set("seatDetailType","000")
	data.Set("whatsSelect","1")
	data.Set("roomType","00")
	data.Set("dwAll","N")
	data.Set("_json_att","")
	data.Set("REPEAT_SUBMIT_TOKEN",order.RepeatSubmitToken)

	httpErr := client.Post(urlVal,data, func(req *http.Request) {
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

		var commonResponse OrderCommonResponse
		err = json.Unmarshal(b,&commonResponse)
		if err != nil {
			fmt.Printf(">>> 提交订单错误：%v\n",string(len(b)))
		}else{
			if commonResponse.Status == false{
				err = errors.New(">>> 订单提交失败")
			}else if commonResponse.Data["submitStatus"] == true {
				fmt.Printf(">>> 下单成功:%v\n",string(b))
			}else{
				fmt.Printf(">>> 订单提交错误:%v\n",commonResponse.Data["errMsg"])
				err = errors.New(">>> 订单提交失败")
			}
		}

	})

	if httpErr != nil{
		err = httpErr
	}

	return err
}

/**
 * 获取站点信息
 */
func (order *Order) GetStation(fromStation string,toStation string) ([]string,[]string,error) {
	stationMap := make(map[string][]string)
	urlVal := "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9138"
	var client HttpClient
	err := client.Get(urlVal, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("User-Agent", "application/x-www-form-urlencoded")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("Host", "kyfw.12306.cn")
		req.Header.Set("Referer", "https://kyfw.12306.cn/otn/leftTicket/init?linktypeid=dc")
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
		stationStr := string(b)

		indexStr := "var station_names ='"
		index := strings.Index(stationStr,indexStr)
		stationStrInfo := stationStr[(index+len(indexStr)):len(stationStr)-2]
		stationSlice := strings.Split(stationStrInfo,"@")

		for _,v := range stationSlice {
			if len(v) > 0{
				stationInfo := strings.Split(v,"|")
				if stationInfo != nil{
					stationMap[stationInfo[1]] = stationInfo
				}
			}
		}

		//order.Station = stationMap
	})

	f := stationMap[fromStation]
	if f == nil{
		return nil,nil,errors.New("始发站不存在")
	}


	t := stationMap[toStation]
	if t == nil{
		return nil,nil,errors.New("终点站不存在")
	}

	return stationMap[fromStation],stationMap[toStation],err
}

/**
 * 组织车票查询字符串
 */
func (queryTicketForm *QueryTicketForm) String() string {
	return "?leftTicketDTO.train_date="+queryTicketForm.TrainDate+"&leftTicketDTO.from_station="+queryTicketForm.FromStation[2]+"&leftTicketDTO.to_station="+queryTicketForm.ToStation[2]+"&purpose_codes="+queryTicketForm.PurposeCodes
}