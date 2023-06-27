//go:build !js
// +build !js

package engine

import (
	"golang.design/x/clipboard"
)

func InitClipboard() error {
	return clipboard.Init()
}

func WriteClipboard(text string) error {
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}

func ReadClipboard() (string, error) {
	cvalue := clipboard.Read(clipboard.FmtText)
	return string(cvalue), nil
}

func IsMobileBrowser() bool {
	return false
}
