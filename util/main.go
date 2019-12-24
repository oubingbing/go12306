package util

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

/**
 * 执行命令
 */
func execCommand(com string,arg ...string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(com, arg[0],arg[1])
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

/**
 * 获取12306验证码
 */
func GetAnswer(imgPath string) string {
	dir, _ := os.Getwd()
	err, out, _ := execCommand("python",dir+"/easy12306/main.py",dir+"/upload/image/"+imgPath)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	ret := strings.Split(out,";")
	keys := ret[0]
	position := ret[1]

	keyString := strings.Split(keys,",")
	positionString := strings.Split(position," ")

	mp := make(map[string]string)
	mp["00"] = "35,45"
	mp["01"] = "111,40"
	mp["02"] = "185,50"
	mp["03"] = "260,45"
	mp["10"] = "44,115"
	mp["11"] = "115,115"
	mp["12"] = "185,115"
	mp["13"] = "260,115"

	var code []string
	for i,v := range positionString {
		if i <= 7 {
			ps := strings.Split(v,",")
			if len(ps[0]) > 0 {
				for _,kv := range keyString {
					if ps[1] == kv {
						code = append(code, mp[ps[0]])
					}
				}
			}
		}
	}

	var answer string
	for i,v := range code{
		if i < len(code)-1 {
			answer += v+","
		}else{
			answer += v
		}
	}

	return answer
}