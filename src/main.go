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
	Send(message string, description string)
}

const (
	SLACK    = "slack"
	TELEGRAM = "telegram"
)

type SlackRequest struct {
	text string
}

type Slack struct{}

func (c *Slack) Send(message string, description string) {
	var jsonStr = []byte(fmt.Sprintf("{'text': '%s" + "\n>```%s```'}", message, description))
	fmt.Println(bytes.NewBuffer(jsonStr))
	fmt.Println(os.Getenv("SLACK_WEBHOOK_URL"))
	res, err := http.Post(os.Getenv("SLACK_WEBHOOK_URL"), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	} else {
		if res.StatusCode != 200 {
			fmt.Println(fmt.Sprintf("[%s] 에러 메시지 전송에 실패했습니다: %s", strconv.Itoa(res.StatusCode), err))
			fmt.Println(fmt.Sprintf("%s", res.Body))
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
	Description string   `json:"description"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("start!")
	fmt.Printf("start!!")
	loadEnv(".env")
	body := Body{}
	err := json.Unmarshal([]byte(string(request.Body)), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "json parse error"}, err
	}
	clients := GetClients(body.Types)
	message := fmt.Sprintf("[%s] 문제가 발생했습니다", body.ServiceName)
	description := fmt.Sprintf("%s", body.Description)
	for _, c := range clients {
		c.Send(message, description)
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "OK"}, nil
}

//func test() (string, error) {
//	loadEnv(".env")
//	types := make([]string, 0)
//	types = append(types, "slack")
//	clients := GetClients(types)
//	message := fmt.Sprintf("[%s] 문제가 발생했습니다", "test")
//	description := fmt.Sprintf("%s", `{"message":"timeout of 1ms exceeded","name":"Error","stack":"Error: timeout of 1ms exceeded\n    at createError (http://localhost:3000/static/js/vendors~main.chunk.js:49233:15)\n    at XMLHttpRequest.handleTimeout (http://localhost:3000/static/js/vendors~main.chunk.js:48740:14)","config":{"url":"https://api.ttbkk.com/api/places/count/?bottom_left=37.39570484631542%2C126.94976752096481&top_right=37.664390757822%2C127.09647377011702","method":"get","headers":{"Accept":"application/json, text/plain, */*"},"transformRequest":[null],"transformResponse":[null],"timeout":1,"xsrfCookieName":"XSRF-TOKEN","xsrfHeaderName":"X-XSRF-TOKEN","maxContentLength":-1,"maxBodyLength":-1},"code":"ECONNABORTED","request":{},"raw":{"message":"timeout of 1ms exceeded","name":"Error","stack":"Error: timeout of 1ms exceeded\n    at createError (http://localhost:3000/static/js/vendors~main.chunk.js:49233:15)\n    at XMLHttpRequest.handleTimeout (http://localhost:3000/static/js/vendors~main.chunk.js:48740:14)","config":{"url":"https://api.ttbkk.com/api/places/count/?bottom_left=37.39570484631542%2C126.94976752096481&top_right=37.664390757822%2C127.09647377011702","method":"get","headers":{"Accept":"application/json, text/plain, */*"},"transformRequest":[null],"transformResponse":[null],"timeout":1,"xsrfCookieName":"XSRF-TOKEN","xsrfHeaderName":"X-XSRF-TOKEN","maxContentLength":-1,"maxBodyLength":-1},"code":"ECONNABORTED"}}`)
//	for _, c := range clients {
//		c.Send(message, description)
//	}
//	return "success", nil
//}

func main() {
	//test()
	lambda.Start(handler)
}