package main

import (
	"easy/util"
	"fmt"
)

func main()  {
	/*var kyfw util.Kyfw
	kyfw.GetAnswer("test.jpg")
	fmt.Println(kyfw.Answer)*/

	var kyfw *util.Kyfw
	var err error

	kyfw,err = util.AuthKyf()
	if err != nil{
		fmt.Println(err.Error())
	}

	var order util.Order
	order.KyfwPrt = kyfw
	order.CheckUser()

	err = order.QueryTicket("G6248")
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println(order.Secret)

	order.SubmitOrderRequest()

	order.InitDc()

	order.GetPassengerDTOs()

	order.CheckOrderInfo()

	order.GetQueueCount()

	order.ConfirmSingleForQueue()
}
