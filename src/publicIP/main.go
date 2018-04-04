package main

import (
	"cfg"
	"flag"
	"fmt"
	"ipops"
	"log"
	"os"
	"sendmail"
	"strings"
	"time"
)

// VERSION : the version program
const VERSION = "v1.0"

func filePath() string {
	dirname := strings.Split(os.Args[0], string(os.PathSeparator))
	return strings.Join(dirname[:len(dirname)-1], string(os.PathSeparator))
}

func main() {

	v := flag.Bool("v", false, "show version ")
	f := flag.String("f", "", "configPath (default:"+filePath()+string(os.PathSeparator)+"publicIP.cfg)")
	defaultConfig := flag.Bool("default-config", false, "show default config contetns")
	flag.Parse()
	if *v {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if *defaultConfig {
		configConets := `fromEmail = xxxx@xx.com   #发送邮件的邮箱
password = xxxxx #邮箱密码
SMTPServer = xxxx:25 #smtp服务器
toEmails = xxxxxx #如果有多个邮箱以，为分隔符
saveLastIP = true  #是否记录最后的IP到文件
savePublicIP = true #是否记录每次获取到的IP
interval = 1m    #合法的单位有"s"、"m"、"h" ,默认是1m
`
		fmt.Println(configConets)
		os.Exit(0)
	}
	configPath := filePath() + string(os.PathSeparator) + "publicIP.cfg"
	if *f != "" {
		configPath = *f
	}

	lastIPfile := filePath() + string(os.PathSeparator) + "lastIP.txt"
	publicIPfile := filePath() + string(os.PathSeparator) + "publicIP.txt-" + time.Now().Format("2006-01-02")
	// record last IP

	// get from config or set default
	// interval, _ := time.ParseDuration("1m")
	// lastIP := ""
	saveLastIP := true
	savePublicIP := true
	subject := ""
	body := ""
	config, err := cfg.InitConfig(configPath)
	if err != nil {
		panic(err)
	}
	logfilename := filePath() + string(os.PathSeparator) + "run.log-" + time.Now().Format("2006-01-02")
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0664)
	defer logfile.Close()
	if err != nil {
		panic("can not open logfile : " + logfilename)
	}
	logger := log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)
	fromEmail := config["fromEmail"]
	password := config["password"]
	SMTPServer := config["SMTPServer"]
	toEmails := config["toEmails"]
	interval, err := time.ParseDuration(config["interval"])
	if err != nil {
		logger.Println(err)
		logger.Println("interval is invalid,set default 1m")
		interval, _ = time.ParseDuration("1m")
	}
	if config["saveLastIP"] != "true" {
		saveLastIP = false
	}
	if config["savePublicIP"] != "true" {
		saveLastIP = false
	}
	lastIP := ""
	// saveLastIP := true
	// savePublicIP := true
	// set logfile

	logger.Println("start...")
	defer func() {
		fmt.Println("stop")
	}()
	for {
		// start get
		go func() {
			logger.Println("start to get public ip ")
			publicip, err := ipops.PublicIP()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			logger.Println("get location information")
			locationInfo, _ := ipops.IPInfo(publicip)
			logger.Println(publicip, locationInfo)
			// if last ip is empty  and saveLastIP is true ,then get last from lastIPfile
			if lastIP == "" && saveLastIP {
				lastIP = ipops.GetLastIP(lastIPfile)
			}
			// fmt.Println(lastIP)
			if lastIP != publicip {
				go func() {
					subject = "IP切换提醒" + time.Now().Format("2006-01-02 15:04:05")
					body = "<p>当前IP：" + publicip + "</p>" + "<p>上次IP：" + lastIP + "</p>"
					logger.Println(fromEmail, toEmails, subject, body)
					err = sendmail.SendToMail(fromEmail, password, SMTPServer, toEmails, subject, body, "html")
					if err != nil {
						logger.Println("sendmail Fail", err)
						return
					}
					lastIP = publicip
				}()
			}
			// save last IP
			if saveLastIP == true {
				go func() {
					logger.Println("save last ip ")
					err = ipops.SaveLastIP(publicip, locationInfo, lastIPfile)
					if err != nil {
						logger.Println("save last ip fail")
						logger.Println(err)
					}
				}()
			}
			//save public IP
			if savePublicIP == true {
				go func() {
					// switch logfile
					if publicIPfile != filePath()+string(os.PathSeparator)+"publicIP.txt-"+time.Now().Format("2006-01-02") {
						publicIPfile = filePath() + string(os.PathSeparator) + "publicIP.txt-" + time.Now().Format("2006-01-02")
					}
					logger.Println("save public ip ")
					err = ipops.SavePublicIP(publicip, locationInfo, publicIPfile)
					if err != nil {
						logger.Println("save public ip fail")
						logger.Println(err)
					}
				}()
			}
		}()

		// sleep interval time
		time.Sleep(interval)

		// switch logfile
		if logfilename != filePath()+string(os.PathSeparator)+"run.log-"+time.Now().Format("2006-01-02") {
			logfilename = filePath() + string(os.PathSeparator) + "run.log-" + time.Now().Format("2006-01-02")
			logfile.Close()
			logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
			defer logfile.Close()
			if err != nil {
				panic("can not open logfile : " + logfilename)
			}
			logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)
		}

	}

}
