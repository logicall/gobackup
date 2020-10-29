package config

import (
	"fmt"
	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

var (
	// Exist Is config file exist
	Exist bool
	// Models configs
	Models []ModelConfig
	// IsTest env
	IsTest bool
	// HomeDir of user
	HomeDir  string
	TempPath string
)

// ModelConfig for special case
type ModelConfig struct {
	Name         string
	DumpPath     string
	CompressWith SubConfig
	EncryptWith  SubConfig
	StoreWith    SubConfig
	Archive      *viper.Viper
	Databases    []SubConfig
	Storages     []SubConfig
	Viper        *viper.Viper
}

// SubConfig sub config info
type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

// loadConfig from:
// - ./gobackup.yml
// - ~/.gobackup/gobackup.yml
// - /etc/gobackup/gobackup.yml
func init() {
	viper.SetConfigType("yaml")

	IsTest = os.Getenv("GO_ENV") == "test"
	HomeDir = os.Getenv("HOME")

	if IsTest {
		viper.SetConfigName("gobackup_test")
		HomeDir = "../"
	} else {
		viper.SetConfigName("gobackup")
	}

	// ./gobackup.yml
	viper.AddConfigPath(".")
	if IsTest {
		viper.AddConfigPath("../")
	} else {
		// ~/.gobackup/gobackup.yml
		viper.AddConfigPath("$HOME/.gobackup") // call multiple times to add many search paths
		// /etc/gobackup/gobackup.yml
		viper.AddConfigPath("/etc/gobackup/") // path to look for the config file in
	}
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Load gobackup config faild", err)
		return
	}

	Exist = true
	Models = []ModelConfig{}
	for key := range viper.GetStringMap("models") {
		Models = append(Models, loadModel(key))
	}

	return
}

func loadModel(key string) (model ModelConfig) {
	model.Name = key
	if viper.GetString("models."+key+".temp_path") != "" {
		if helper.IsExistsPath(viper.GetString("models."+key+".temp_path")) == true {
			TempPath = path.Join(viper.GetString("models."+key+".temp_path"), "gobackup", fmt.Sprintf("%d", time.Now().UnixNano()))
		} else {
			helper.MkdirP(viper.GetString("models." + key + ".temp_path"))
			TempPath = path.Join(viper.GetString("models."+key+".temp_path"), "gobackup", fmt.Sprintf("%d", time.Now().UnixNano()))
		}
	} else {
		TempPath = path.Join(os.TempDir(), "gobackup", fmt.Sprintf("%d", time.Now().UnixNano()))
	}
	model.DumpPath = path.Join(TempPath, key)
	model.Viper = viper.Sub("models." + key)

	model.CompressWith = SubConfig{
		Type:  model.Viper.GetString("compress_with.type"),
		Viper: model.Viper.Sub("compress_with"),
	}

	model.EncryptWith = SubConfig{
		Type:  model.Viper.GetString("encrypt_with.type"),
		Viper: model.Viper.Sub("encrypt_with"),
	}

	model.StoreWith = SubConfig{
		Type:  model.Viper.GetString("store_with.type"),
		Viper: model.Viper.Sub("store_with"),
	}

	model.Archive = model.Viper.Sub("archive")

	loadDatabasesConfig(&model)
	loadStoragesConfig(&model)

	return
}

func loadDatabasesConfig(model *ModelConfig) {
	subViper := model.Viper.Sub("databases")
	for key := range model.Viper.GetStringMap("databases") {
		dbViper := subViper.Sub(key)
		model.Databases = append(model.Databases, SubConfig{
			Name:  key,
			Type:  dbViper.GetString("type"),
			Viper: dbViper,
		})
	}
}

func loadStoragesConfig(model *ModelConfig) {
	subViper := model.Viper.Sub("storages")
	for key := range model.Viper.GetStringMap("storages") {
		dbViper := subViper.Sub(key)
		model.Storages = append(model.Storages, SubConfig{
			Name:  key,
			Type:  dbViper.GetString("type"),
			Viper: dbViper,
		})
	}
}

// GetModelByName get model by name
func GetModelByName(name string) (model *ModelConfig) {
	for _, m := range Models {
		if m.Name == name {
			model = &m
			return
		}
	}
	return
}

// GetDatabaseByName get database config by name
func (model *ModelConfig) GetDatabaseByName(name string) (subConfig *SubConfig) {
	for _, m := range model.Databases {
		if m.Name == name {
			subConfig = &m
			return
		}
	}
	return
}
