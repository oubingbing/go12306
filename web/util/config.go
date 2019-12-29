package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetMysqlConfig() (string,error) {
	configs := GetAppConfig()

	for key,_ := range configs {
		_,ok := configs[key]
		if !ok {
			return "",errors.New(key+"配置信息错误")
		}
	}

	var builder strings.Builder
	builder.WriteString(strings.TrimSpace(configs["mysql_username"]))
	builder.WriteString(":")
	builder.WriteString(strings.TrimSpace(configs["mysql_psw"]))
	builder.WriteString("@(")
	builder.WriteString(strings.TrimSpace(configs["mysql_host"]))
	builder.WriteString(":")
	builder.WriteString(strings.TrimSpace(configs["mysql_port"]))
	builder.WriteString(")/")
	builder.WriteString(strings.TrimSpace(configs["mysql_dbname"]))
	builder.WriteString("?charset=utf8")
	builder.WriteString("&parseTime=True&loc=Local")
	return builder.String(),nil
}

func GetAppConfig() map[string]string {
	dir, _ := os.Getwd()
	f,err := os.OpenFile(dir+"/app.conf",os.O_RDONLY,0777)
	if err != nil {
		Error(fmt.Sprintf("获取配置文件失败：%v\n",err.Error()))
	}

	reader := bufio.NewReader(f)
	configMp := make(map[string]string)
	for {
		line, err := reader.ReadString('\n') //以'\n'为结束符读入一行
		fmt.Printf("line:%v\n",line)
		if err != nil || io.EOF == err {
			break
		}

		if strings.Index(line,"=") == -1 {
			continue
		}

		configKey := line[0:strings.Index(line,"=")]
		configValue := line[strings.Index(line,"=")+1:]
		configMp[configKey] = configValue
	}

	return configMp
}
