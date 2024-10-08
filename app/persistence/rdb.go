package persistence

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shashjar/redis-in-go/app/store"
)

var RDB_DIR string
var RDB_FILENAME string

const DEFAULT_RDB_DIR = "/redis-data"
const DEFAULT_RDB_FILENAME = "dump.rdb"

// NOTE: for ease of implementation, uses a custom RDB format (see redis-data/example.rdb) instead of the official
// RDB binary file format (https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/ + https://rdb.fnordig.de/file_format.html)

// Persists the data stored in the RDB file at the given path into this running server's key-value store
func PersistFromRDB(filePath string) {
	lines, err := readFile(filePath)
	if err != nil {
		log.Println("Unable to persist from RDB file into Redis server:", err.Error())
		return
	}
	processRDBKeyValuePairs(lines)
}

// Dumps the current state of this running server's key-value store to an RDB file that is saved to disk
func DumpToRDB() error {
	filepath := "." + RDB_DIR + "/" + RDB_FILENAME
	rdbBytes := GetRDBBytes()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(rdbBytes)
	if err != nil {
		return err
	}
	file.Sync()

	return nil
}

// Gets a byte slice representing the data that would be written to an RDB file
// NOTE: only dumps KV pairs of type "string" into the RDB file (does not persist streams)
func GetRDBBytes() []byte {
	data := store.Data()
	numKeyValuePairs := 0
	var bytes []byte

	for key, value := range data {
		if value.Type == "string" {
			if value.IsExpired() {
				store.DeleteKey(key)
				continue
			}

			bytes = append(bytes, []byte(key+"\n")...)
			bytes = append(bytes, []byte(value.Value.(string)+"\n")...)
			if value.HasExpiration() {
				bytes = append(bytes, []byte(strconv.Itoa(int(value.Expiration.Unix()))+"\n")...)
			}
			bytes = append(bytes, []byte("\n")...)

			numKeyValuePairs += 1
		}
	}
	bytes = bytes[:len(bytes)-1]

	bytes = append([]byte(strconv.Itoa(numKeyValuePairs)+"\n\n"), bytes...)
	return bytes
}

func readFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
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

	return lines, nil
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

		store.Set(key, value, expiresAt)
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
