package main

import (
	"os"
	"testing"
)

func TestQuerySetIdByMd5(t *testing.T) {
	_ = os.Remove("querycache_test.db")
	defer func() {
		_ = os.Remove("querycache_test.db")
	}()
	conf := newConf()

	var api Api
	if len(conf.OsuApiKey) == 0 {
		t.Error("osu_api_key unset")
	} else {
		api = newApi(conf.OsuApiKey, "querycache_test.db")
		defer api.close()
	}

	md5 := "df1b615c3588932f554ed314e1a04924"
	want := "660914"

	get1 := api.QuerySetIdByMd5(md5)
	if get1 != want {
		t.Errorf("error in query to api.\nget: %v\nwant: %v\n", get1, want)
	}
	get2 := api.QuerySetIdByMd5(md5)
	if get2 != want {
		t.Errorf("error in query to cache db.\nget: %v\nwant: %v\n", get2, want)
	}
}
