package setting

import (
	"encoding/json"
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
	Logs     []LogConfig
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
	Cfg = new(Configs)
	log.NewLogger(0, "console", `{"level": 0}`)
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

	appConfig()
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

var logLevels = map[string]string{
	"Trace":    "0",
	"Debug":    "1",
	"Info":     "2",
	"Warn":     "3",
	"Error":    "4",
	"Critical": "5",
}

func newLogService() {
	f := "logs.toml"
	configPath := getUserConfigFile(f)

	m := multiconfig.NewWithPath(configPath)
	err := m.Load(Cfg)
	if err != nil {
		panic(err)
	}

	for _, conf := range Cfg.Logs {
		if !conf.ENABLE {
			continue
		}
		level, ok := logLevels[conf.LEVEL]
		if !ok {
			log.Fatal(4, "Unknown log level: %s", conf.LEVEL)
		}

		str := ""
		switch conf.MODE {
		case "console":
			str = fmt.Sprintf(`{"level":%s}`, level)
		case "file":
			str = fmt.Sprintf(
				`{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`,
				level,
				conf.FILE_NAME,
				conf.LOG_ROTATE,
				conf.MAX_LINES,
				1<<uint(conf.MAX_SIZE_SHIFT),
				conf.DAILY_ROTATE, conf.MAX_DAYS,
			)
		case "conn":
			str = fmt.Sprintf(`{"level":%s,"reconnectOnMsg":%v,"reconnect":%v,"net":"%s","addr":"%s"}`, level,
				conf.RECONNECT_ON_MSG,
				conf.RECONNECT,
				conf.PROTOCOL,
				conf.ADDR)
		case "smtp":

			tos, err := json.Marshal(conf.RECEIVERS)
			if err != nil {
				log.Error(4, "json.Marshal(conf.RECEIVERS) err %v", err)
				continue
			}

			str = fmt.Sprintf(`{"level":%s,"username":"%s","password":"%s","host":"%s","sendTos":%s,"subject":"%s"}`, level,
				conf.USER,
				conf.PASSWD,
				conf.HOST,
				tos,
				conf.SUBJECT)
		case "database":
			str = fmt.Sprintf(`{"level":%s,"driver":"%s","conn":"%s"}`, level,
				conf.DRIVER,
				conf.CONN)
		default:
			continue
		}

		log.Info(str)
		log.NewLogger(conf.BUFFER_LEN, conf.MODE, str)
		log.Info("Log Mode: %s(%s)", conf.MODE, conf.LEVEL)
	}

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
