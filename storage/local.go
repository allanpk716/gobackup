package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
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

	sysType := runtime.GOOS
	if sysType == "windows" {
		// windows系统
		logger.Info("Windows upload")
		logger.Info("ctx.archivePath", ctx.archivePath)
		logger.Info("ctx.destPath", ctx.destPath)
		logger.Info("fileKey", fileKey)
		err := MoveFile(ctx.archivePath, path.Join(ctx.destPath, fileKey))
		if err != nil {
			return err
		}
	} else {
		_, err = helper.Exec("cp", ctx.archivePath, ctx.destPath)
		if err != nil {
			return err
		}
	}

	logger.Info("Store successed", ctx.destPath)
	return nil
}

func (ctx *Local) delete(fileKey string) (err error) {
	sysType := runtime.GOOS
	if sysType == "windows" {
		// windows系统
		logger.Info("Windows delete")
		logger.Info("ctx.archivePath", ctx.archivePath)
		logger.Info("ctx.destPath", ctx.destPath)
		logger.Info("fileKey", fileKey)
		err := os.Remove(path.Join(ctx.destPath, fileKey))
		if err != nil {
			return err
		}
	} else {
		_, err = helper.Exec("rm", path.Join(ctx.destPath, fileKey))
	}
	return
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
