package main

import (
	"errors"
	"fmt"
	"os"
	"pin-creator/accessToken"
	"pin-creator/config"
	"pin-creator/pinterest"
	"pin-creator/schedule"

	log "github.com/sirupsen/logrus"
)

var cfg *config.Config

func main() {

	readConfig()

	log.Infof("Checking for pins to create in '%s'", cfg.ScheduleFilePath)

	scheduleReader := schedule.NewScheduleReader(cfg.ScheduleFilePath)
	nextPinData, err := scheduleReader.Next()
	if err != nil {
		log.Fatal(err.Error())
	}
	if nextPinData == nil {
		log.Info("No pin scheduled for creation")
		return
	}

	createPin(nextPinData)

	err = scheduleReader.SetCreated(nextPinData.Index)
	if err != nil {
		log.Fatalf("Error setting pin created to true. Error: %s", err.Error())
	}
}

func readConfig() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("config.yaml file not provided")
	}

	configFilePath := args[1]

	cr := config.NewReader(configFilePath)
	c, err := cr.Read()
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg = c
}

func getToken() string {
	fmt.Println(cfg)
	tokenFileHandler := accessToken.NewAccessTokenFileHandler(cfg.AccessTokenPath)

	log.Info("Reading access token from file")
	token, err := tokenFileHandler.Read()
	if err == nil {
		return token
	} else {
		log.Info("No access token file found. Creating new token")

		tokenCreator := accessToken.NewAccessAccessTokenCreator(cfg.BrowserPath, cfg.RedirectPort)
		appId := os.Getenv("APP_ID")
		appSecret := os.Getenv("APP_SECRET")

		token, err := tokenCreator.NewToken(appId, appSecret)
		if err != nil {
			log.Fatalf("error creating new access token. Error: %s", err.Error())
		}

		log.Info("Writing access token to file")
		tokenFileHandler.Write(token)

		return token
	}
}

func createPin(scheduledPinData *schedule.NextPinData) {
	token := getToken()
	client := pinterest.NewClient(token)

	boards, err := client.ListBoards()
	if err != nil {
		log.Fatal(err.Error())
	}

	boardId, err := boardIdByName(boards, scheduledPinData.BoardName)
	if err != nil {
		log.Fatal(err.Error())
	}

	pinData := pinterest.PinData{
		BoardId:     boardId,
		ImgPath:     scheduledPinData.ImagePath,
		Link:        scheduledPinData.Link,
		Title:       scheduledPinData.Title,
		Description: scheduledPinData.Description,
		AltText:     scheduledPinData.Description,
	}

	err = client.CreatePin(pinData)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("Created Pin '%s' in board '%s'\n", pinData.Title, scheduledPinData.BoardName)
}

func boardIdByName(boards []pinterest.BoardInfo, boardName string) (string, error) {

	for _, board := range boards {
		if board.Name == boardName {
			return board.Id, nil
		}
	}

	return "", errors.New(fmt.Sprintf("board %s not found\n", boardName))
}
