package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
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
	SLACK    = "slack"
	TELEGRAM = "telegram"
)

type SlackRequest struct {
	text string
}

type Slack struct{}

func (c *Slack) Send(message string) {
	var jsonStr = []byte(fmt.Sprintf(`{"text": "%s"}`, message))
	fmt.Println(os.Getenv("SLACK_WEBHOOK_URL"))
	res, err := http.Post(os.Getenv("SLACK_WEBHOOK_URL"), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	} else {
		if res.StatusCode != 200 {
			fmt.Println(fmt.Sprintf("[%s] 에러 메시지 전송에 실패했습니다: %s", strconv.Itoa(res.StatusCode), err))
		} else {
			fmt.Println(fmt.Sprintf("[%s] 에러 메시지가 전송되었습니다: %s", strconv.Itoa(res.StatusCode), message))
		}
	}
}

func getNotificationClient(notificationType string) NotificationClient {
	switch notificationType {
	case SLACK:
		return new(Slack)
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

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
	fmt.Println("start!")
	fmt.Printf("start!!")
	loadEnv(".env")
	body := Body{}
	err := json.Unmarshal([]byte(string(request.Body)), &body)
	if err != nil {
		return "json parse error", err
	}
	clients := GetClients(body.Types)
	message := fmt.Sprintf("[%s] 문제가 발생했습니다", body.ServiceName)
	for _, c := range clients {
		c.Send(message)
	}
	return message, nil
}

//func test() (string, error) {
//	loadEnv(".env")
//	types := make([]string, 0)
//	types = append(types, "slack")
//	clients := GetClients(types)
//	message := fmt.Sprintf("[%s] 문제가 발생했습니다", "test")
//	for _, c := range clients {
//		c.Send(message)
//	}
//	return "success", nil
//}

func main() {
	//test()
	lambda.Start(handler)
}
