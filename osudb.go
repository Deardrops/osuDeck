package main

import (
	"bufio"
	"encoding/binary"
	"github.com/Deardrops/binpacker"
	"log"
	"os"
	"strconv"
)

type Beatmap struct {
	size                     int32
	artist                   string
	artistUnicode            string
	title                    string
	titleUnicode             string
	creator                  string
	version                  string
	audioFileName            string
	beatmapChecksum          string
	filename                 string
	submissionStatus         byte
	countNormal              uint16
	countSlider              uint16
	countSpinner             uint16
	dataModified             int64
	ar                       float32
	cs                       float32
	hp                       float32
	od                       float32
	sliderMultiplier         float64
	tomStars                 []map[int32]float64
	drainLength              int32
	totalLength              int32
	previewTime              int32
	controlPoints            []ControlPoint
	beatmapId                int32
	beatmapSetId             int32
	beatmapTopicId           int32
	playerRankOsu            byte
	playerRankFruits         byte
	playerRankTaiko          byte
	playerRankMania          byte
	playerOffset             int16
	stackLeniency            float32
	playMode                 byte
	source                   string
	tags                     string
	onlineOffset             int16
	onlineDisplayTitle       string
	newFile                  bool
	dataLastPlayed           int64
	inOszContainer           bool
	containingFolderAbsolute string
	lastInfoUpdate           int64
	disableSamples           bool
	disableSkins             bool
	disableStoryboard        bool
	disableVideo             bool
	visualSettingsOverride   bool
	lastEditTime             int32
	maniaSpeed               byte
}

type BeatmapList struct {
	version         int32
	username        string
	folderFileCount int32
	accountUnlocked bool
	unlockDate      int64
	list            []Beatmap
}

type ControlPoint struct {
	beatLength   float64
	offset       float64
	timingChange bool
}

type Set map[string]bool

type HashBidMap map[string]int32

type BeatmapsetIds map[int32]bool

func (BL *BeatmapList) getHashBidMap() HashBidMap {
	res := make(HashBidMap)
	for _, b := range BL.list {
		res[b.beatmapChecksum] = b.beatmapId
	}
	return res
}

func (BL *BeatmapList) getMd5Sets() Set {
	localMapMd5s := make(Set)
	for _, b := range BL.list {
		localMapMd5s[b.beatmapChecksum] = true
	}
	return localMapMd5s
}

func (BL *BeatmapList) getBeatmapsetIdSets() BeatmapsetIds {
	beatmapsetIdSets := make(BeatmapsetIds)
	for _, b := range BL.list {
		beatmapsetIdSets[b.beatmapSetId] = true
	}
	return beatmapsetIdSets
}

func (BL *BeatmapList) downloadMissedMapsets(file string) {
	mapsetIds, err := readFromFile(file)
	if err != nil {
		log.Println("error in read file ", file)
	}
	localMapMd5s := BL.getMd5Sets()
	var downloadList []string
	for _, setId := range mapsetIds {
		if _, ok := localMapMd5s[setId]; ok {
			continue
		}
		downloadList = append(downloadList, setId)
	}
	downloader := newDownloader("bloodcat")
	for i, id := range downloadList {
		err := downloader.download(id)
		if err != nil {
			log.Printf("download %s failed, %s", id, err)
		}
		percent := i / len(mapsetIds) * 100
		if percent%5 == 0 {
			log.Printf("Current Progress: %d \n", percent)
		}
	}
}

func (BL *BeatmapList) getBeatmapsetTotal() int {
	beatmapsetList := make(map[int32]bool)
	total := 0
	for _, beatmap := range BL.list {
		if beatmap.beatmapSetId == 0 {
			total++
			continue
		}
		if _, ok := beatmapsetList[beatmap.beatmapSetId]; !ok {
			beatmapsetList[beatmap.beatmapSetId] = true
			total++
		}
	}
	return total
}

func (BL *BeatmapList) getBeatmapTotal() int {
	return len(BL.list)
}

