package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	TARGET                 = "http://127.0.0.1:80/status"
	TARGET_EXPECTED_STATUS = 200
	INTERVAL               = 60 * time.Second
)

func main() {
	ticker := time.NewTicker(INTERVAL)
	for {

		<-ticker.C

		resp, err := http.Get(TARGET)
		if err != nil {
			reportError(err)
			continue
		}

		if resp.StatusCode != TARGET_EXPECTED_STATUS {
			reportError(nil)
			continue
		}

		reportSucceed()

	}
}

func reportError(checkError error) {
	fmt.Printf("check error happens: %#v\n", checkError) // TODO
}

func reportSucceed() {
	fmt.Printf("check succeeded\n") // TODO
}
