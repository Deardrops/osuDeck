package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var conf Config
var api Api

var root string

const logFolder = "logs"
const downloadFolder = "download"

func main() {
	root, _ = filepath.Abs("./")
	checkAndCreateFolder()
	conf = newConf()
	defer conf.output()

	if len(conf.OsuApiKey) != 0 {
		api = newApi(conf.OsuApiKey)
		defer api.close()
	} else {
		log.Printf("missed osu_api_key value in conf.yaml.\n")
	}

	now := time.Now().Format("20060102_150405")
	f, err := os.OpenFile(path.Join(root, logFolder, now+".txt"), os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	buildGui()
}

func checkAndCreateFolder() {
	list := []string{
		path.Join(root, downloadFolder),
		path.Join(root, logFolder),
	}
	for _, folder := range list {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			err = os.Mkdir(folder, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// TODO：谱面上传脚本
// https://bloodcat.com/osu/custom.php?mod=add&id=453990&m=s
// sayobot 对于未缓存谱面会触发 http2 的某个 bug