func (BL *BeatmapList) getMissedMapsetId(localBeatmapsetIds *BeatmapsetIds) []string {
	setIdMaps := Set{}
	for _, beatmap := range BL.list {
		if _, ok := (*localBeatmapsetIds)[beatmap.beatmapSetId]; ok {
			continue
		}
		if beatmap.beatmapSetId == 0 {
			continue
		}
		setId := strconv.Itoa(int(beatmap.beatmapSetId))
		if _, ok := setIdMaps[setId]; ok {
			continue
		}
		setIdMaps[setId] = true
	}
	var res []string
	for setId := range setIdMaps {
		res = append(res, setId)
	}
	return res
}

func (BL *BeatmapList) getMissedTotal(localBeatmapsetIds *BeatmapsetIds) int {
	return len(BL.getMissedMapsetId(localBeatmapsetIds))
}

func (BL *BeatmapList) ParseFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fInfo, err := f.Stat()
	if err != nil {
		panic(err)
	}
	buffer := bufio.NewReaderSize(f, int(fInfo.Size()))
	r := binpacker.NewUnpacker(binary.LittleEndian, buffer)

	r.FetchInt32(&BL.version)
	r.FetchInt32(&BL.folderFileCount)

	r.FetchBool(&BL.accountUnlocked)
	r.FetchInt64(&BL.unlockDate)

	r.FetchVarString(&BL.username)

	var count int32
	r.FetchInt32(&count)

	for i := 0; i < int(count); i++ {
		b := Beatmap{}
		r.FetchInt32(&b.size)
		r.FetchVarString(&b.artist)
		r.FetchVarString(&b.artistUnicode)
		r.FetchVarString(&b.title)
		r.FetchVarString(&b.titleUnicode)
		r.FetchVarString(&b.creator)
		r.FetchVarString(&b.version)
		r.FetchVarString(&b.audioFileName)
		r.FetchVarString(&b.beatmapChecksum)
		r.FetchVarString(&b.filename)
		r.FetchByte(&b.submissionStatus)
		r.FetchUint16(&b.countNormal)
		r.FetchUint16(&b.countSlider)
		r.FetchUint16(&b.countSpinner)
		r.FetchInt64(&b.dataModified)

		r.FetchFloat32(&b.ar)
		r.FetchFloat32(&b.cs)
		r.FetchFloat32(&b.hp)
		r.FetchFloat32(&b.od)

		r.FetchFloat64(&b.sliderMultiplier)

		for i := 0; i < 4; i++ {
			pair := make(map[int32]float64)
			r.FetchIntDoublePair(&pair)
			b.tomStars = append(b.tomStars, pair)
		}

		r.FetchInt32(&b.drainLength)
		r.FetchInt32(&b.totalLength)
		r.FetchInt32(&b.previewTime)

		var c int32
		r.FetchInt32(&c)
		for i := 0; i < int(c); i++ {
			controlPoint := ControlPoint{}
			r.FetchFloat64(&controlPoint.beatLength)
			r.FetchFloat64(&controlPoint.offset)
			r.FetchBool(&controlPoint.timingChange)
			b.controlPoints = append(b.controlPoints, controlPoint)
		}

		r.FetchInt32(&b.beatmapId)
		r.FetchInt32(&b.beatmapSetId)
		r.FetchInt32(&b.beatmapTopicId)
		r.FetchByte(&b.playerRankOsu)
		r.FetchByte(&b.playerRankFruits)
		r.FetchByte(&b.playerRankTaiko)
		r.FetchByte(&b.playerRankMania)
		r.FetchInt16(&b.playerOffset)
		r.FetchFloat32(&b.stackLeniency)
		r.FetchByte(&b.playMode)
		r.FetchVarString(&b.source)
		r.FetchVarString(&b.tags)
		r.FetchInt16(&b.onlineOffset)
		r.FetchVarString(&b.onlineDisplayTitle)
		r.FetchBool(&b.newFile)
		r.FetchInt64(&b.dataLastPlayed)
		r.FetchBool(&b.inOszContainer)
		r.FetchVarString(&b.containingFolderAbsolute)
		r.FetchInt64(&b.lastInfoUpdate)
		r.FetchBool(&b.disableSamples)
		r.FetchBool(&b.disableSkins)
		r.FetchBool(&b.disableStoryboard)
		r.FetchBool(&b.disableVideo)
		r.FetchBool(&b.visualSettingsOverride)

		r.FetchInt32(&b.lastEditTime)
		r.FetchByte(&b.maniaSpeed)

		BL.list = append(BL.list, b)
	}
}
