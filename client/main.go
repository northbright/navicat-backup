package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"

	"github.com/northbright/httputil"
	"github.com/northbright/pathhelper"
)

// Config is the client configuration.
type Config struct {
	NavicatBackupDir string   `json:"navicat_backup_dir"`
	RemoteUploadURLs []string `json:"remote_upload_urls"`
}

var (
	execDir    = ""
	configFile = ""
	config     Config
)

func init() {
	execDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(execDir, "config.json")
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

func getLatestBackupFile() (string, error) {
	f, err := os.Open(config.NavicatBackupDir)
	if err != nil {
		return "", err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return "", err
	}

	if len(names) <= 0 {
		return "", fmt.Errorf("no backup file found")
	}

	// Sort names.
	sort.Strings(names)

	// Get latest file name which matches the "YYMMDDHHMMSS.psc" pattern(e.g. 170601220002.psc).
	p := `^\d{12}\.psc$`
	re := regexp.MustCompile(p)

	l := len(names)
	for i := l - 1; i >= 0; i-- {
		if re.MatchString(names[i]) {
			return names[i], nil
		}
	}

	return "", fmt.Errorf("no psc file found")
}

func uploadBackupFile(latestBackupFile string) error {
	p := path.Join(config.NavicatBackupDir, latestBackupFile)

	for _, uri := range config.RemoteUploadURLs {
		req, err := httputil.NewUploadFileRequest("POST", uri, p, "upload", nil)
		if err != nil {
			err = fmt.Errorf("httputil.NewUploadFileRequest() error: %v", err)
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			err = fmt.Errorf("http.DefaultClient.Do() error: %v", err)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("Status Code is not ok: %v", resp.StatusCode)
			return err
		}
		log.Printf("upload %v to %v successfully", latestBackupFile, uri)
	}

	return nil
}

func main() {
	var err error

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if config, err = loadConfig(configFile); err != nil {
		err = fmt.Errorf("loadConfig() error: %v", err)
		return
	}

	f, err := getLatestBackupFile()
	if err != nil {
		err = fmt.Errorf("getLatestBackupFile() error: %v", err)
		return
	}

	log.Printf("latest file name: %v", f)

	if err = uploadBackupFile(f); err != nil {
		err = fmt.Errorf("uploadBackupFile() error: %v", err)
		return
	}
}
