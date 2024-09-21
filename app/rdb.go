package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var RDB_DIR string
var RDB_FILENAME string

const DEFAULT_RDB_DIR = "/redis-data"
const DEFAULT_RDB_FILENAME = "dump.rdb"

// TODO: doesn't use the actual RDB binary file format - maybe implement this later. Refer to https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/
// & https://rdb.fnordig.de/file_format.html
func persistFromRDB() {
	filepath := "." + RDB_DIR + "/" + RDB_FILENAME
	lines := readFile(filepath)
	processRDBKeyValuePairs(lines)
}

func readFile(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Error reading RDB file:", err.Error())
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var lines []string
	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.TrimSpace(strings.Split(line, "//")[0])
		lines = append(lines, line)
	}

	return lines
}

func processRDBKeyValuePairs(lines []string) {
	groupedLines := splitOnDelimiter(lines, "")

	numKeyValuePairs, err := strconv.Atoi(groupedLines[0][0])
	if err != nil {
		log.Println("Error reading number of key-value pairs stored in RDB file:", err.Error())
		return
	}

	if numKeyValuePairs != len(groupedLines)-1 {
		log.Println("Number of key-value pairs in RDB file does not match written value")
		return
	}

	for _, lineGroup := range groupedLines[1:] {
		if len(lineGroup) < 2 || len(lineGroup) > 3 {
			log.Println("Invalid key-value pair configuration in RDB file")
			return
		}

		key := lineGroup[0]
		value := lineGroup[1]

		expiresAt := time.Time{}
		if len(lineGroup) == 3 {
			unixExpirationTimestamp, err := strconv.Atoi(lineGroup[2])
			if err != nil {
				log.Println("Invalid expiration timestamp provided in RDB file:", err.Error())
				return
			}

			expiresAt = time.Unix(int64(unixExpirationTimestamp), 0)
		}

		REDIS_STORE.Set(key, value, expiresAt)
	}
}

func splitOnDelimiter(lines []string, delimiter string) [][]string {
	var groupedLines [][]string

	currentGroup := []string{}
	for _, line := range lines {
		if line == delimiter {
			groupedLines = append(groupedLines, currentGroup)
			currentGroup = []string{}
		} else {
			currentGroup = append(currentGroup, line)
		}
	}

	if len(currentGroup) > 0 {
		groupedLines = append(groupedLines, currentGroup)
	}

	return groupedLines
}
