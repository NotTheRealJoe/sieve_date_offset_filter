package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func readToNewlineOr1K(reader *bufio.Reader) (line string, atEndOfFile bool) {
	var bs [1024]byte

	for i := 0; i < 1024; i++ {
		b, err := reader.ReadByte()
		if atEof(err) {
			return string(bs[0:i]), true
		}

		if b == '\n' {
			return string(bs[0:i]), false
		}

		bs[i] = b
	}

	// if we've used up all 1024 bytes, advance the reader to the next line, effectively
	// discarding the rest of the current one
	atEof := advanceToNextLine(reader)
	return string(bs[:]), atEof
}

func advanceToNextLine(reader *bufio.Reader) bool {
	for {
		b, err := reader.ReadByte()
		if atEof(err) {
			return true
		}

		if b == '\n' {
			return false
		}
	}
}

func atEof(err error) bool {
	if err == nil {
		return false
	}

	if err == io.EOF {
		return true
	}

	fmt.Printf("%s:%v\n", "Error in reading input.", err)
	// we exit with success on any *internal* errors, as exiting
	// with a failure status causes the email to be dropped
	os.Exit(0)

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
	reader := bufio.NewReader(os.Stdin)

	for {
		line, atEof := readToNewlineOr1K(reader)

		// remove the \r if the message used \r\n; also cut any spaces that
		// might be hanging around
		line = strings.Trim(line, " \r\n")

		if strings.Index(strings.ToLower(line), "date:") == 0 {
			dateField := strings.Trim(string(line[strings.Index(line, ":")+1:]), " ")
			msgTime := parseEmailTime(dateField)

			if msgTime == nil {
				fmt.Printf("Error getting date from email; dateField: %s", dateField)
				os.Exit(0)
			}

			curTime := time.Now()

			threshold := curTime.Add(time.Hour * 48)

			if msgTime.Equal(threshold) || msgTime.After(threshold) {
				// Test failed; message is too far in the future!
				os.Exit(1)
			}

			// Message time is OK.
			os.Exit(0)
		}

		if atEof {
			break
		}
	}
}
