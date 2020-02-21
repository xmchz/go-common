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
	} else {
		t.Logf("config: %s, %s", conf.App.Name, conf.App.Mode)
	}

}
