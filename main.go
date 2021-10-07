package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Writer struct{}

func (Writer) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}

func main() {

	response, err := http.Get("https://google.com/")
	if err != nil {
		fmt.Println(err)
	}

	w := Writer{}
	io.Copy(w, response.Body)

	start := time.Now()
	servers := []string{
		"http://platzi.com",
		"http://google.com",
		"http://instagram.com",
		"http://facebook.com",
		"http://twitter.com",
	}

	for _, server := range servers {
		checkServer(server)
	}

	timeTaken := time.Since(start)
	fmt.Printf("exec time %s\n", timeTaken)

}

func checkServer(server string) {
	_, err := http.Get(server)
	if err != nil {
		fmt.Println(server, " is not working")
	} else {
		fmt.Println(server, " is working correctly")
	}
}
