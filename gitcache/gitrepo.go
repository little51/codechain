package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/robfig/cron"
)

var _PATH_DEPTH = 2
var _IS_SYNC = false

func fetchMirrorFromRemoteUnshallow(repository string) {
	remote := "https:/" + strings.Replace(repository, g_Basedir, "", -1)
	local := repository
	fmt.Printf("%s %s\n", remote, local)
	fetchMirrorFromRemote(remote, local, "unshallow")
}

func walkDir(dirpath string, depth int, f func(string)) {
	if depth > _PATH_DEPTH {
		return
	}
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			walkDir(dirpath+"/"+file.Name(), depth+1, f)
			if strings.HasSuffix(file.Name(), ".git") {
				f(dirpath + "/" + file.Name())
			}
			continue
		}
	}
}

func SyncLocalMirrorFromRemote() {
	if _IS_SYNC {
		return
	}
	log.Println("SyncLocalMirrorFromRemote")
	_IS_SYNC = true
	walkDir(g_Basedir, 0, fetchMirrorFromRemoteUnshallow)
	_IS_SYNC = false
}

func GetLocalMirrorsInfo() string {
	return ""
}

func GetMirrorProgress(repoName string) string {
	return ""
}

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 */2 * * *", func() {
		SyncLocalMirrorFromRemote()
	})
	c.Start()
	log.Println("cron start")
	return
}
