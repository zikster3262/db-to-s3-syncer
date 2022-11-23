package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	uuid string
	time string
}

func main() {

	ch := make(chan Request)

	ctx := context.Background()

	for i := 0; i < 20; i++ {
		go generateRequest(ch)
	}

	for {
		select {
		case r := <-ch:
			go send(&r)
		case <-ctx.Done():
			return
		}
	}

}

func generateRequest(r chan Request) {
	rq := Request{
		uuid: uuid.New().String(),
		time: time.Now().String(),
	}
	r <- rq
}

func send(r *Request) {
	time.Sleep(time.Second * 1)
	//Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"uuid": r.uuid,
		"time": r.time,
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:8080/api/v1/rq", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
