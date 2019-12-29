package main

import (
	"easy/util"
	"fmt"
)

func main()  {
	var kyfw *util.Kyfw
	var err error

	kyfw,err = util.AuthKyf()
	if err != nil{
		fmt.Println(err.Error())
	}

	var order util.Order
	order.KyfwPrt = kyfw
	order.CheckUser()

	//需要确认站点是否存在
	fromStation,toStation,stationErr := order.GetStation("广州东","深圳")
	if stationErr != nil{
		fmt.Println(stationErr)
	}
	fmt.Printf("获取站点结果：%v,%v\n",fromStation,toStation)

	var queryTicketForm util.QueryTicketForm
	queryTicketForm.FromStation 	= fromStation
	queryTicketForm.ToStation 		= toStation
	queryTicketForm.TrainDate 		= "2020-01-05"
	queryTicketForm.PurposeCodes 	= "ADULT"
	queryTicketForm.TrainNo 		= "C7027"
	queryTicketForm.PassengerName 	= "区志彬"

	order.TicketForm = &queryTicketForm

	//获取车票信息，尝试五次
	queryTickeTry := 3
	for i:=0;i<queryTickeTry;i++  {
		err = order.QueryTicket()
		if err != nil{
			fmt.Println(err.Error())
		}else{
			fmt.Println(order.Secret)
			break
		}
	}

	if len(order.Secret) <= 0{
		fmt.Println("获取车票失败")
		return
	}

	order.SubmitOrderRequest()

	order.InitDc()

	//获取车票信息，尝试五次
	etPassengerDTOTry := 3
	for i:=0;i<etPassengerDTOTry;i++  {
		err = order.GetPassengerDTOs()
		if err != nil{
			fmt.Println(err.Error())
		}else{
			order.GetPassengerDTOs()
			break
		}
	}


	if len(order.Secret) <= 0{
		fmt.Println("获取乘客信息失败")
		return
	}


	err = order.CheckOrderInfo()
	if err != nil{
		order.CheckOrderInfo()
	}

	order.GetQueueCount()

	order.ConfirmSingleForQueue()
}
