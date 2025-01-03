package main

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	httpStatusError = errors.New("http status error")
	httpMethodError = errors.New("http method error")
)

func main() {
	status, err := doHttp()
	if err = errors.Unwrap(err); err != nil {
		fmt.Println(err)
	}
	if errors.Is(err, httpMethodError) {
		fmt.Println(err)
		return
	}
	if errors.Is(err, httpStatusError) {
		fmt.Println(err)
		return
	}

	fmt.Println(status)
}

func doHttp() (string, error) {
	resp, err := http.Get("https://www.baidu.com/metrics")
	if err != nil {
		return "", fmt.Errorf("http get error: %w", httpMethodError)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http get error: %w", httpStatusError)
	}
	return resp.Status, nil
}
