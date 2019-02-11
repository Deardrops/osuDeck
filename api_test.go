package main

import (
	"log"
	"os"
	"testing"
)

func TestQuerySetIdByMd5(t *testing.T) {
	err := os.Remove("test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove("test.db")
		if err != nil {
			log.Fatal(err)
		}
	}()
	md5 := "df1b615c3588932f554ed314e1a04924"
	ID := "660914"
	api := Api{}
	api.init("test.db")
	defer api.destruct()
	ID1 := api.QuerySetIdByMd5(md5)
	if ID1 != ID {
		t.Error("error in query to api")
	}
	ID2 := api.QuerySetIdByMd5(md5)
	if ID2 != ID {
		t.Error("error in query to cache db")
	}
}
