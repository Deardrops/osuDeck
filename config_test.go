package main

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Run("create new file", func(t *testing.T) {
		_ = os.Remove("conf_test.yaml")
		target := newConf("conf_test.yaml")
		defer func() {
			_ = os.Remove("conf_test.yaml")
		}()
		exception := Config{
			Mirror:   "bloodcat",
			FilePath: "conf_test.yaml",
		}
		if !reflect.DeepEqual(target, exception) {
			t.Errorf("error, \nget: %v \nwant:%v", target, exception)
		}
	})
	t.Run("load exist config file", func(t *testing.T) {
		target := newConf("example.conf.yaml")
		exception := Config{
			"path/to/osu!/folder",
			"path/to/osu!.db",
			"path/to/collection.db",
			"change_to_your_own_osu_api_key",
			"change_to_your_osu_username",
			"change_to_your_osu_password",
			"bloodcat",
			"example.conf.yaml",
		}
		if !reflect.DeepEqual(target, exception) {
			t.Errorf("error, \nget: %v \nwant:%v", target, exception)
		}
	})
}
