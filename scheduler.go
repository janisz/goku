package main
import "net/http"
import (
	log "github.com/Sirupsen/logrus"
	"bytes"
	"bufio"
	"encoding/json"
	"strconv"
	"strings"
)

const MASTER = "http://172.17.0.7:5050"
const API_V1 = "/api/v1/scheduler"

type BaseEvent struct {
	Type string `json:"type"`
}

type SubscribedEvent struct {
	Type string `json:"type"`
	Subscribed string `json:"subscribed"`
}

type Subscribed struct {
	FrameworkId string `json:"framework_id"`
}

type FrameworkId struct {
	Value string `json:"value"`
}

func EventType(jsonBlob []byte) (string, error) {
	event := BaseEvent{}
	err := json.Unmarshal(jsonBlob, &event)
	if err != nil {
		return "", err
	} else if event.Type == "" {
		return "", nil
	} else {
		return event.Type, nil
	}
}

func Subscribe() error {
	SUBSCRIBE_BODY := `{
		"type": "SUBSCRIBE",
		"subscribe": {
		"framework_info": {
			"user" :  "root",
			"name" :  "Goku"
		},
		"force" : true
		}
	}`

	res, err := http.Post(MASTER + API_V1, "application/json", bytes.NewBuffer([]byte(SUBSCRIBE_BODY)))
	if err != nil {
		log.WithError(err)
		return err
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)
	line, err := reader.ReadString('\n')
	bytesCount, err := strconv.Atoi(strings.Trim(line, "\n"))
	if err != nil {
		log.WithError(err).Error("")
		return err
	}
	for {
		line, err = reader.ReadString('\n')
		line := strings.Trim(line, "\n")
		data := line[:bytesCount]
		bytesCount, err = strconv.Atoi((line[bytesCount:]))

		if err != nil {
			log.WithError(err).Error("")
			return err
		}
		eventType, err := EventType([]byte(data))
		if err != nil {
			log.WithError(err).Error("Invalid JSON")
			return err
		}

		switch eventType {
		case "ERROR":
			log.WithField("Type", eventType).Error(line)
		case "OFFERS":
			log.WithField("Type", eventType).Info(line)
		case "SUBSCRIBED":
			var sub SubscribedEvent
			json.Unmarshal([]byte(data), &sub)
			log.WithField("Type", eventType).Info(sub)
		case "HEARTBEAT":
			log.WithField("Type", eventType).Debug("ping")
		}
	}

	return nil
}