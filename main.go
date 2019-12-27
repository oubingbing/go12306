package main

import (
	"easy/util"
	"fmt"
)

func main()  {
	var kyfw util.Kyfw
	var err error
	kyfw.Cookies = map[string]string{}
	err = kyfw.InitLogin()
	if err != nil {
		fmt.Println(err.Error())
	}

	var imageData string
	imageData,err = kyfw.GetBase64Image()
	if err != nil {
		fmt.Println(err.Error())
	}

	var image string
	image,err = util.SaveImage(imageData)
	kyfw.GetAnswer(image)
	fmt.Println(kyfw.Answer)

	err = kyfw.CheckCode()
	if err != nil {
		fmt.Println(err)
	}

	err = kyfw.Login()
	if err != nil {
		fmt.Println(err)
	}

	err = kyfw.Uamtk()
	if err != nil {
		fmt.Println(err)
	}

	err = kyfw.Uamauthclient()
	if err != nil {
		fmt.Println(err)
	}
}
