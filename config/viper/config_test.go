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
