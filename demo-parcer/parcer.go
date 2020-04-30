package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/satyrius/gonx"
)

type Score struct {
	Remote_addr     string
	Body_bytes_sent string
}

var format string
var logFile string

func main() {

	// Set parcer params
	flag.StringVar(&format, "format",
		"$remote_addr - $remote_user [$time_local] "+
			"\"$request\" $status $body_bytes_sent $request_time "+
			"\"$http_referer\" \"$http_user_agent\" "+
			"[upstream: $upstream_addr $upstream_status] "+
			"request_id=$upstream_http_x_request_id", "Log format")

	flag.StringVar(&logFile, "log",
		"-", "Log file name to read. Read from STDIN if file name is '-'")

	flag.Parse()

	// Create a parser based on given format
	parser := gonx.NewParser(format)

	// Read given file or from STDIN
	var logReader io.Reader
	if logFile == "dummy" {
		logReader = strings.NewReader(`195.77.230.188 - - [29/Mar/2020:16:05:50 +0000] "GET /some-url HTTP/1.1" 200 1793 0.009 "https://example.com" "Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36" [upstream: 127.0.0.1:8000 200] request_id=-`)
	} else if logFile == "-" {
		logReader = os.Stdin
	} else {
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		logReader = file
		defer file.Close()
	}

	// Get stats from log file
	reducer := gonx.NewGroupBy(
		// Fields to group by
		[]string{"remote_addr"},
		// Result reducers
		&gonx.Sum{[]string{"body_bytes_sent"}},
	)

	results := gonx.MapReduce(logReader, parser, reducer)

	// Pull results to Score struct
	// Own struct - own rules!
	// We may implement custom methods or types or structs in our struct
	scores := []Score{}
	for res := range results {
		var score = new(Score)
		score.Remote_addr, _ = res.Field("remote_addr")
		score.Body_bytes_sent, _ = res.Field("body_bytes_sent")

		scores = append(scores, *score)
	}

	// Print the report
	for rep := range scores {
		fmt.Println(scores[rep].Remote_addr,
			scores[rep].Body_bytes_sent)
	}

}
