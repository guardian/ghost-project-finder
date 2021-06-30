package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type DataPacket struct {
	Hostname string    `json:"hostname"`
	Fullpath string    `json:"fullpath"`
	Size     int64     `json:"size"`
	Modtime  time.Time `json:"modtime"`
}

func SendFileEntry(rec *FileEntry, hostname string, destBaseUrl string) error {
	toSend := DataPacket{
		Hostname: hostname,
		Fullpath: rec.FullPath,
		Size:     rec.Size(),
		Modtime:  rec.ModTime(),
	}

	toSendBytes, marshalErr := json.Marshal(&toSend)
	if marshalErr != nil {
		return marshalErr
	}

	log.Printf("DEBUG %s", string(toSendBytes))
	url := fmt.Sprintf("%s/foundfile", destBaseUrl)
	response, err := http.Post(url, "application/json", bytes.NewReader(toSendBytes))
	if err != nil {
		return err
	}

	io.Copy(ioutil.Discard, response.Body)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Server returned error %d", response.StatusCode))
	}
	return nil
}

func AsyncSender(inputCh chan *FileEntry, destBaseUrl string, nosend bool) chan error {
	myHostname, hostNameErr := os.Hostname()
	if hostNameErr != nil {
		log.Print("Could not get hostname: ", hostNameErr)
		myHostname = "unknown"
	}

	errCh := make(chan error, 1)

	go func() {
		for {
			entry := <-inputCh
			if entry == nil {
				log.Print("INFO AsyncSender got stream complete, exiting")
				errCh <- nil
				return
			}

			if nosend {
				log.Printf("INFO found %s of size %d from %s", entry.FullPath, entry.Size(), entry.ModTime())
			} else {
				err := SendFileEntry(entry, myHostname, destBaseUrl)
				if err != nil {
					log.Printf("ERROR AsyncSender could not send %s to server: %s", entry.FullPath, err)
				}
			}
		}
	}()

	return errCh
}
