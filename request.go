package main

import (
	"errors"
	"strconv"
	"strings"
)

type Request struct {
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	Method      string            `json:"method"`
	RequestPath string            `json:"-"`
}

func NewRequest(mqttAnswerChannel string, body string) Request {
	request := Request{
		Headers: map[string]string{
			"mqtt_answer_channel": mqttAnswerChannel,
		},
		Body: body,
	}

	return request
}

func NewSocketRequest(body string) Request {
	request := Request{
		Headers: map[string]string{
			"is_socket_message": "true",
		},
		Body: body,
	}

	return request
}

func (req *Request) GetAnswerChannel() (string, error) {
	if req.IsSocketMessage() {
		return "", errors.New("the request does not have an answer channel, because it is a socket message")
	}

	if answerChannel, exists := req.Headers["mqtt_answer_channel"]; exists {
		if channel := strings.TrimSpace(answerChannel); len(channel) != 0 {
			return channel, nil
		}
		return "", errors.New("the answer channel is empty")
	}
	return "", errors.New("the request does not have an answer channel")
}

func (req *Request) IsSocketMessage() bool {
	if socket_message, exists := req.Headers["is_socket_message"]; exists {
		if is_socket_message, err := strconv.ParseBool(socket_message); err == nil {
			return is_socket_message
		}
		return false
	}

	if _, exists := req.Headers["mqtt_answer_channel"]; exists {
		return false
	}
	return true
}
