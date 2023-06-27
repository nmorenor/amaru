package net

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AvailableSession struct {
	ID              string `json:"id"`
	SessionHostName string `json:"name"`
	Size            int    `json:"size"`
}

func GetAvailableSessions() *[]AvailableSession {
	resp, err := http.Get(AvailableSessionsURL)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var sessions []AvailableSession
	err = json.Unmarshal(body, &sessions)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &sessions
}
