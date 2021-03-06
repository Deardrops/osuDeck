package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/parnurzeal/gorequest"
	"log"
	"time"
)

type apiBeatmap struct {
	gorm.Model
	Approved        string `json:"approved"`
	ApprovedDate    string `json:"approved_date"`
	LastUpdate      string `json:"last_update"`
	Artist          string `json:"artist"`
	BeatmapId       string `json:"beatmap_id"`
	BeatmapsetId    string `json:"beatmapset_id"`
	BPM             string `json:"bpm"`
	Creator         string `json:"creator"`
	CreatorId       string `json:"creator_id"`
	Diffcultyrating string `json:"diffcultyrating"`
	DiffSize        string `json:"diff_size"`
	DiffOverall     string `json:"diff_overall"`
	DiffApproach    string `json:"diff_approach"`
	DiffDrain       string `json:"diff_drain"`
	HitLength       string `json:"hit_length"`
	Source          string `json:"source"`
	GenreId         string `json:"genre_id"`
	LanguageId      string `json:"language_id"`
	Title           string `json:"title"`
	TotalLength     string `json:"total_length"`
	Version         string `json:"version"`
	FileMd5         string `json:"file_md5"`
	Mode            string `json:"mode"`
	Tags            string `json:"tags"`
	FavouriteCount  string `json:"favourite_count"`
	PlayCount       string `json:"playcount"`
	PassCount       string `json:"passcount"`
	MaxCombo        string `json:"max_combo"`
}

type Api struct {
	DB      *gorm.DB
	BaseUrl string
}

func (api *Api) close() {
	_ = api.DB.Close()
}

func (api *Api) QuerySetIdByMd5(md5 string) string {
	var beatmap apiBeatmap
	err := api.DB.Where("file_md5 = ?", md5).First(&beatmap).Error
	if err != nil {
		log.Println(err)
	}
	if beatmap.BeatmapId != "" {
		return beatmap.BeatmapsetId
	}
	if len(api.BaseUrl) == 0 {
		log.Println("please set your osu_api_key in conf.yaml file")
		return "0"
	}
	url := api.BaseUrl + md5
	_, body, _ := gorequest.New().Get(url).End()
	var beatmaps []apiBeatmap
	err = json.Unmarshal([]byte(body), &beatmaps)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(time.Second)
	if len(beatmaps) > 0 {
		api.DB.Create(&beatmaps[0])
		return beatmaps[0].BeatmapsetId
	} else {
		return "0"
	}
}

func newApi(osuApiKey string, paths ...string) Api {
	api := Api{}
	var err error
	if len(paths) == 0 {
		api.DB, err = gorm.Open("sqlite3", "querycache.db")
	} else {
		api.DB, err = gorm.Open("sqlite3", paths[0])
	}
	if err != nil {
		log.Fatalf("failed to connect database. %T:%s\n", err, err)
	}
	api.DB.AutoMigrate(&apiBeatmap{})

	baseUrl := "https://osu.ppy.sh/api/get_beatmaps?k=%s&h="
	api.BaseUrl = fmt.Sprintf(baseUrl, osuApiKey)
	return api
}

//TODO: 初始化时检测 osu_api_key 的合法性
