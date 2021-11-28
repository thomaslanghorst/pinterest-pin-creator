package schedule

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	timestampLayout = time.RFC1123 //Mon, 02 Jan 2006 15:04:05 MST
)

type NextPinData struct {
	Created     bool
	Timestamp   time.Time
	BoardName   string
	Title       string
	Description string
	ImagePath   string
	Link        string
	Index       int
}

type ScheduleReaderInterface interface {
	Next() (*NextPinData, error)
	SetCreated(index int) error
}

type ScheduleReader struct {
	filePath string
}

func NewScheduleReader(filePath string) *ScheduleReader {
	return &ScheduleReader{
		filePath: filePath,
	}
}

func (r *ScheduleReader) Next() (*NextPinData, error) {

	now := time.Now()

	allLines, err := readFile(r.filePath)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(allLines); i++ {

		// skip header
		if i == 0 {
			continue
		}

		line := allLines[i]

		created, err := strconv.ParseBool(line[0])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to parse created value %s in csv file. Error: %s", line[0], err.Error()))
		}

		if created == true {
			continue
		}

		timestamp, err := time.Parse(time.RFC1123, line[1])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unable to parse timestamp %s in csv file. Error: %s", line[1], err.Error()))
		}

		if timestamp.After(now) {
			continue
		}

		nextPinData := &NextPinData{}
		nextPinData.Index = i
		nextPinData.Created = created
		nextPinData.Timestamp = timestamp
		nextPinData.BoardName = line[2]
		nextPinData.Title = line[3]
		nextPinData.Description = line[4]
		nextPinData.ImagePath = line[5]
		nextPinData.Link = line[6]

		return nextPinData, nil

	}

	return nil, nil
}

func (r *ScheduleReader) SetCreated(index int) error {
	allFiles, err := readFile(r.filePath)
	if err != nil {
		return err
	}

	allFiles[index][0] = "true"

	err = writeFile(r.filePath, allFiles)
	if err != nil {
		return err
	}
	return nil
}

func readFile(csvFile string) ([][]string, error) {

	csvfile, err := os.Open(csvFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read csv file. Error: %s", err.Error()))
	}

	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.Comma = ';'

	allLines, err := r.ReadAll()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read csv file. Error: %s", err.Error()))
	}

	return allLines, nil

}

func writeFile(csvFile string, allLines [][]string) error {
	csvfile, err := os.OpenFile(csvFile, os.O_WRONLY, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to write csv file. Error: %s", err.Error()))
	}

	defer csvfile.Close()

	w := csv.NewWriter(csvfile)
	w.Comma = ';'

	defer w.Flush()

	err = w.WriteAll(allLines)

	if err != nil {
		return errors.New(fmt.Sprintf("unable to write csv file. Error: %s", err.Error()))
	}

	return nil
}
