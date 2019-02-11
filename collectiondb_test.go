package main

import (
	"fmt"
	"path"
	"testing"
)

func TestParseFile(t *testing.T) {
	collList := new(CollectionList)
	collList.ParseFile("import/2016collection.db")
	fmt.Println(len(collList.list))
}

func TestOutputFile(t *testing.T) {
	osuRootPath := "d:/osu!"
	beatmapList := new(BeatmapList)
	beatmapList.ParseFile(path.Join(osuRootPath, "osu!.db"))
	H := beatmapList.getHashBidMap()
	collList := new(CollectionList)
	collList.ParseFile(path.Join(osuRootPath, "collection.db"))
	collList.outputFile(H)
}

func TestMergeColl(t *testing.T) {

}
