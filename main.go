package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

var ostype = runtime.GOOS

func getCurrentPath() string {

	path, _ := os.Getwd()
	if ostype == "windows" {
		path = path + "\\"
	} else if ostype == "linux" {
		path = path + "/"
	}
	return path
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("error: %s, exit!\n", err.Error())
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)

}

var conf = struct {
	HTTPProxy  string `json:"http_proxy"`
	VpngateAPI string `json:"vpngate_api"`
}{
	HTTPProxy:  "http://127.0.0.1:1080",
	VpngateAPI: "http://www.vpngate.net/api/iphone/",
}

func main() {
	fmt.Scan()
	if isExist("./conf.json") {
		dat, err := ioutil.ReadFile("./conf.json")
		checkErr(err)
		err = json.Unmarshal(dat, &conf)
		checkErr(err)
	}

	urli := url.URL{}
	urlproxy, err := urli.Parse(conf.HTTPProxy)
	checkErr(err)
	client := &http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}}
	resp, err := client.Get(conf.VpngateAPI)
	checkErr(err)
	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)

	f, err := os.Create(fmt.Sprintf("%svpngate_%s.csv", getCurrentPath(), time.Now().Format("20060102150405"))) //创建文件
	defer f.Close()
	checkErr(err)
	var count = 0
	for {

		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		ret := strings.Split(strings.Trim(line, "\n"), ",")
		if count == 1 {
			f.WriteString(line)
		} else if count > 1 && len(ret) == 15 {
			decodeBytes, err := base64.StdEncoding.DecodeString(ret[14])
			if err != nil {
				fmt.Printf("base64 decode error :%s \n", err.Error())
			}
			sp := strings.Split(string(decodeBytes), "\n")
			for _, v := range sp {
				if strings.HasPrefix(v, "remote") {
					ret[14] = v
					f.WriteString(strings.Join(ret, ","))
					break
				}
			}
			// fmt.Println(ret[14])
		}
		count++
	}
	fmt.Println("complete")
	time.Sleep(time.Second)

}
