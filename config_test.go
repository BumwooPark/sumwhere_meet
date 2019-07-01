package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)
import "github.com/spf13/viper"

type Config struct {
	Database struct {
		Httpport string
		Driver string
	}
}

func Test_Config(t *testing.T){
	var config Config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	assert.NoError(t,viper.ReadInConfig())

	//f, err := os.Open(filepath.Join(".", "config.yml"))
	//if err != nil {
	//	fmt.Errorf("Fatal error config file: %s \n", err)
	//}
	//defer f.Close()
	//viper.MergeConfig(f)

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Errorf("Fatal error config file: %s \n", err)
	}

	fmt.Println(config)

}