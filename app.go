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

func main() {
	urls := readJson()

	f(urls)

	fmt.Scanln()
}

func f(urls Urls) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	for i := 0; i < len(urls); i++ {
		go doReq(i, urls[i].IP, urls[i].Port)
	}
}

func doReq(i int, ip string, port string) {
	t_start := time.Now()

	err := os.Setenv("HTTP_PROXY", "http://"+ip+":"+port)
	if err != nil {
		fmt.Println(i, err)
		return
	}

	proxyUrl, err := url.Parse("http://" + ip + ":" + port)
	if err != nil {
		fmt.Println(i, err)
		return
	}

	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: 10 * time.Second}
	res, err := myClient.Get("http://2children.ru/news/index?NewsSearch%5Bword%5D=%D0%B4%D0%B5%D1%82%D0%B8")
	if err != nil {
		fmt.Println(i, err)
		return
	}

	t_end := time.Now()
	fmt.Println(i, res.StatusCode, t_end.Sub(t_start))
}

func readJson() Urls {
	var urls Urls
	jsonFile, err := os.Open("proxies.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &urls)

	return urls
}
