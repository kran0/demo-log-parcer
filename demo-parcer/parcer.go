package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/satyrius/gonx"
)

type AppConfig struct {
	Limit         int    `default:"10"`
	FileName      string `required:"true"`
	LogFormat     string `default:"$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent $request_time \"$http_referer\" \"$http_user_agent\" [upstream: $upstream_addr $upstream_status] request_id=$upstream_http_x_request_id"`
	HumanReadable bool   `default:"false"`
}

type Score struct {
	Remote_addr     string
	Body_bytes_sent uint64
}

var format string
var logFile string

func main() {

	// Get app config from env
	var c AppConfig
	err := envconfig.Process("parcer", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set parcer params
	flag.StringVar(&format, "format", c.LogFormat, "Log format")
	flag.StringVar(&logFile, "log", c.FileName, "Log file name to read. Read from STDIN if file name is '-'")
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

		body_bytes_sent, _ := res.Field("body_bytes_sent")
		score.Body_bytes_sent, _ = strconv.ParseUint(strings.TrimSuffix(body_bytes_sent, ".00"), 10, 64)

		scores = append(scores, *score)
	}

	// Sort []Scores
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Body_bytes_sent > scores[j].Body_bytes_sent
	})

	// Set the limit
	limit := len(scores)
	if c.Limit < limit {
		limit = c.Limit
	}

	// Print the report
	for rep := range scores[:limit] {
		if c.HumanReadable {
			fmt.Println(scores[rep].Remote_addr,
				bytefmt.ByteSize(scores[rep].Body_bytes_sent))
		} else {
			fmt.Println(scores[rep].Remote_addr)
		}
	}

}
