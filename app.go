package main

import(
	"fmt"
	"net/http"
	"time"
	"os"
	"net/url"
	"encoding/json"
    "io/ioutil"
)

type Urls []struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

var urls Urls

func main(){
	readJson()

	for i:=0; i<len(urls); i++ {
		go doReq(i)
	}
	fmt.Scanln()
}

func doReq(i int){
	t_start:= time.Now()

	os.Setenv("HTTP_PROXY", "http://" + urls[i].IP + ":" + urls[i].Port)
	
	proxyUrl, err := url.Parse("http://" + urls[i].IP + ":" + urls[i].Port)
	if err != nil {
		fmt.Println(i, err)
	}

	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	res, err := myClient.Get("http://2children.ru/news/index")
	if err != nil {
		fmt.Println(i, err)
	}

	t_end := time.Now()
	fmt.Println(i, res.StatusCode, t_end.Sub(t_start))
}

func readJson(){
	jsonFile, err := os.Open("proxies.json")
    if err != nil {
        fmt.Println(err)
    }

    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &urls)
}