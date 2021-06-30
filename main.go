package main

import (
	"flag"
	"log"
	"regexp"
	"time"
)

func generateExcludeList(source *string) *[]string {
	splitter := regexp.MustCompile("\\s*,\\s*")
	nullpath := regexp.MustCompile("^\\s*$")

	parts := splitter.Split(*source, -1)

	output := make([]string, 0)

	for _, part := range parts {
		if len(part) > 0 && !nullpath.MatchString(part) {
			output = append(output, part)
		}
	}
	return &output
}

func main() {
	excludeListRawPtr := flag.String("exclude", "/Applications,/System,/dev", "comma-separated list of paths to exclude")
	serverPtr := flag.String("server", "https://localhost:9000", "base URL of the server to send data to")
	noSendPtr := flag.Bool("nosend", false, "if set don't try to upload anything")
	startPathPtr := flag.String("start", "/System/Volumes/Data", "start recursive scan from this path")
	flag.Parse()

	excludeListPtr := generateExcludeList(excludeListRawPtr)

	log.Printf("Exclude list is %v", *excludeListPtr)
	matcher := regexp.MustCompile("(.*\\.prproj|.*\\.plproj|.*\\.aep$|.*\\.cpr)")
	startTime := time.Now()
	filesCh, scanErrCh := AsyncScanner(startPathPtr, matcher, excludeListPtr)
	sendErrCh := AsyncSender(filesCh, *serverPtr, *noSendPtr)

	func() {
		for {
			select {
			//case file := <-filesCh:
			//	if file == nil {
			//		endTime := time.Now()
			//		runtime := endTime.Sub(startTime)
			//		log.Printf("All done, run completed in %s", runtime.String())
			//		return
			//	}
			//	log.Printf("Got %s of size %d", file.FullPath, file.Size())
			case err := <-scanErrCh:
				log.Print("ERROR: ", err)
				return
			case err := <-sendErrCh:
				if err == nil {
					endTime := time.Now()
					runtime := endTime.Sub(startTime)
					log.Printf("All done, run completed in %s", runtime.String())
					return
				} else {
					log.Print("ERROR: ", err)
					return
				}
			}
		}
	}()
}
