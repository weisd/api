package setting

import (
	// "encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/koding/multiconfig"
	"github.com/weisd/com"
	"github.com/weisd/log"
)

var (
	ConfigPath string
	UserPath   string
	WorkDir    string
)

var Cfg *Configs

type Configs struct {
	Debug    bool
	Logs     map[string][]LogConfig
	Hosts    map[string]HostsConf
	DBs      map[string]DataBaseConfig
	Redis    map[string]RedisConfig
	DBMaster []string
	DBSlave  []string
}

func (c *Configs) GetUserDir(hostname string) string {
	if info, ok := c.Hosts[hostname]; ok {
		return info.Config
	}

	return "/"
}

type HostsConf struct {
	Config   string
	Hostname string
}

func init() {
	Cfg = newCfg()
	log.NewLogger(0, "console", `{"level": 0}`)
}

func newCfg() *Configs {
	return &Configs{}
}

func ResetConfig() {
	Cfg = newCfg()
}

func InitConfig() {
	var err error

	WorkDir, err = com.WorkDir()
	if err != nil {
		panic(err)
	}

	ConfigPath = path.Join(WorkDir, "conf")

	Hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println(Hostname)

	err = initHosts()
	if err != nil {
		panic(err)
	}

	fmt.Println(Cfg)

	UserPath = path.Join(ConfigPath, Cfg.GetUserDir(Hostname))

	fmt.Println(ConfigPath, UserPath)

	InitServices()
}

func appConfig() {
	configPath := getUserConfigFile("app.toml")

	m := multiconfig.NewWithPath(configPath)
	err := m.Load(Cfg)
	if err != nil {
		panic(err)
	}
}

func initHosts() error {
	m := multiconfig.NewWithPath("conf/hosts.toml")
	err := m.Load(Cfg)
	if err != nil {
		return err
	}

	return nil
}

func InitServices() {
	appConfig()
	newLogService()
	newDBService()
	newRedisService()
}

//// log /////
type LogConfig struct {
	ENABLE bool
	MODE   string `required:"true"`
	LEVEL  string `required:"true"`

	BUFFER_LEN int64 `default:"10000"`

	// file
	FILE_NAME      string `default:"/tmp/dadalog.log"`
	LOG_ROTATE     bool   `default:"true"`
	MAX_LINES      int    `default:"1000000"`
	MAX_SIZE_SHIFT int    `default:"28"`
	DAILY_ROTATE   bool   `default:"true"`
	MAX_DAYS       int    `default:"7"`

	// conn
	RECONNECT_ON_MSG bool
	RECONNECT        bool
	PROTOCOL         string
	ADDR             string

	// smtp
	USER      string
	PASSWD    string
	HOST      string
	RECEIVERS []string
	SUBJECT   string

	// database
	DRIVER string
	CONN   string
}

func newLogService() {
	f := "logs.toml"
	configPath := getUserConfigFile(f)

	m := multiconfig.NewWithPath(configPath)
	err := m.Load(Cfg)
	if err != nil {
		panic(err)
	}

	log.Info("k \n %v", Cfg.Logs)

}

func getUserConfigFile(f string) string {
	configPath := path.Join(UserPath, f)
	if !com.IsFile(configPath) {
		configPath = path.Join(ConfigPath, f)
	}
	return configPath
}

//// log /////

//// database ////
type DataBaseConfig struct {
	NAME     string
	TYPE     string
	HOST     string
	DB       string
	USER     string
	PASSWD   string
	SSL_MODE string
	PATH     string

	MaxIdle int
	MaxOpen int

	ShowSQL   bool
	ShowDebug bool
	ShowError bool
	ShowWarn  bool

	LogPath string
}

func newDBService() {
	configPath := getUserConfigFile("database.toml")

	m := multiconfig.NewWithPath(configPath)
	err := m.Load(Cfg)
	if err != nil {
		panic(err)
	}
}

//// database ////

//// redis ////
type RedisConfig struct {
	ADDR   string
	PASSWD string

	MAX_IDLE    int  `default:"10"`
	MAX_ACTIVE  int  `default:"0"`
	IdleTimeout int  `default:"30"`
	Wait        bool `default:"true"`
}

func newRedisService() {
	configPath := getUserConfigFile("redis.toml")

	m := multiconfig.NewWithPath(configPath)
	err := m.Load(Cfg)
	if err != nil {
		panic(err)
	}
}
