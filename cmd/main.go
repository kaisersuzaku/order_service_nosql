package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	url := "http://localhost:7789/order-product"
	fmt.Println("URL:>", url)

	id := os.Args[1]
	loop, _ := strconv.Atoi(os.Args[2])
	start := time.Now()
	ch := make(chan string)
	for i := 0; i < loop; i++ {
		go hit(url, id, ch)
	}
	for i := 0; i < loop; i++ {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func hit(url, id string, ch chan<- string) {
	var jsonStr = []byte(`{"id":"` + id + `", "qty":1}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, len(body), url)
}
