package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	fileDialog "github.com/sqweek/dialog"
	"log"
	"path"
)

func selectFileDialog(name string) string {
	filename, err := fileDialog.File().Filter(name+" db file", "db").Load()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return filename
}

func selectFolderDialog() string {
	directory, err := fileDialog.Directory().Title("osu! folder").Browse()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return directory
}

type Label struct {
	text  string
	value int
}

type EntryGroup struct {
	list    []Label
	entries []*widget.Label
	group   *widget.Group
}

func (o *Label) String() string {
	return fmt.Sprintf(o.text, o.value)
}

func (o *EntryGroup) createEntries(name string) *widget.Group {
	o.group = widget.NewGroup(name)
	for _, object := range o.list {
		w := widget.NewLabel(object.String())
		o.entries = append(o.entries, w)
		o.group.Append(w)
	}
	o.group.Hide()
	return o.group
}

func (o *EntryGroup) updateEntries() {
	for i, object := range o.list {
		o.entries[i].SetText(object.String())
	}
}

func (o *EntryGroup) setValues(values ...int) {
	if len(values) != len(o.list) {
		log.Fatalf("error when Set values for entries, get %d, wanted %d", len(values), len(o.list))
	}
	for i := range o.list {
		o.list[i].value = values[i]
	}
	o.updateEntries()
}

var localBeatmapMd5s Set

var localBeatmapsetIds BeatmapsetIds

var downloadSetIds Set

func buildGui() {
	a := app.New()
	w := a.NewWindow("osu!Deck")
	localEntries := EntryGroup{
		list: []Label{
			{"beatmapset: %d", 0},
			{"beatmap: %d", 0},
			{"colletion: %d", 0},
			{"map in collection: %d", 0},
		},
	}
	osuDbEntries := EntryGroup{
		list: []Label{
			{"beatmap: %d", 0},
			{"beatmapset: %d", 0},
			{"missed beatmapset: %d", 0},
		},
	}
	collectionDbEntries := EntryGroup{
		list: []Label{
			{"collection: %d", 0},
			{"beatmap: %d", 0},
			{"missed beatmap: %d", 0},
		},
	}

	downloader := newDownloader(conf.Mirror)

	wipInfo := NewWipInfo(w)
	progress := widget.NewProgressBar()
	w.SetContent(widget.NewVBox(
		widget.NewHBox(
			widget.NewLabel("Open: "),
			widget.NewButton("osu! folder", func() {
				conf.LocalOsuRootPath = selectFolderDialog()
				ParseAndRenderLocalOsu(conf.LocalOsuRootPath, &localEntries)
			}),
			widget.NewButton("osu!.db", func() {
				conf.OsuDbPath = selectFileDialog("osu!")
				ParseAndRenderOsuDbFile(conf.OsuDbPath, &osuDbEntries)
			}),
			widget.NewButton("collection.db", func() {
				conf.CollectionDbPath = selectFileDialog("colletion")
				ParseAndRenderCollectionDbFile(conf.CollectionDbPath, &collectionDbEntries)
			}),
		),
		widget.NewTabContainer(
			widget.NewTabItem("Local",
				widget.NewVBox(
					localEntries.createEntries("Local"),
					widget.NewGroup("You can:",
						widget.NewButton("output all collection with beatmapsetId", func() {
							wipInfo.Show()
						}),
					),
				),
			),
			widget.NewTabItem("osu!db",
				widget.NewVBox(
					osuDbEntries.createEntries("osu!db"),
					widget.NewGroup("You can:",
						widget.NewButton("download missed beatmaps", func() {
							if len(localBeatmapMd5s) == 0 {
								dialog.NewInformation("Warning", "Please specified your local osu! folder firstly", w)
								log.Println("local osu! folder unspecified")
								return
							}
							missedSetId := importedBeatmapList.getMissedMapsetId(&localBeatmapsetIds)
							defer progress.Hide()
							for i, setId := range missedSetId {
								log.Println("start download beatmapset", setId)
								err := downloader.download(setId)
								if err != nil {
									log.Println(err)
								}
								percent := float64(i+1) / float64(len(missedSetId))
								progress.SetValue(percent)
							}
						}),
					),
				),
			),
			widget.NewTabItem("collection",
				widget.NewVBox(
					collectionDbEntries.createEntries("collection.db"),
					widget.NewGroup("You can",
						widget.NewButton("load beatmaps in collection", func() {
							if len(conf.OsuApiKey) == 0 {
								dialog.NewInformation("Warning", "Please set your osu_api_key in conf.yaml file firstly", w)
								log.Println("require osu_api_key unspecified")
								return
							}
							if len(localBeatmapMd5s) == 0 {
								dialog.NewInformation("Warning", "Please specified your local osu! folder firstly", w)
								log.Println("local osu! folder unspecified")
								return
							}
							downloadSetIds = make(Set)
							list := importedCollectionList.getMissedMd5s(&localBeatmapMd5s)
							for i, md5 := range list {
								setId := api.QuerySetIdByMd5(md5)
								downloadSetIds[setId] = true
								percent := float64(i+1) / float64(len(list))
								progress.SetValue(percent)
							}
						}),
						widget.NewButton("download missed beatmaps", func() {
							if len(downloadSetIds) == 0 {
								dialog.NewInformation("Warning", "Please load beatmaps firstly", w)
								return
							}
							defer progress.Hide()
							var i int
							for setId := range downloadSetIds {
								log.Println("start download beatmapset", setId)
								err := downloader.download(setId)
								if err != nil {
									log.Println(err)
								}
								i++
								percent := float64(i+1) / float64(len(downloadSetIds))
								progress.SetValue(percent)
							}
						}),
						widget.NewButton("merge into local collection", func() {
							wipInfo.Show()
						}),
						widget.NewButton("output all collection with beatmapid ", func() {
							wipInfo.Show()
						}),
					),
				),
			),
		),
		progress,
	))
	ParseAndRenderLocalOsu(conf.LocalOsuRootPath, &localEntries)
	ParseAndRenderOsuDbFile(conf.OsuDbPath, &osuDbEntries)
	ParseAndRenderCollectionDbFile(conf.CollectionDbPath, &collectionDbEntries)
	localEntries.group.Show()
	w.ShowAndRun()
}

