package src

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

var (
	handled           = make(map[string]bool)
	packageSourcesDir = ""
)

func Sync() {
	// 初始化 home dir
	initPackageDir()

	// 执行命令所在目录应当存在 package.toml 文件
	workdir, err := os.Getwd()
	if err != nil {
		throw("cannot get workdir, err=%v", err)
	}
	configFile := fmt.Sprintf("%s/%s", workdir, PackageFile)
	SyncByConfig(configFile)
}

// ~/.nature/package/sources
// ~/.nature/package/caches
func initPackageDir() {
	u, err := user.Current()
	if err != nil {
		throw("cannot get current user, err=%v", err)
	}
	homeDir := u.HomeDir
	if homeDir == "" {
		throw("cannot get home dir")
	}

	packageSourcesDir = fmt.Sprintf("%s/.nature/package/sources", homeDir)
	err = os.MkdirAll(packageSourcesDir, 0755)
	if err != nil {
		throw("cannot create dir=%v, err=%v", packageSourcesDir, err)
	}

	packageCachesDir := fmt.Sprintf("%s/.nature/package/caches", homeDir)
	err = os.MkdirAll(packageCachesDir, 0755)
	if err != nil {
		throw("cannot create dir=%v, err=%v", packageCachesDir, err)
	}
}

func SyncByConfig(configFile string) {
	// 已经处理过的 package 避免重复处理
	if _, ok := handled[configFile]; ok {
		return
	}

	// - 判断 package.toml 文件是否存在
	p, err := Parser(configFile)
	if err != nil {
		throw("file=%s parser err=%v", configFile, err)
	}

	// 循环便利下载依赖
	for _, dep := range p.Dependencies {
		if dep.Version == "" {
			throw("version cannot be empty")
		}

		depPath := ""
		if dep.Type == DependencyTypeGit {
			depPath = SyncGit(dep.Url, dep.Version)
		} else {
			if dep.Path == "" {
				throw("path cannot be empty")
			}
			depPath = dep.Path
			log("use local dst=%v", depPath)
		}

		// 配置文件, 递归
		configFile := fmt.Sprintf("%s/%s", depPath, PackageFile)
		SyncByConfig(configFile)
	}

	handled[configFile] = true
}

func SyncGit(url, version string) string {
	if url == "" || version == "" {
		throw("url or version cannot be empty")
	}

	// - url 合法检测(不能携带 http 或者 https 前缀)
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		throw("url=%v cannot start with http or https", url)
	}

	dst := dstDir(url, version)
	if !dirExists(dst) {
		err := gitPull(url, version, dst)
		if err != nil {
			throw("git pull failed %v", err)
		}
		log("sync git success dst=%v", dst)
	} else {
		log("sync git exists dst=%v ", dst)
	}

	return dst
}
