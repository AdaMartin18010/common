package common

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// 执行源文件(or main.go)的编译之前的 文件目录
func CompiledExectionFilePath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		log.Default().Printf("Can not get current file info")
		return ""
	}
	file = filepath.ToSlash(file)
	file = filepath.Clean(file)
	return filepath.FromSlash(file)
}

// 项目编译后的执行文件 执行时 原命令所在的目录
// 可执行文件 所在的目录 ./路径 配置文件等 需要基于该目录
func ExecutedCurrentFilePath() (string, error) {
	fp, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fp, err
	}
	fp = filepath.ToSlash(fp)
	fp = filepath.Clean(fp)
	return filepath.FromSlash(fp), nil
}

// 可执行文件 执行时的当前目录 log 和数据库需要基于该目录  默认./就是基于该目录
func ExecutingCurrentFilePath() string {
	getwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	getwd = filepath.ToSlash(getwd)
	getwd = filepath.Clean(getwd)
	return filepath.FromSlash(getwd)
}

// 以执行文件放置的目录 返回补全 整理过的文件路径
func DealWithExecutedCurrentFilePath(fp string) (string, error) {
	isAbs := filepath.IsAbs(fp)
	// 如果是绝对路径
	if isAbs {
		// 替换路径中的'/'为路径分隔符
		fs := filepath.ToSlash(fp)
		fs = filepath.Clean(fs)
		// 返回处理过的路径名
		return filepath.FromSlash(fs), nil
	} else {
		// 如果是相对路径 以放置二进制文件的当前目录 为基准返回完整的路径名
		fs, err := ExecutedCurrentFilePath()
		if err != nil {
			return fp, err
		}
		fs = filepath.ToSlash(fs)
		fs = filepath.Clean(fs)
		fpd := filepath.ToSlash(fp)
		fpd = filepath.Clean(fpd)
		fs = filepath.Join(fs, fpd)
		return filepath.FromSlash(fs), nil
	}
}

// 以执行文件执行的目录 返回补全 整理过的文件路径
func DealWithExecutingCurrentFilePath(fp string) string {
	isAbs := filepath.IsAbs(fp)
	// 如果是绝对路径
	if isAbs {
		// 替换路径中的'/'为路径分隔符
		fs := filepath.ToSlash(fp)
		fs = filepath.Clean(fs)
		// 返回处理过的路径名
		return filepath.FromSlash(fs)
	} else {
		// 如果是相对路径 以放置二进制文件的当前目录 为基准返回完整的路径名
		fs := ExecutingCurrentFilePath()
		fs = filepath.ToSlash(fs)
		fs = filepath.Clean(fs)
		fpd := filepath.ToSlash(fp)
		fpd = filepath.Clean(fpd)
		fs = filepath.Join(fs, fpd)
		return filepath.FromSlash(fs)
	}
}

func PathJoin(rootPath, leafPath string) (filePath string) {
	rootPath = filepath.ToSlash(rootPath)
	rootPath = filepath.Clean(rootPath)
	leafPath = filepath.ToSlash(leafPath)
	leafPath = filepath.Clean(leafPath)
	filePath = filepath.Join(rootPath, leafPath)
	return filepath.FromSlash(filePath)
}

func PathExists(path string) (bool, error) {
	fOrPath, err := os.Stat(path)
	if err == nil {
		if fOrPath.IsDir() {
			return true, nil
		}
		return false, errors.New("exists same name file - 存在与目录同名的文件")
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func FileExists(path string) (bool, error) {
	fOrPath, err := os.Stat(path)
	if err == nil {
		if fOrPath.IsDir() {
			return false, errors.New("exists same name file path - 存在与文件同名的目录")
		}
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CreatePathDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return err
}