var localBeatmapList, importedBeatmapList BeatmapList
var localCollectionList, importedCollectionList CollectionList

func ParseAndRenderLocalOsu(osuRootPath string, entries *EntryGroup) {
	if len(osuRootPath) == 0 {
		return
	}
	localBeatmapList = BeatmapList{}
	localBeatmapList.ParseFile(path.Join(osuRootPath, "osu!.db"))
	localBeatmapMd5s = localBeatmapList.getMd5Sets()
	localBeatmapsetIds = localBeatmapList.getBeatmapsetIdSets()
	localCollectionList = CollectionList{}
	localCollectionList.ParseFile(path.Join(osuRootPath, "collection.db"))
	entries.setValues(
		localBeatmapList.getBeatmapsetTotal(),
		localBeatmapList.getBeatmapTotal(),
		localCollectionList.getCollectionTotal(),
		localCollectionList.getMapInCollectionTotal(),
	)
}

func ParseAndRenderOsuDbFile(osuDbFilePath string, entries *EntryGroup) {
	if len(osuDbFilePath) == 0 {
		return
	}
	importedBeatmapList = BeatmapList{}
	importedBeatmapList.ParseFile(osuDbFilePath)
	entries.setValues(
		importedBeatmapList.getBeatmapsetTotal(),
		importedBeatmapList.getBeatmapTotal(),
		importedBeatmapList.getMissedTotal(&localBeatmapsetIds),
	)
}

func ParseAndRenderCollectionDbFile(collectionDbFilePath string, entries *EntryGroup) {
	if len(collectionDbFilePath) == 0 {
		return
	}
	importedCollectionList = CollectionList{}
	importedCollectionList.ParseFile(collectionDbFilePath)
	entries.setValues(
		importedCollectionList.getCollectionTotal(),
		importedCollectionList.getMapInCollectionTotal(),
		importedCollectionList.getMissedTotal(&localBeatmapMd5s),
	)
}

func NewWipInfo(w fyne.Window) dialog.Dialog {
	d := dialog.NewInformation("WIP", "This function is working in progress", w)
	return d
}
