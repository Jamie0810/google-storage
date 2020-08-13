package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server        Server        `mapstructure:"server"`
	Database      Database      `mapstructure:"database"`
	GoogleStorage GoogleStorage `mapstructure:"google_storage"`
	Log           Log           `mapstructure:"log"`
	IsAuthorized  bool          `mapstructure:"is_authorized"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Database struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         uint   `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	InstanceName string `mapstructure:"instance_name"`
}

type GoogleStorage struct {
	ProjectID      string `mapstructure:"project_id"`
	StorageKeyFile string `mapstructure:"storage_key_file"`
	BucketName     string `mapstructure:"bucket_name"`
	URL            string `mapstructure:"url"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func InitConfig(configPath string) (config Config, err error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	/* default */
	v.SetDefault("log_level", "INFO")
	v.SetDefault("log_format", "console")

	defaultPath := `./configs`

	if configPath == "" {
		configPath = defaultPath
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath(configPath)

	files, _ := ioutil.ReadDir(configPath)
	index := 0

	for _, file := range files {
		if filepath.Ext("./"+file.Name()) != ".yaml" && filepath.Ext("./"+file.Name()) != ".yml" {
			continue
		}

		v.SetConfigName(file.Name())
		var err error
		if index == 0 {
			err = v.ReadInConfig()
		} else {
			err = v.MergeInConfig()
		}
		if err == nil {
			index++
		}
	}

	if err = v.Unmarshal(&config); err != nil {
		return
	}

	return
}
