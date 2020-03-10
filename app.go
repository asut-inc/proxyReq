package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Urls []struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

var urls Urls

func main() {
	readJson()

	f()

	fmt.Scanln()
}

func f() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	for i := 0; i < len(urls); i++ {
		doReq(i)
	}
}

func doReq(i int) {
	t_start := time.Now()

	err := os.Setenv("HTTP_PROXY", "http://"+urls[i].IP+":"+urls[i].Port)
	if err != nil {
		fmt.Println(i, err)
		return
	}

	proxyUrl, err := url.Parse("http://" + urls[i].IP + ":" + urls[i].Port)
	if err != nil {
		fmt.Println(i, err)
		return
	}

	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: 5 * time.Second}
	res, err := myClient.Get("http://2children.ru/news/index")
	if err != nil {
		fmt.Println(i, err)
		return
	}

	t_end := time.Now()
	fmt.Println(i, res.StatusCode, t_end.Sub(t_start))

	doReq(i + 1)
}

func readJson() {
	jsonFile, err := os.Open("proxies.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &urls)
}
