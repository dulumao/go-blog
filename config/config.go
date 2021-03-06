package config

import (
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/go-ini/ini"
	"github.com/zzayne/go-blog/utils"
)

var (
	cfg *ini.File
)

type dbConfig struct {
	Dialect      string
	Database     string
	User         string
	Password     string
	Host         string
	Port         int
	Name         string
	Charset      string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

type appConfig struct {
	PageSize    int
	TokenSecret string
	PassSalt    string
}
type serverConfig struct {
	Port               int
	Env                string
	APIPrefix          string
	MaxMultipartMemory int
	LogDir             string
	LogFile            string
}

func init() {
	var err error
	cfg, err = ini.Load("./app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	loadDBConifg()
	loadAppConfig()
	loadServerConfig()
}

// DBConfig 数据库相关配置
var DBConfig dbConfig

func loadDBConifg() {
	err := cfg.Section("database").MapTo(&DBConfig)
	if err != nil {
		log.Fatalf("Fail to load database config: %v", err)
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.Charset)
	DBConfig.URL = url
}

//AppConfig 服务内部相关配置
var AppConfig appConfig

func loadAppConfig() {
	err := cfg.Section("app").MapTo(&AppConfig)
	if err != nil {
		log.Fatalf("Fail to load app config: %v", err)
	}
}

//ServerConfig 服务运行环境配置
var ServerConfig serverConfig

func loadServerConfig() {
	err := cfg.Section("server").MapTo(&ServerConfig)
	if err != nil {
		log.Fatalf("Fail to load server config: %v", err)
	}
	sep := string(os.PathSeparator)
	execPath, _ := os.Getwd()
	length := utf8.RuneCountInString(execPath)
	lastChar := execPath[length-1:]
	if lastChar != sep {
		execPath = execPath + sep
	}
	ymdStr := utils.GetTodayYMD("-")

	if ServerConfig.LogDir == "" {
		ServerConfig.LogDir = execPath
	} else {
		length := utf8.RuneCountInString(ServerConfig.LogDir)
		lastChar := ServerConfig.LogDir[length-1:]
		if lastChar != sep {
			ServerConfig.LogDir = ServerConfig.LogDir + sep
		}
	}
	ServerConfig.LogFile = ServerConfig.LogDir + ymdStr + ".log"
}
