package main

import (
	"log"
	"os"
	"path"
	"time"
)

var conf Config
var api Api

const logFolderPath = "./logs"

func main() {
	checkAndCreateFolder()
	conf = Config{}
	conf.init()

	api = Api{}
	api.init()

	now := time.Now().Format("20060102_150405")
	f, err := os.OpenFile(path.Join(logFolderPath, now+".txt"), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	defer conf.output()
	buildGui()
}

func checkAndCreateFolder() {
	list := []string{
		downloadFolderPath,
		logFolderPath,
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
