package main

import (
	"fmt"
	"os"

	"github.com/hsn/tiny-redis/config"
	"github.com/hsn/tiny-redis/logger"
	"github.com/hsn/tiny-redis/memdb"
	"github.com/hsn/tiny-redis/server"
)

func init() {
	memdb.RegisterKeyCommand()
	memdb.RegisterStringCommands()
	memdb.RegisterHashCommands()
}
func main() {
	cfg, err := config.Setup()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = logger.SetUp(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = server.Start(cfg)
	if err != nil {
		os.Exit(1)
	}
}
