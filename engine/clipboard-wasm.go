//go:build js
// +build js

package engine

import (
	"regexp"
	"syscall/js"
)

func InitClipboard() error {
	return nil
}

func WriteClipboard(text string) error {
	setResult := make(chan struct{}, 1)
	js.Global().Get("navigator").Get("clipboard").Call("writeText", text).
		Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				setResult <- struct{}{}
				return nil
			}),
		).Call("catch",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			println("failed to set clipboard: " + args[0].String())
			setResult <- struct{}{}
			return nil
		}),
	)
	<-setResult
	return nil
}

func ReadClipboard() (string, error) {
	resultChan := make(chan string, 1)
	js.Global().Get("navigator").Get("clipboard").Call("readText").
		Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				resultChan <- args[0].String()
				return nil
			}),
		).Call("catch",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			println("failed to read clipboard: " + args[0].String())
			resultChan <- ""
			return nil
		}),
	)
	val := <-resultChan
	return val, nil
}

var MOBILE_BROWSER_REGEX = regexp.MustCompile("(?i)Android|webOS|iPhone|iPad|iPod|BlackBerry|Windows Phone")

func IsMobileBrowser() bool {
	navigator := js.Global().Get("navigator")
	userAgent := navigator.Get("userAgent")
	return MOBILE_BROWSER_REGEX.Match([]byte(userAgent.String()))
}
