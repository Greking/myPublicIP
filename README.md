# myPublicIP
间断检测公网IP，如果有变，就发送邮件提醒。



### linux安装：
需要有go环境，设置好GOPATH之后，在代码目录运行
```
bash install.sh
```
### windows安装:
需要有go环境，设置好GOPATH之后，在代码目录运行
```
install.bat
```


### 使用方法：
./publicIP -h
```
Usage of bin/publicIP:
  -default-config
    	show default config contetns
  -f string
    	configPath (default:bin/publicIP.cfg)
  -v	show version 

```
需要先创建一个配置文件，配置内容可以使用-default-config参数查看
```
fromEmail = xxxx@xx.com   #发送邮件的邮箱
password = xxxxx #邮箱密码
SMTPServer = xxxx:25 #smtp服务器
toEmails = xxxxxx #如果有多个邮箱以，为分隔符
saveLastIP = true  #是否记录最后的IP到文件
savePublicIP = true #是否记录每次获取到的IP
interval = 1m    #合法的单位有"s"、"m"、"h" ,默认是1m
```
启动的时候如果没有指定配置文件的路径，默认是找跟运行程序同级的目录下的publicIP.cfg
