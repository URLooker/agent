package g

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/toolkits/file"
)

type WebConf struct {
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
}

type GlobalConfig struct {
	Debug    bool     `json:"debug"`
	Idc      string   `json:"idc"`
	Hostname string   `json:"hostname"`
	Worker   int      `json:"worker"`
	ItemAddr string   `json:"itemAddr"`
	Dns      string   `json:"dns"`
	Web      *WebConf `json:"web"`
}

var (
	Config *GlobalConfig
)

func Hostname() (string, error) {
	hostname := Config.Hostname
	if hostname != "" {
		return hostname, nil
	}

	return os.Hostname()
}

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	Config = &c

	log.Println("load configuration file", cfg, "successfully")
	return nil
}
