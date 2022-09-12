package log

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/astaxie/beego/logs"
)

var defaultLogConf = "./logger.json"

type configFile struct {
	Types     []interface{}     `json:"types"`
	Levels    map[string]string `json:"levels"`
	ChanSizes map[string]int    `json:"chanSizes"`
	LogPath   string            `json:"logPath"`
}

type categoryConf struct {
	level    int        // log level
	chanSize int        // chan size
	logConfs []*logConf // log writer configs for this category
}

type logConf struct {
	typ    string
	config string
}

var loggers = make(map[string]*logs.BeeLogger)
var loggerCategories = make(map[string]*categoryConf)
var LevelMap = map[string]int{
	"EMERGENCY":     logs.LevelEmergency,
	"ALERT":         logs.LevelAlert,
	"CRITICAL":      logs.LevelCritical,
	"CRIT":          logs.LevelCritical,
	"ERROR":         logs.LevelError,
	"ERR":           logs.LevelError,
	"WARNING":       logs.LevelWarning,
	"WARN":          logs.LevelWarning,
	"NOTICE":        logs.LevelNotice,
	"INFORMATIONAL": logs.LevelInformational,
	"INFO":          logs.LevelInformational,
	"DEBUG":         logs.LevelDebug,
}

func init() {
	// because flag.Parse must be called after all flags defined, but we want to get logConf in init func,
	// so here we parse logConf manually
	var logConfFile = getLogConf()
	if fileInfo, err := os.Stat(logConfFile); err == nil && !fileInfo.IsDir() {
		err = initLogger(logConfFile)
		if err != nil {
			panic(err)
		}
	}
	if _, ok := loggerCategories["default"]; !ok {
		loggerCategories["default"] = &categoryConf{
			level:    logs.LevelInformational,
			logConfs: []*logConf{{typ: "console", config: `{"color":false}`}},
		}
	}
	fmt.Println("init:", loggerCategories)
	Default = Get("default", true)
}

func getLogConf() string {
	for _, arg := range os.Args[1:] {
		fields := strings.Split(arg, "=")
		if len(fields) != 2 || (fields[0] != "-logconf" && fields[0] != "--logconf") || len(fields[1]) == 0 {
			continue
		}
		return fields[1]
	}
	return defaultLogConf
}

func getAsString(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}
	str, ok := v.(string)
	return str, ok
}

func checkLevels(levels map[string]string) error {
	for category, level := range levels {
		if _, ok := LevelMap[strings.ToUpper(level)]; !ok {
			return errors.New("unknown log level for " + category)
		}
	}
	return nil
}

// 防止float64在Marshal时产生科学计数，比如file类型的maxsize字段
func fixFloatingPoint(values map[string]interface{}) {
	for k, v := range values {
		switch value := v.(type) {
		case float64:
			values[k] = int(value)
		}
	}
}

func getLogPath(logPath interface{}) (string, error) {
	logDir, _ := getAsString(logPath)
	if len(logDir) == 0 {
		return "", nil
	}
	logDir, err := filepath.Abs(logDir)
	if err != nil {
		return "", err
	}
	return logDir, nil
}

func ensureFilePath(fileName string) {
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
}

func replaceVariable(content []byte) []byte {
	_, programName := filepath.Split(os.Args[0])
	content = regexp.MustCompile("\\$\\{programName\\}").ReplaceAll(content, []byte(programName))
	return content
}

func initLogger(configFileName string) error {
	content, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}
	var config = configFile{
		Levels:    make(map[string]string),
		ChanSizes: make(map[string]int),
	}
	err = json.Unmarshal(replaceVariable(content), &config)
	if err != nil {
		return err
	}
	err = checkLevels(config.Levels)
	if err != nil {
		return err
	}
	var types = config.Types
	for _, v := range types {
		var typeConfig = v.(map[string]interface{})
		fixFloatingPoint(typeConfig)
		typ, ok := getAsString(typeConfig["type"])
		if !ok {
			continue
		}
		if typ != "console" && typ != "file" {
			continue
		}
		delete(typeConfig, "type")
		category, ok := getAsString(typeConfig["category"])
		if !ok {
			category = "default"
		}
		delete(typeConfig, "category")
		if typ == "file" && len(config.LogPath) > 0 {
			fileName, _ := getAsString(typeConfig["filename"])
			fileName = path.Join(config.LogPath, fileName)
			ensureFilePath(fileName)
			typeConfig["filename"] = fileName
		}

		var categories = strings.Split(category, ",")
		for _, category := range categories {
			category = strings.TrimSpace(category)
			level, ok := LevelMap[strings.ToUpper(config.Levels[category])]
			if !ok {
				level = logs.LevelInformational
			}
			chanSize, ok := config.ChanSizes[category]
			if !ok {
				chanSize = 10000
			}
			rb, _ := json.Marshal(typeConfig)
			conf := &logConf{typ: typ, config: string(rb)}
			if confs, ok := loggerCategories[category]; ok {
				confs.logConfs = append(confs.logConfs, conf)
			} else {
				loggerCategories[category] = &categoryConf{level: level, chanSize: chanSize, logConfs: []*logConf{conf}}
			}
		}
	}
	return nil
}

// Init register command line flag, but only for a placeholder in command usage, as this package parses command line itself
func Init() {
	flag.String("logconf", defaultLogConf, "specify log config file")
}

// Get get a logger by category
func Get(name string, isAsync bool) *logs.BeeLogger {
	if beeLog, ok := loggers[name]; ok {
		return beeLog
	}
	configs, ok := loggerCategories[name]
	if !ok {
		configs = loggerCategories["default"]
	}
	fmt.Println(loggerCategories)
	beeLog := logs.NewLogger(int64(configs.chanSize))
	if isAsync {
		beeLog = beeLog.Async()
	}
	for _, lc := range configs.logConfs {
		beeLog.SetLogger(lc.typ, lc.config)
	}
	beeLog.SetLevel(configs.level)
	loggers[name] = beeLog
	return beeLog
}

// Close close all opened logger, call this when program exits
func Close() {
	for _, logger := range loggers {
		logger.Close()
	}
}

// Flush all opened logger
func Flush() {
	for _, logger := range loggers {
		logger.Flush()
	}
}
