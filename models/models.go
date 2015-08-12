package models

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"../modules/polling"
	"../modules/setting"

	_ "github.com/go-sql-driver/mysql"
	"github.com/weisd/log"
	// "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var (
	XormEngines   map[string]*xorm.Engine
	masterPolling *polling.Polling
	slavePolling  *polling.Polling
)

func init() {
	XormEngines = make(map[string]*xorm.Engine)
}

// Round-Robin

func Master() *xorm.Engine {
	name := setting.Cfg.DBMaster[masterPolling.Index()]
	x, ok := XormEngines[name]
	if !ok {
		panic("Unknown master name %s", name)
	}

	log.Debug("Master use db name %s", name)

	return x
}

func Slave() *xorm.Engine {
	name := setting.Cfg.DBSlave[slavePolling.Index()]
	x, ok := XormEngines[name]
	if !ok {
		panic("Unknown Slave name %s", name)
	}

	log.Debug("Slave use db name %s", name)

	return x
}

func InitDatabaseConn() {

	if len(setting.Cfg.DBMaster) == 0 || len(setting.Cfg.DBSlave) == 0 {
		panic("setting.Cfg.DBMaster & DBSlave must be set ")
	}

	for name, conf := range setting.Cfg.DBs {
		x, err := newEngine(conf)
		if err != nil {
			log.Error(4, "newEngine failed name %s  %v", name, err)
			continue
		}

		XormEngines[name] = x
	}

	masterPolling = polling.NewPolling(len(setting.Cfg.DBMaster))
	masterPolling = polling.NewPolling(len(setting.Cfg.DBSlave))

	log.Debug("初始化 models done %v", XormEngines)
}

func newEngine(conf setting.DataBaseConfig) (*xorm.Engine, error) {
	cnnstr := ""
	switch conf.TYPE {
	case "mysql":
		if conf.HOST[0] == '/' { // looks like a unix socket
			cnnstr = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8",
				conf.USER, conf.PASSWD, conf.HOST, conf.NAME)
		} else {
			cnnstr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
				conf.USER, conf.PASSWD, conf.HOST, conf.NAME)
		}
	case "postgres":
		var host, port = "127.0.0.1", "5432"
		fields := strings.Split(conf.HOST, ":")
		if len(fields) > 0 && len(strings.TrimSpace(fields[0])) > 0 {
			host = fields[0]
		}
		if len(fields) > 1 && len(strings.TrimSpace(fields[1])) > 0 {
			port = fields[1]
		}
		cnnstr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			conf.USER, conf.PASSWD, host, port, conf.NAME, conf.SSL_MODE)
	case "sqlite3":
		os.MkdirAll(path.Dir(conf.PATH), os.ModePerm)
		cnnstr = "file:" + conf.PATH + "?cache=shared&mode=rwc"
	default:
		return nil, fmt.Errorf("Unknown database type: %s", conf.TYPE)
	}
	x, err := xorm.NewEngine(conf.TYPE, cnnstr)
	if err != nil {
		return nil, err
	}
	// 连接池的空闲数大小
	x.SetMaxIdleConns(conf.MaxIdle)
	// 最大打开连接数
	x.SetMaxOpenConns(conf.MaxOpen)

	x.ShowSQL = true
	// 则会在控制台打印出生成的SQL语句；
	x.ShowDebug = true
	// 则会在控制台打印调试信息；
	x.ShowErr = true
	// 则会在控制台打印错误信息；
	x.ShowWarn = true
	// 则会在控制台打印警告信息；

	logpath, _ := filepath.Abs(conf.LogPath)
	os.MkdirAll(path.Dir(logpath), os.ModePerm)
	// 日志
	f, err := os.Create(logpath)
	if err != nil {
		log.Error(4, "create xorm log file failed %v", err)
		return nil, err
	}
	// defer f.Close()
	x.Logger = xorm.NewSimpleLogger(f)
	return x, nil
}
