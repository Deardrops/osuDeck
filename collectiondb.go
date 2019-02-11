package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Deardrops/binpacker"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Collection struct {
	name string
	md5s Set
}

type CollectionList struct {
	version int32
	list    []Collection
}

func (CL *CollectionList) ParseBytes(buffer io.Reader) {
	r := binpacker.NewUnpacker(binary.LittleEndian, buffer)

	r.FetchInt32(&CL.version)
	var count int32
	r.FetchInt32(&count)

	for i := 0; i < int(count); i++ {
		coll := Collection{}
		r.FetchVarString(&coll.name)
		var mapCount int32
		r.FetchInt32(&mapCount)
		coll.md5s = make(Set)
		for j := 0; j < int(mapCount); j++ {
			var checksum string
			r.FetchVarString(&checksum)
			coll.md5s[checksum] = true
		}
		CL.list = append(CL.list, coll)
	}
}

func (CL *CollectionList) ParseFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		panic(err)
	}
	buffer := bufio.NewReaderSize(f, int(fInfo.Size()))
	CL.ParseBytes(buffer)
}

func (CL *CollectionList) saveToDisk(path string) {
	buffer := new(bytes.Buffer)
	p := binpacker.NewPacker(binary.LittleEndian, buffer)
	p.PushInt32(CL.version)
	count := len(CL.list)
	p.PushInt32(int32(count))

	for _, coll := range CL.list {
		p.PushVarString(coll.name)
		mapCount := len(coll.md5s)
		p.PushInt32(int32(mapCount))
		for md5 := range coll.md5s {
			p.PushVarString(md5)
		}
	}
	f, err := os.Create(path)
	if os.IsExist(err) {
		panic(err)
	}
	defer f.Close()
	f.Write(buffer.Bytes())
}

type SimpleCollectionList struct {
	Name       string `json:"name"`
	ExistCount int    `json:"existCount"`
	TotalCount int    `json:"totalCount"`
}

func (CL *CollectionList) output(localMapMd5s Set) string {
	var res []SimpleCollectionList
	for _, coll := range CL.list {
		var existCount int
		for md5 := range coll.md5s {
			if stringInSet(md5, localMapMd5s) {
				existCount++
			}
		}
		res = append(res, SimpleCollectionList{
			Name:       coll.name,
			ExistCount: existCount,
			TotalCount: len(coll.md5s),
		})
	}
	output, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(output)
}

func (CL *CollectionList) outputFile(H HashBidMap) {
	for _, coll := range CL.list {
		var res []string
		if strings.HasPrefix(coll.name, "!") {
			continue
		}
		for md5 := range coll.md5s {
			if bid, ok := H[md5]; ok {
				res = append(res, strconv.Itoa(int(bid)))
			}
		}
		if err := saveToFile("export/"+coll.name+".txt", res); err != nil {
			fmt.Println(err)
		}
	}
}

func (CL *CollectionList) getMapInCollectionTotal() int {
	var total int
	for _, coll := range CL.list {
		total += len(coll.md5s)
	}
	return total
}

func (CL *CollectionList) getCollectionTotal() int {
	return len(CL.list)
}

func (CL *CollectionList) getMissedTotal(localMd5s *Set) int {
	return len(CL.getMissedMd5s(localMd5s))
}

func (CL *CollectionList) getMissedMd5s(localMd5s *Set) []string {
	var res []string
	for _, coll := range CL.list {
		for md5 := range coll.md5s {
			if _, ok := (*localMd5s)[md5]; !ok {
				res = append(res, md5)
			}
		}
	}
	return res
}

func (CL *CollectionList) merge(importedCL CollectionList, newFileName string) {
	var addedCollNames []string
	for _, coll := range CL.list {
		addedCollNames = append(addedCollNames, strings.TrimSpace(coll.name))
		if !strings.HasPrefix(coll.name, " ") {
			coll.name = " " + coll.name
		}
	}
	log.Printf("scaned %d collections in CL db file.", len(addedCollNames))
	var missedMapCount, addedCollCount, addedMapCount int
	for _, coll := range importedCL.list {
		if stringInArray(coll.name, addedCollNames) {
			fmt.Printf("failed to add colletion: Name duplicated - %s.", coll.name)
			continue
		}
		for md5 := range coll.md5s {
			if !stringInSet(md5, localBeatmapMd5s) {
				delete(coll.md5s, md5)
				missedMapCount++
			} else {
				addedMapCount++
			}
		}
		CL.list = append(CL.list, coll)
		addedCollCount++
	}
	log.Printf("created %d collections with %d beatmaps\n", addedCollCount, addedMapCount)
	log.Printf("missed beatmap: %d\n", missedMapCount)
	CL.saveToDisk(newFileName)
}
