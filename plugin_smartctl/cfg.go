package main

import (
	"encoding/json"
	"fmt"
	"github.com/toolkits/file"
	"os"
	"sync"
)

type GlobalConfig struct {
	IP       string            `json:"ip"`
	Path     string            `json:"path"`
	Count    int               `json:"count"`
	Plugins  []string          `json:"plugins"`
	Enabled  map[string]bool   `json:"enabled"`
	Dir      map[string]string `json:"dir"`
	Interval map[string]int    `json:"interval"`
	Addition map[string]string `json:"additon"`
}

type Plugin struct {
	name     string
	enabled  bool
	dir      string
	interval int
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func count() int {
	count := Config().Count
	if count == 0 {
		fmt.Println("No Plugins")
		os.Exit(1)
	}

	return count
}

func IP() string {
	ip := Config().IP
	if ip == "" {
		LogRun(plu_name + "*****" + "get hostip failed")
		return ""
	}

	return ip
}

func path() string {
	path := Config().Path
	if path == "" {
		fmt.Println("the path of falcon-plugins is incorrect")
		os.Exit(1)
	}

	return path
}

func dir(name string) string {
	Dir := Config().Dir
	dir := Dir[name]
	if dir == "" {
		return ""
	}

	return dir
}

func plugins() []Plugin {
	plugins := []Plugin{}
	var plugin Plugin
	plu_name := Config().Plugins
	plu_enabled := Config().Enabled
	plu_dir := Config().Dir
	plu_interval := Config().Interval
	for _, name := range plu_name {
		plugin.name = name
		plugin.enabled = plu_enabled[name]
		plugin.dir = plu_dir[name]
		plugin.interval = plu_interval[name]
		plugins = append(plugins, plugin)
	}

	return plugins
}

func ParseConfig(cfg string) {
	if cfg == "" {
		fmt.Println("use -c to specify configuration file")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		fmt.Println("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		fmt.Println("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	fmt.Println("read config file:", cfg, "successfully")
}
