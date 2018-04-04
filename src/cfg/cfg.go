package cfg

import (
	"io/ioutil"
	"strings"
)

// InitConfig : read configFile to return map config and error info
func InitConfig(configPath string) (config map[string]string, err error) {
	config = make(map[string]string)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// spik empty line , remark line and the line without =
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") || (strings.Contains(line, "=") == false) {
			continue
		}
		// wipe out remark content
		line = strings.Split(line, "#")[0]
		// get key and value
		item := strings.Split(line, "=")
		if strings.TrimSpace(item[0]) == "" {
			continue
		}
		config[strings.TrimSpace(item[0])] = strings.TrimSpace(item[1])
	}
	return config, nil
}
