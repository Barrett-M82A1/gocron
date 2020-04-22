package setting

import (
	"errors"
	"os"
	"strconv"

	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"gopkg.in/ini.v1"
)

const DefaultSection = "default"

type Setting struct {
	Db struct {
		Engine       string
		Host         string
		Port         int
		User         string
		Password     string
		Database     string
		Prefix       string
		Charset      string
		MaxIdleConns int
		MaxOpenConns int
	}
	AllowIps      string
	AppName       string
	ApiKey        string
	ApiSecret     string
	ApiSignEnable bool

	EnableTLS bool
	CAFile    string
	CertFile  string
	KeyFile   string

	ConcurrencyQueue int
	AuthSecret       string
}

func ReadEnv(key string, value string) string {
	if os.Getenv(key) == "" {
		return value
	}

	return os.Getenv(key)
}

// 读取配置
func Read(filename string) (*Setting, error) {

	var s Setting

	// 如果存在环境变量配置优先读取
	if ReadEnv("ENV_CONFIG", "off") == "open" {
		s.Db.Engine = ReadEnv("DB_ENGINE", "mysql")
		s.Db.Host = ReadEnv("DB_HOST", "127.0.0.1")
		s.Db.Port, _ = strconv.Atoi(ReadEnv("DB_PORT", "3306"))
		s.Db.User = ReadEnv("DB_USER", "")
		s.Db.Password = ReadEnv("DB_PASSWORD", "")
		s.Db.Database = ReadEnv("DB_DATABASE", "gocron")
		s.Db.Prefix = ReadEnv("DB_PREFIX", "")
		s.Db.Charset = ReadEnv("DB_CHARSET", "utf8")
		s.Db.MaxIdleConns, _ = strconv.Atoi(ReadEnv("DB_MAXIDLECONNS", "30"))
		s.Db.MaxOpenConns, _ = strconv.Atoi(ReadEnv("DB_MAXOPENCONNS", "100"))

		s.AllowIps = ReadEnv("ALLOW_IPS", "")
		s.AppName = ReadEnv("APP_NAME", "定时任务管理系统")
		s.ApiKey = ReadEnv("API_KEY", "")
		s.ApiSecret = ReadEnv("API_SECRET", "")
		s.ApiSignEnable, _ = strconv.ParseBool(ReadEnv("API_SIGN_ENABLE", "true"))
		s.ConcurrencyQueue, _ = strconv.Atoi(ReadEnv("CONCURRENCY_QUEUE", "500"))
		s.AuthSecret = ReadEnv("AUTH_SECRET", "")
		if s.AuthSecret == "" {
			s.AuthSecret = utils.RandAuthToken()
		}

		s.EnableTLS, _ = strconv.ParseBool(ReadEnv("ENABLE_TLS", "false"))
		s.CAFile = ReadEnv("CA_FILE", "")
		s.CertFile = ReadEnv("CERT_FILE", "")
		s.KeyFile = ReadEnv("KEY_FILE", "")

		if s.EnableTLS {
			if !utils.FileExist(s.CAFile) {
				logger.Fatalf("failed to read ca cert file: %s", s.CAFile)
			}

			if !utils.FileExist(s.CertFile) {
				logger.Fatalf("failed to read client cert file: %s", s.CertFile)
			}

			if !utils.FileExist(s.KeyFile) {
				logger.Fatalf("failed to read client key file: %s", s.KeyFile)
			}
		}

		return &s, nil
	}

	config, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}
	section := config.Section(DefaultSection)

	s.Db.Engine = section.Key("db.engine").MustString("mysql")
	s.Db.Host = section.Key("db.host").MustString("127.0.0.1")
	s.Db.Port = section.Key("db.port").MustInt(3306)
	s.Db.User = section.Key("db.user").MustString("")
	s.Db.Password = section.Key("db.password").MustString("")
	s.Db.Database = section.Key("db.database").MustString("gocron")
	s.Db.Prefix = section.Key("db.prefix").MustString("")
	s.Db.Charset = section.Key("db.charset").MustString("utf8")
	s.Db.MaxIdleConns = section.Key("db.max.idle.conns").MustInt(30)
	s.Db.MaxOpenConns = section.Key("db.max.open.conns").MustInt(100)

	s.AllowIps = section.Key("allow_ips").MustString("")
	s.AppName = section.Key("app.name").MustString("定时任务管理系统")
	s.ApiKey = section.Key("api.key").MustString("")
	s.ApiSecret = section.Key("api.secret").MustString("")
	s.ApiSignEnable = section.Key("api.sign.enable").MustBool(true)
	s.ConcurrencyQueue = section.Key("concurrency.queue").MustInt(500)
	s.AuthSecret = section.Key("auth_secret").MustString("")
	if s.AuthSecret == "" {
		s.AuthSecret = utils.RandAuthToken()
	}

	s.EnableTLS = section.Key("enable_tls").MustBool(false)
	s.CAFile = section.Key("ca_file").MustString("")
	s.CertFile = section.Key("cert_file").MustString("")
	s.KeyFile = section.Key("key_file").MustString("")

	if s.EnableTLS {
		if !utils.FileExist(s.CAFile) {
			logger.Fatalf("failed to read ca cert file: %s", s.CAFile)
		}

		if !utils.FileExist(s.CertFile) {
			logger.Fatalf("failed to read client cert file: %s", s.CertFile)
		}

		if !utils.FileExist(s.KeyFile) {
			logger.Fatalf("failed to read client key file: %s", s.KeyFile)
		}
	}

	return &s, nil
}

// 写入配置
func Write(config []string, filename string) error {
	if len(config) == 0 {
		return errors.New("参数不能为空")
	}
	if len(config)%2 != 0 {
		return errors.New("参数不匹配")
	}

	file := ini.Empty()

	section, err := file.NewSection(DefaultSection)
	if err != nil {
		return err
	}
	for i := 0; i < len(config); {
		_, err = section.NewKey(config[i], config[i+1])
		if err != nil {
			return err
		}
		i += 2
	}
	err = file.SaveTo(filename)

	return err
}
