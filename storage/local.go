package storage

import (
	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
	"path"
	"strings"
)

// Local storage
//
// type: local
// path: /data/backups
type Local struct {
	Base
	destPath string
}

func (ctx *Local) open() (err error) {
	ctx.destPath = ctx.model.StoreWith.Viper.GetString("path")
	helper.MkdirP(ctx.destPath)
	return
}

func (ctx *Local) close() {}

func (ctx *Local) upload(fileKey string) (err error) {

	if ctx.Base.sync == "rsync" {
		rsyncArgs := []string{}
		if len(ctx.Base.syncLimit) > 0 {
			rsyncArgs = append(rsyncArgs, "--bwlimit="+ctx.Base.syncLimit)
		}
		rsyncCMD := "rsync " + strings.Join(rsyncArgs, " ")
		_, err = helper.Exec(rsyncCMD, ctx.archivePath, ctx.destPath)
	} else if ctx.Base.sync == "cp" {
		_, err = helper.Exec("cp", ctx.archivePath, ctx.destPath)
	}

	if err != nil {
		return err
	}
	logger.Info("Store successed", ctx.destPath)
	return nil
}

func (ctx *Local) delete(fileKey string) (err error) {
	_, err = helper.Exec("rm", path.Join(ctx.destPath, fileKey))
	return
}
