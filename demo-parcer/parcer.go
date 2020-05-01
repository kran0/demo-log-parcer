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
	InputFile           string `required:"true"`
	InputFileFormat     string `default:"$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent $request_time \"$http_referer\" \"$http_user_agent\" [upstream: $upstream_addr $upstream_status] request_id=$upstream_http_x_request_id"`
	OutputFile          string `default:"-"`
	OutputLimit         int    `default:"10"`
	OutputHumanReadable bool   `default:"false"`
}

type Score struct {
	Remote_addr     string
	Body_bytes_sent uint64
}

var format string
var logFile string

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
		panic(e)
	}
}

func main() {

	// Get app config from env
	var c AppConfig
	err := envconfig.Process("parcer", &c)
	check(err)
	// Set parcer params
	flag.StringVar(&format, "format", c.InputFileFormat, "Log format")
	flag.StringVar(&logFile, "log", c.InputFile, "Log file name to read. Read from STDIN if file name is '-'")
	flag.Parse()

	// Create a parser based on given format
	parser := gonx.NewParser(format)

	// Read given file or from STDIN
	var logReader io.Reader
	if logFile == "-" {
		logReader = os.Stdin
	} else {
		file, err := os.Open(logFile)
		check(err)
		logReader = file
		defer file.Close()
	}

	var OutWriter io.Writer
	if c.OutputFile == "-" {
		OutWriter = os.Stdout
	} else {
		file1, err := os.Create(c.OutputFile)
		check(err)
		OutWriter = file1
		defer file1.Close()
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
		score.Body_bytes_sent, err = strconv.ParseUint(strings.TrimSuffix(body_bytes_sent, ".00"), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping. Cannot convert str to uint64: %v", err)
		}

		scores = append(scores, *score)
	}

	// Sort []Scores
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Body_bytes_sent > scores[j].Body_bytes_sent
	})

	// Set the limit
	limit := len(scores)
	if c.OutputLimit < limit {
		limit = c.OutputLimit
	}

	// Print the report
	for rep := range scores[:limit] {
		if c.OutputHumanReadable {
			io.WriteString(OutWriter, scores[rep].Remote_addr+" "+bytefmt.ByteSize(scores[rep].Body_bytes_sent)+"\r\n")
		} else {
			io.WriteString(OutWriter, scores[rep].Remote_addr+"\r\n")
		}
	}

}
