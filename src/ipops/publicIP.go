package ipops

/*
Get server networkinfo

*/

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// PublicIPInfo : PublicIPInfo struct
type PublicIPInfo struct {
	PublicIP string
	Info     string
}

//IPinfo : ip information
type IPinfo struct {
	Code int               `json:"code"`
	Data map[string]string `json:"data"`
}

// NewPublicIPInfo : get public IP info
func (publicIPinfo PublicIPInfo) NewPublicIPInfo() (ipinfo PublicIPInfo, err error) {
	ipinfo.PublicIP, err = PublicIP()
	if err != nil {
		return ipinfo, err
	}
	ipinfo.Info, err = IPInfo(ipinfo.PublicIP)
	return ipinfo, err

}

// GetPublicIP : get Public IP from http://www.ipip.net
func GetPublicIP(url, complieString string) (publicIP string, err error) {
	// create client and set timeout
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	// Get page
	resp, err := client.Get(url)
	// if error return empty string and error
	if err != nil {
		return "", err
	}
	// close open body
	defer resp.Body.Close()
	// []byte body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// commplie regex to get the ip string
	// re, _ := regexp.Compile(".*您当前的IP：(.*)</div>.*")
	re, _ := regexp.Compile(complieString)
	// find match from the html body byte
	results := re.FindSubmatch(body)
	// get the first match,if nomatch return  empty string and error no match
	if len(results) > 1 {
		publicIP = string(results[1])
		err = nil
	} else {
		publicIP = ""
		err = errors.New("no match")
	}
	return publicIP, err

}

// PublicIP : get public ip from internet
func PublicIP() (publicIP string, err error) {
	// get public from http://www.ipip.net and find match content
	publicIP, err = GetPublicIP("http://www.ipip.net", ".*您当前的IP：(.*)</div>.*")
	if err != nil {
		//  get public from https://ip.cn and find match content
		publicIP, err = GetPublicIP("https://ip.cn", ".*<p>您现在的 IP：<code>(.*)</code></p><p>所.*")
		if err != nil {
			//  get public from https://ip.cn and find match content

			publicIP, err = GetPublicIP("https://www.boip.net/api/myip", ".*(\\d\\.\\d\\.\\d\\.\\d).*")
		}
	}
	return publicIP, err

}

// IPInfo : get ip location-information
func IPInfo(ip string) (info string, err error) {
	// use taobao api to get IP location-information
	taobaoapi := "http://ip.taobao.com/service/getIpInfo.php?ip=" + ip
	response, err := http.Get(taobaoapi)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	// read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	// string json to mamp
	var dat IPinfo
	err = json.Unmarshal(body, &dat)
	if err != nil {
		return "", err
	}
	if dat.Code != 0 {
		return "", errors.New("search fail")
	}
	// return localtion-information string
	return dat.Data["country"] + " " + dat.Data["region"] + " " + dat.Data["city"] + " " + dat.Data["isp"], nil
}
