package log

import (
	"encoding/json"
	"fmt"
	"os"

	"../../modules/setting"

	baseLog "github.com/weisd/log"
)

type loggerMap map[string]*baseLog.Logger

var LogsMap map[string]loggerMap

var logLevels = map[string]string{
	"Trace":    "0",
	"Debug":    "1",
	"Info":     "2",
	"Warn":     "3",
	"Error":    "4",
	"Critical": "5",
}

func InitLogs() {

	LogsMap = make(map[string]loggerMap)

	for name, v := range setting.Cfg.Logs {
		loggerMaper := make(map[string]*baseLog.Logger)
		for _, conf := range v {

			if !conf.ENABLE {
				continue
			}
			level, ok := logLevels[conf.LEVEL]
			if !ok {
				baseLog.Fatal(4, "Unknown log level: %s", conf.LEVEL)
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
					baseLog.Error(4, "json.Marshal(conf.RECEIVERS) err %v", err)
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

			baseLog.Info(str)

			loggerMaper[conf.MODE] = baseLog.NewCustomLogger(conf.BUFFER_LEN, conf.MODE, str)
			baseLog.Info("Log Mode: %s(%s)", conf.MODE, conf.LEVEL)

		}

		LogsMap[name] = loggerMaper
	}

	Info("logs 初始化完成 map %v", LogsMap)

}

func (loggers loggerMap) Trace(format string, v ...interface{}) {

	for _, logger := range loggers {
		logger.Trace(format, v...)
	}
}

func (loggers loggerMap) Debug(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Debug(format, v...)
	}
}

func (loggers loggerMap) Info(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Info(format, v...)
	}
}

func (loggers loggerMap) Warn(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Warn(format, v...)
	}
}

func (loggers loggerMap) Error(skip int, format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Error(skip, format, v...)
	}
}

func (loggers loggerMap) Critical(skip int, format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Critical(skip, format, v...)
	}
}

func (loggers loggerMap) Fatal(skip int, format string, v ...interface{}) {
	loggers.Error(skip, format, v...)
	for _, l := range loggers {
		l.Close()
	}
	os.Exit(1)
}

func Get(name string) loggerMap {
	l, ok := LogsMap[name]
	if !ok {
		baseLog.Fatal(1, "Unknown log %s", name)
		return nil
	}

	return l
}

func Trace(format string, v ...interface{}) {

	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Trace(format, v...)
	// for _, logger := range loggers {
	// 	logger.Trace(format, v...)
	// }
}

func Debug(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Debug(format, v...)
	// for _, logger := range loggers {
	// 	logger.Debug(format, v...)
	// }
}

func Info(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Info(format, v...)
	// for _, logger := range loggers {
	// 	logger.Info(format, v...)
	// }
}

func Warn(format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Warn(format, v...)
	// for _, logger := range loggers {
	// 	logger.Warn(format, v...)
	// }
}

func Error(skip int, format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Error(skip, format, v...)
	// for _, logger := range loggers {
	// 	logger.Error(skip, format, v...)
	// }
}

func Critical(skip int, format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Critical(skip, format, v...)
	// for _, logger := range loggers {
	// 	logger.Critical(skip, format, v...)
	// }
}

func Fatal(skip int, format string, v ...interface{}) {
	loggers, ok := LogsMap["default"]
	if !ok {
		return
	}
	loggers.Fatal(skip, format, v...)
	// Error(skip, format, v...)
	// for _, l := range loggers {
	// 	l.Close()
	// }
	// os.Exit(1)
}

func Close() {
	for _, loggers := range LogsMap {
		for _, l := range loggers {
			l.Close()
		}
	}

}
