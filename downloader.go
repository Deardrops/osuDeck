package main

import (
	"crypto/md5"
	"fmt"
	"github.com/cavaliercoder/grab"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	ErrUnavailable = DownloaderErr("beatmapset download url unavailable")
)

type DownloaderErr string

func (e DownloaderErr) Error() string {
	return string(e)
}

type Downloader interface {
	download(string) error
	Init()
}

var downloadFolderPath = "./download"

type OsuOfficialClient struct {
	baseUrl string
}

func (C *OsuOfficialClient) download(setId string) error {
	baseUrl := "https://osu.ppy.sh/d/%sn?u=%s&h=%s"
	password := fmt.Sprintf("%x", md5.Sum([]byte(conf.Password)))
	downloadUrl := fmt.Sprintf(baseUrl, setId, conf.Username, password)
	resp, err := grab.Get(downloadFolderPath, downloadUrl)
	if err != nil {
		log.Printf("%T: %v\n", err, err)
		return err
	}

	if resp.Size == -1 {
		_ = os.Remove(path.Join(downloadFolderPath, setId+"n"))
		return ErrUnavailable
	}
	log.Printf("finished %s\n", path.Base(resp.Filename))
	return nil
}

func (C *OsuOfficialClient) Init() {
	C.baseUrl = "https://osu.ppy.sh/d/%s"
	if len(conf.Username) == 0 || len(conf.Password) == 0 {
		return
	}
	password := fmt.Sprintf("%x", md5.Sum([]byte(conf.Password)))
	C.baseUrl = C.baseUrl + fmt.Sprintf("n?u=%s&h=%s", conf.Username, password)
}

type SayobotClient struct{}

func (C *SayobotClient) download(setId string) error {
	baseUrl := "https://txy1.sayobot.cn/download/osz/novideo/%s"
	downloadUrl := strings.Replace(baseUrl, "%s", setId, 1)
	resp, err := grab.Get(downloadFolderPath, downloadUrl)
	if err != nil {
		log.Printf("%T: %v\n", err, err)
		return err
	}
	log.Printf("finished %s \n", path.Base(resp.Filename))
	return nil
}

func (C *SayobotClient) Init() {

}

type BloodcatClient struct {
	*grab.Client
}

func (c *BloodcatClient) Init() {
	c.Client = grab.NewClient()
	// disable Storyboard and background video
	cookie := []*http.Cookie{
		{
			Name:  "DLOPT",
			Value: "%7B%22bg%22%3Afalse%2C%22video%22%3Atrue%2C%22skin%22%3Atrue%2C%22cdn%22%3Afalse%7D",
		},
	}
	var bloodcatHost = "https://bloodcat.com/osu"
	var err error
	u, err := url.Parse(bloodcatHost)
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("error in init BloodcatClient")
	}
	jar.SetCookies(u, cookie)
	c.Client.HTTPClient.Jar = jar
}

func (C *BloodcatClient) download(setId string) error {
	baseUrl := "https://bloodcat.com/osu/s/%s"
	downloadUrl := strings.Replace(baseUrl, "%s", setId, 1)

	req, err := grab.NewRequest(downloadFolderPath, downloadUrl)
	resp := C.Client.Do(req)
	if resp.Err() != nil {
		if grab.IsStatusCodeError(err) {
			log.Printf("expected IsStatusCodeError to return true for %T: %v\n", err, err)
		} else {
			log.Printf("%T: %v\n", err, err)
		}
		return err
	}
	log.Printf("finished %s \n", path.Base(resp.Filename))
	return nil
}

func newDownloader(name string) Downloader {
	var downloader Downloader
	switch name {
	case "official":
		downloader = new(OsuOfficialClient)
	case "bloodcat":
		downloader = new(BloodcatClient)
	case "sayobot":
		downloader = new(SayobotClient)
	default:
		log.Fatalf("error in create downloader, wrong name %s\n", name)
	}
	downloader.Init()
	return downloader
}
