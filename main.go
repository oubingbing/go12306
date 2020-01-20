package main

import (
	"easy/util"
	"fmt"
	"sync"
	"time"
)

func main()  {
	passenger := "区志彬"					//1、购票乘客
	username := ""							//2、登录账号密码
	password := ""							//3、登录秘密
	date := "2020-01-29"					//4、乘车日期
	fromStation := "怀化南"					//5、始发站
	toStation := "广州南"					//6、终点站
	getTicketTime := "2019-12-31 17:29:50"	//7、开始抢票时间

	wait(getTicketTime)

	kyfw := login(username,password)

	trainSlice := []string{"G6166","G6175","G6173","G6141","G16166"}	//8、购买的车次，可以多选，支持并发
	dotime := len(trainSlice)

	var wg sync.WaitGroup
	chn := make(chan int,dotime)
	wg.Add(dotime)

	for i:=0; i<dotime;i++  {
		trainNo := trainSlice[i]
		var queryTicketForm util.QueryTicketForm
		queryTicketForm.TrainDate 		= date
		queryTicketForm.PurposeCodes 	= "ADULT"
		queryTicketForm.TrainNo 		= trainNo
		queryTicketForm.PassengerName 	= passenger
		go order(chn,&wg,kyfw,&queryTicketForm,fromStation,toStation)
	}

	go done(chn,&wg)

	for c := range chn {
		fmt.Println(c)
	}

	fmt.Println("结束")
}

func done(chn chan int,group *sync.WaitGroup)  {
	group.Wait()
	close(chn)
}

/**
 * 定时抢票
 */
func wait(getTicketTime string)  {

	fmt.Println(">>> 等待抢票...")

	if getTicketTime != ""{
		format := "2006-01-02 15:04:05"
		getTicketTimeFormat, _ := time.Parse(format, getTicketTime)

		//判断是否可以抢票
		var now time.Time
		standby := true
		for standby  {
			now,_ = time.Parse(format,time.Now().Format(format))
			if now.After(getTicketTimeFormat) {
				standby = false
				break
			}else{
				time.Sleep(time.Millisecond * 500)
			}
		}
	}
}

/**
 * 登录12306
 */
func login(username,password string) *util.Kyfw {
	var kyfw *util.Kyfw
	var err error

	kyfw,err = util.AuthKyf(username,password)
	if err != nil{
		fmt.Println(err.Error())
	}

	return kyfw
}

/**
 * 进行下单操作
 */
func order(chn chan int,wg *sync.WaitGroup,kyfw *util.Kyfw,queryTicketForm *util.QueryTicketForm,from string,to string)  {
	var order util.Order
	order.KyfwPrt = kyfw

	order.CheckUser()

	//需要确认站点是否存在
	fromStation,toStation,stationErr := order.GetStation(from,to)
	if stationErr != nil{
		fmt.Println(stationErr)
	}
	fmt.Printf(">>> 获取站点结果：%v,%v\n",fromStation,toStation)

	//var queryTicketForm util.QueryTicketForm
	queryTicketForm.FromStation 	= fromStation
	queryTicketForm.ToStation 		= toStation

	order.TicketForm = queryTicketForm

	//获取车票信息，尝试五次
	getTicketTry := 5
	var getTicketErr error
	var getTicketResult bool = false
	for i:=0;i<getTicketTry;i++  {
		getTicketErr = order.QueryTicket()
		if getTicketErr == nil{
			getTicketResult = true
			break
		}
	}
	if !getTicketResult {
		fmt.Println(">>> 下单失败,获取车票信息失败")
		return
	}

	submitOrderRequestTry := 5
	var submitOrderRequestErr error
	var submitOrderRequestResult bool = false
	for i:=0;i<submitOrderRequestTry;i++  {
		submitOrderRequestErr = order.SubmitOrderRequest()
		if submitOrderRequestErr == nil{
			submitOrderRequestResult = true
			break
		}
	}

	if !submitOrderRequestResult{
		fmt.Println(">>> 下单失败,发起订单请求失败")
		return
	}

	order.InitDc()

	//获取车票信息，尝试五次
	etPassengerDTOTry := 5
	var getPassengerDTOsErr error
	var getPassengerDTOsResult bool = false
	for i:=0;i<etPassengerDTOTry;i++  {
		getPassengerDTOsErr = order.GetPassengerDTOs()
		if getPassengerDTOsErr == nil{
			getPassengerDTOsResult = true
			break
		}
	}

	if !getPassengerDTOsResult {
		fmt.Println(">>> 下单失败,获取乘客信息失败")
		return
	}

	checkOrderTry := 5
	var checkOrderInfoErr error
	checkOrderInfoResult := false
	for i:=0;i<checkOrderTry;i++  {
		checkOrderInfoErr = order.CheckOrderInfo()
		if checkOrderInfoErr == nil{
			checkOrderInfoResult = true
			break
		}
	}

	if !checkOrderInfoResult {
		fmt.Println(">>> 下单失败,检测订单失败")
		return
	}

	order.GetQueueCount()

	submitOrderTry := 5
	var submitOrderErr error
	submitOrderfoResult := false
	for i:=0;i<submitOrderTry;i++  {
		submitOrderErr = order.ConfirmSingleForQueue()
		if submitOrderErr == nil{
			submitOrderfoResult = true
			break
		}
	}

	if !submitOrderfoResult {
		fmt.Println(">>> 下单失败")
	}else{
		fmt.Println(">>> 下单成功")
	}

	fmt.Printf("抢票结束：%v\n",queryTicketForm.TrainNo)

	defer func() {
		chn <- 1
		wg.Done()
	}()
}