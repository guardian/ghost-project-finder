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
	flag.Parse()

	excludeListPtr := generateExcludeList(excludeListRawPtr)

	startPath := "/"

	log.Printf("Exclude list is %v", *excludeListPtr)
	matcher := regexp.MustCompile("(.*\\.prproj|.*\\.plproj|.*\\.aep$|.*\\.cpr)")
	startTime := time.Now()
	filesCh, errCh := AsyncScanner(&startPath, matcher, excludeListPtr)

	func() {
		for {
			select {
			case file := <-filesCh:
				if file == nil {
					endTime := time.Now()
					runtime := endTime.Sub(startTime)
					log.Printf("All done, run completed in %s", runtime.String())
					return
				}
				log.Printf("Got %s of size %d", file.FullPath, file.Size())
			case err := <-errCh:
				log.Print("ERROR: ", err)
				return
			}
		}
	}()
}
