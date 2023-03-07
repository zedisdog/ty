package viper

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestNormal(t *testing.T) {
	viper.SetConfigType("yml")
	viper.SetConfigFile("./config.yaml")
	viper.ReadInConfig()
	fmt.Printf("%#v", viper.AllSettings())
}

func TestSet(t *testing.T) {
	v := viper.New()
	v.MergeConfigMap(map[string]interface{}{
		"a": "b",
		"c": map[string]interface{}{
			"d": "e",
		},
	})
	//fmt.Printf("%#v\n", v.AllSettings())
	println(v.GetString("c.d"))
}
