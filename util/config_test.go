package util

import "testing"

type config struct {
	App *appConfig
}

type appConfig struct {
	Name string
	Mode string
}

func TestInitConfig(t *testing.T) {
	conf := new(config)
	if err := InitConfig(conf, "./config.yml"); err != nil {
		t.Fatal(err)
	}
	if conf.App.Name != "go-common" || conf.App.Mode != "release" {
		t.Fatal("conf read failed")
	}
}
