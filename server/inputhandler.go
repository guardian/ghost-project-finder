package main

import (
	"context"
	"encoding/base64"
	"github.com/guardian/ghost-project-finder/common"
	"github.com/olivere/elastic/v6"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type InputHandler struct {
	elasticSearchClient *elastic.Client
	indexName           string
	timeout             time.Duration
}

type IndexEntry struct {
	common.DataPacket

	Filename string
	Tags     []string
}

func makeDocId(from *common.DataPacket) string {
	maxLength := 512

	maybeId := base64.StdEncoding.EncodeToString([]byte(from.Hostname + ":" + from.Fullpath))
	firstGuessLength := len(maybeId)

	if firstGuessLength >= maxLength {
		charCountToRemove := firstGuessLength - maxLength
		cutPoint := (firstGuessLength / 2) - (charCountToRemove / 2)
		parts := []string{
			maybeId[0:cutPoint],
			maybeId[cutPoint+charCountToRemove:],
		}
		final := strings.Join(parts, "")
		return final
	} else {
		return maybeId
	}
}

func makeIndexEntry(from *common.DataPacket) IndexEntry {
	tagsToApply := []string{}

	if strings.HasPrefix(from.Fullpath, "/Downloads/Internet Downloads") {
		tagsToApply = append(tagsToApply, "InternetDownloads")
	}

	if strings.Contains(from.Fullpath, "Adobe Premiere Pro Auto-Save") {
		tagsToApply = append(tagsToApply, "Autosave")
	}
	if strings.HasPrefix(from.Fullpath, "/Users/Shared/Adobe/Premiere Pro/14.0/Tutorial/") {
		tagsToApply = append(tagsToApply, "Tutorials")
	}
	if strings.Contains(from.Fullpath, "Library/Preferences/Cubase 10/Project Templates/") {
		tagsToApply = append(tagsToApply, "Templates")
	}

	return IndexEntry{
		DataPacket: *from,
		Filename:   path.Base(from.Fullpath),
		Tags:       tagsToApply,
	}
}

func (h InputHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//requestUrl, _ := url.ParseRequestURI(r.RequestURI)

	if !AssertHttpMethod(r, w, "POST") {
		return //an error message has already been output
	}

	var incomingData common.DataPacket

	readErr := ReadJsonBody(r.Body, &incomingData)
	if readErr != nil {
		log.Printf("ERROR Could not understand incoming data: %s", readErr)
		responseMsg := GenericErrorResponse{
			Status: "bad_request",
			Detail: readErr.Error(),
		}
		WriteJsonContent(&responseMsg, w, 400)
		return
	}

	if incomingData.Hostname == "" || incomingData.Fullpath == "" {
		log.Printf("ERROR received an invalid data input: %v", incomingData)
		responseMsg := GenericErrorResponse{
			Status: "bad_request",
			Detail: "Need a hostname and a path",
		}
		WriteJsonContent(&responseMsg, w, 400)
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), h.timeout)
	defer cancelFunc()

	_, sendErr := h.elasticSearchClient.
		Index().
		Index(h.indexName).
		Type("ghost-project").
		Id(makeDocId(&incomingData)).
		BodyJson(makeIndexEntry(&incomingData)).
		Do(ctx)

	if sendErr != nil {
		log.Printf("ERROR could not write incoming data about %s from %s to index: %s", incomingData.Fullpath, incomingData.Hostname, sendErr)
		responseMsg := GenericErrorResponse{
			Status: "index_error",
			Detail: sendErr.Error(),
		}
		WriteJsonContent(&responseMsg, w, 500)
		return
	}

	responseMsg := GenericErrorResponse{
		Status: "ok",
		Detail: "indexed",
	}
	WriteJsonContent(&responseMsg, w, 200)
	return
}
