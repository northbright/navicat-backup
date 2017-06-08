# navicat-backup

[![Build Status](https://travis-ci.org/northbright/navicat-backup.svg?branch=master)](https://travis-ci.org/northbright/navicat-backup)
[![Go Report Card](https://goreportcard.com/badge/github.com/northbright/navicat-backup)](https://goreportcard.com/report/github.com/northbright/navicat-backup)

Backup tools for Navicat including both client and server.The tools are written in [Golang](https://golang.org).

#### Use Case
* Navicat has a daily backup schedule(backup to .psc files).
* Run server on one or more remote PCs.
* Run client as scheduled task to upload latest local backup file to remote servers.

#### Client
* Configure Client via `client/config.json`:

        {
            "navicat_backup_dir":"D:\\navicat-bak\\data",
            "remote_upload_urls":[
                    "http://192.168.1.100"
            ]
        }

  * `"navicat_backup_dir"` is navicat backup folder.
  * `"remote_upload_urls"` is the URLs of remote servers to upload latest backup file(.psc).

* Run client as scheduled task(e.g. add a new schedule job on Windows).

#### Server
* Configure Server via `server/config.json`:

        {
            "server_addr":":80"
        }

  * `"server_addr"` is TCP network address to listen.

* Run server as service on remote PCs.

#### Build from Source
* Client

        cd client

        // For Windows x86
        env GOOS=windows GOARCH=386 go build -v

        // For Mac OSX x64
        env GOOS=darwin GOARCH=amd64 go build -v

        // For Linux x64
        env GOOS=linux GOARCH=amd64 go build -v

* Server

        cd ../server
        
        // For Windows x86
        env GOOS=windows GOARCH=386 go build -v

        // For Mac OSX x64
        env GOOS=darwin GOARCH=amd64 go build -v

        // For Linux x64
        env GOOS=linux GOARCH=amd64 go build -v

#### License
* [MIT License](LICENSE)
