package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const EXIT_CODE_OK = 1
const EXIT_CODE_ERR = 2
const EXIT_CODE_SPAM = 0

func readToNewlineOr1K(reader *bufio.Reader, logger *log.Logger) (line string, atEndOfFile bool) {
	var bs [1024]byte

	for i := 0; i < 1024; i++ {
		b, err := reader.ReadByte()
		if atEof(err, logger) {
			return string(bs[0:i]), true
		}

		if b == '\n' {
			return string(bs[0:i]), false
		}

		bs[i] = b
	}

	// if we've used up all 1024 bytes, advance the reader to the next line, effectively
	// discarding the rest of the current one
	atEof := advanceToNextLine(reader, logger)
	return string(bs[:]), atEof
}

func advanceToNextLine(reader *bufio.Reader, logger *log.Logger) bool {
	for {
		b, err := reader.ReadByte()
		if atEof(err, logger) {
			return true
		}

		if b == '\n' {
			return false
		}
	}
}

func atEof(err error, logger *log.Logger) bool {
	if err == nil {
		return false
	}

	if err == io.EOF {
		return true
	}

	logger.Printf("%s:%v\n", "Error in reading input.", err)
	// we exit with success on any *internal* errors, as exiting
	// with a failure status causes the email to be dropped
	os.Exit(EXIT_CODE_OK)

	// unreachable, but the compiler wants it
	return false
}

func parseEmailTime(timestamp string) *time.Time {
	// some emails add a timezone abbreviation in parenthesis after the UTC
	// offset, just cutting this off makes it compatible
	parenFound := strings.Index(timestamp, "(")
	if parenFound >= 0 {
		timestamp = timestamp[0:parenFound-1] + timestamp[strings.Index(timestamp, ")")+1:]
	}

	t, err := time.Parse(time.RFC1123Z, timestamp)
	if err == nil {
		return &t
	}

	t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", timestamp)
	if err == nil {
		return &t
	}

	t, err = time.Parse(time.RFC1123, timestamp)
	if err == nil {
		return &t
	}

	t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", timestamp)
	if err == nil {
		return &t
	}

	return nil
}

func main() {
	logger := log.New(os.Stderr, "", 0)
	reader := bufio.NewReader(os.Stdin)

	for {
		line, atEof := readToNewlineOr1K(reader, logger)

		// remove the \r if the message used \r\n; also cut any spaces that
		// might be hanging around
		line = strings.Trim(line, " \r\n")

		if strings.Index(strings.ToLower(line), "date:") == 0 {
			dateField := strings.Trim(string(line[strings.Index(line, ":")+1:]), " ")
			msgTime := parseEmailTime(dateField)

			if msgTime == nil {
				logger.Printf("Error getting date from email; dateField: %s", dateField)
				os.Exit(EXIT_CODE_OK)
			}

			curTime := time.Now()

			threshold := curTime.Add(time.Hour * 48)

			if msgTime.Equal(threshold) || msgTime.After(threshold) {
				// Test failed; message is too far in the future!
				fmt.Print("too far in future")
				os.Exit(EXIT_CODE_SPAM)
			}

			// Message time is OK.
			os.Exit(EXIT_CODE_OK)
		}

		if atEof {
			break
		}
	}
}
