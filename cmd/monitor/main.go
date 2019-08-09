package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jenska/go2m68k/pkg/cpu"
	"github.com/jenska/go2m68k/pkg/mem"
	"github.com/jenska/go2m68k/pkg/util"
	"github.com/spf13/viper"
)

type Image struct {
	Description string
	Language    string
	Size        string
	Path        string
}
type Configuration struct {
	Version      string
	DefaultImage int
	Images       []Image
}

var configuration Configuration

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file %s", err))
	}

	viper.ReadInConfig()
	if viper.Unmarshal(&configuration) != nil {
		panic(fmt.Errorf("unable to decode configuration %s", err))
	}

	path := configuration.Images[configuration.DefaultImage].Path
	size := configuration.Images[configuration.DefaultImage].Size

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to load image %s", err))
	}

	startROM := cpu.Address(0xe00000)
	if size == "192k" {
		startROM = cpu.Address(0xfc0000)
	}

	bus := mem.NewAddressBus(
		mem.NewProtectedRAM(0, 2048),
		mem.NewRAM(2048, 1024*1024),
		mem.NewROM(startROM, data),
	)

	cpu.NewCPU(bus)
	util.Dump(bus, startROM, 128)
}
