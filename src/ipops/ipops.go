package ipops

import (
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// SavePublicIP : save info information to datfile
func SavePublicIP(ip, info, datfile string) (err error) {
	record := ip + "|" + info + "|" + time.Now().Format("2006-01-02 15:04:05") + "\n"
	f, err := os.OpenFile(datfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(record))
	return err
}

// GetLastIP : from last ip file
func GetLastIP(datfile string) string {
	contents, err := ioutil.ReadFile(datfile)
	if err != nil {
		return ""
	}
	return strings.Split(string(contents), "|")[0]
}

// SaveLastIP : save last ip information to datfile
func SaveLastIP(ip, info, tempfile string) (err error) {
	record := ip + "|" + info + "|" + time.Now().Format("2006-01-02 15:04:05")
	err = ioutil.WriteFile(tempfile, []byte(record), 0644)
	return err
}
