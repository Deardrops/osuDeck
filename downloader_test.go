package main

import (
	"os"
	"reflect"
	"testing"
)

var tests = []struct {
	name string
	in   string
	out  error
}{
	{"ranked", "219380", nil},             // toumei elegy
	{"unranked", "156445", nil},           // Masked bitcH
	{"deleted", "228811", ErrUnavailable}, // kimashitawa
	{"DMCA", "90784", ErrUnavailable},     // anhedonia
}

func InitTestEnv() {
	conf = newConf()
	_ = os.RemoveAll(downloadFolderPath)
	_ = os.MkdirAll(downloadFolderPath, os.ModePerm)
}

func CheckResult(get error, want error, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(get, want) {
		t.Errorf("get: %T: %v\nwant: %T: %v\n", get, get, want, want)
	}
}

func TestOsuOfficialDownload(t *testing.T) {
	InitTestEnv()
	downloader := newDownloader("official")

	t.Run("ranked beatmap", func(t *testing.T) {
		tt := tests[0]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("unranked beatmap", func(t *testing.T) {
		tt := tests[1]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("deleted beatmap", func(t *testing.T) {
		tt := tests[2]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("DMCA beatmap", func(t *testing.T) {
		tt := tests[3]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
}

func TestSayobotDownload(t *testing.T) {
	InitTestEnv()
	downloader := newDownloader("sayobot")
	t.Run("ranked beatmap", func(t *testing.T) {
		tt := tests[0]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("unranked beatmap", func(t *testing.T) {
		tt := tests[1]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("deleted beatmap", func(t *testing.T) {
		tt := tests[2]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("DMCA beatmap", func(t *testing.T) {
		tt := tests[3]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
}

func TestBlootcatDownload(t *testing.T) {
	InitTestEnv()
	downloader := newDownloader("bloodcat")

	t.Run("ranked beatmap", func(t *testing.T) {
		tt := tests[0]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("unranked beatmap", func(t *testing.T) {
		tt := tests[1]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("deleted beatmap", func(t *testing.T) {
		tt := tests[2]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
	t.Run("DMCA beatmap", func(t *testing.T) {
		tt := tests[3]
		err := downloader.download(tt.in)
		CheckResult(err, tt.out, t)
	})
}
