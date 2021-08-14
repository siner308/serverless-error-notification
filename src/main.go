package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func loadEnv(filepath string) {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}

type NotificationClient interface {
	Send(message string)
}

const (
	SLACK = "slack"
	TELEGRAM = "telegram"
)

type SlackRequest struct {
	text string
}

type Slack struct{}
func (c *Slack) Send(message string) {
	req := SlackRequest{message}
	reqBytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(reqBytes)
	_, err := http.Post(os.Getenv("SLACK_WEBHOOK_URL"), "application/json", buff)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("에러 메시지가 전송되었습니다:" + message)
	}
}


func getNotificationClient(notificationType string) NotificationClient {
	switch notificationType {
	case SLACK: return new(Slack)
	}
	return nil
}

func GetClients(types []string) []NotificationClient {
	results := make([]NotificationClient, 0)
	for _, t := range types {
		client := getNotificationClient(t)
		if client != nil {
			results = append(results, client)
		}
	}
	return results
}

type Body struct {
	ServiceName string   `json:"serviceName"`
	Types       []string `json:"types"`
}

func handler(ctx context.Context, body Body) (string, error) {
	loadEnv(".env")
	clients := GetClients(body.Types)
	message := fmt.Sprintf("[%s] 문제가 발생했습니다", body.ServiceName)
	for _, c := range clients {
		c.Send(message)
	}
	return "success", nil
}

func main() {
	lambda.Start(handler)
}
