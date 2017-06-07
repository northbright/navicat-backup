package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/northbright/pathhelper"
)

type Config struct {
	ServerAddr string `json:"server_addr"`
}

var (
	serverRoot = ""
	configFile = ""
	config     Config
	uploadDir  = ""
)

func init() {
	serverRoot, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(serverRoot, "config.json")
	uploadDir = path.Join(serverRoot, "uploaded")
}

func loadConfig(f string) (Config, error) {
	var buf []byte
	var err error

	c := Config{}

	// Load Conifg
	if buf, err = ioutil.ReadFile(f); err != nil {
		return c, err
	}

	if err = json.Unmarshal(buf, &c); err != nil {
		return c, err
	}

	return c, err
}

func uploadFile(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		err = fmt.Errorf("FormFile() error: %v", err)
		return
	}

	fileName := path.Join(uploadDir, header.Filename)

	out, err := os.Create(fileName)
	if err != nil {
		err = fmt.Errorf("os.Create() error: %v", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		err = fmt.Errorf("io.Copy() error: %v", err)
		return
	}
}

func main() {
	var err error

	defer func() {
		if err != nil {
			log.Printf("main() err: %v\n", err)
		}
	}()

	if config, err = loadConfig(configFile); err != nil {
		err = fmt.Errorf("loadConfig() error: %v", err)
		return
	}

	if err = os.MkdirAll(uploadDir, os.ModeDir|os.ModePerm); err != nil {
		err = fmt.Errorf("MkdirAll() error: %v", err)
		return
	}

	r := gin.Default()

	r.POST("/", uploadFile)

	r.Run(config.ServerAddr)
}
