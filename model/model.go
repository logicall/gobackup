package model

import (
	"github.com/spf13/viper"
	"os"
	"path"

	"github.com/huacnlee/gobackup/archive"
	"github.com/huacnlee/gobackup/compressor"
	"github.com/huacnlee/gobackup/config"
	"github.com/huacnlee/gobackup/database"
	"github.com/huacnlee/gobackup/encryptor"
	"github.com/huacnlee/gobackup/logger"
	"github.com/huacnlee/gobackup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

// Perform model
func (ctx Model) Perform() {
	logger.Info("======== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath+"\n")
	defer ctx.cleanup()

	err := database.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	if ctx.Config.Archive != nil {
		err = archive.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	archivePath, err := compressor.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	archivePath, err = encryptor.Run(archivePath, ctx.Config)
	if err != nil {
		logger.Error(err)
		return
	}

	err = storage.Run(ctx.Config, archivePath)
	if err != nil {
		logger.Error(err)
		return
	}

}

// Cleanup model temp files
func (ctx Model) cleanup() {
	tempPath := path.Join(viper.GetString("models."+ctx.Config.Name+".temp_path"), "gobackup")

	//logger.Info("Cleanup temp dir:" + config.TempPath + "...\n")
	logger.Info("Cleanup temp dir:" + tempPath + "...\n")

	err := os.RemoveAll(tempPath)
	if err != nil {
		logger.Error("Cleanup temp dir "+tempPath+" error:", err)
	}
	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")
}
