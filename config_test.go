package main

import (
	"reflect"
	"testing"
)

func TestOutputConfig(t *testing.T) {
	target := Config{
		"change_to_your_own_osu_api_key",
		"path/to/osu!/folder",
		"osu!.db",
		"collection.db",
		"username",
		"password",
		"official",
	}
	target.output()
}

func TestInitConfig(t *testing.T) {
	target := Config{}
	target.init()
	exception := Config{
		"change_to_your_own_osu_api_key",
		"path/to/osu!/folder",
		"osu!.db",
		"collection.db",
		"username",
		"password",
		"official",
	}
	if !reflect.DeepEqual(target, exception) {
		t.Error("error in init config")
	}
}
