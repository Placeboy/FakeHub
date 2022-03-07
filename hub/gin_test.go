package main

import (
	"github.com/Placeboy/FakeHub/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func BenchmarkSendDataRequest(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		GetData("heyhey","temperature")
//	}
//}
//
//func TestSendDataRequest(t *testing.T) {
//	GetData("heyhey", "temperature")
//}

func TestHandleDeviceRequest(t *testing.T) {
	// 假装一条设备消息
	cmd := "ReadTemperature"
	data, _ := rsa.RsaEncrypt([]byte(cmd))
	cryptgraph := base64.StdEncoding.EncodeToString(data)
	fmt.Println(cryptgraph)

	msg := "LinkLab-2022-" + cryptgraph
	id := "heyhey"
	fmt.Println(msg, id)

	router := setupRouter()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/handleDeviceRequest", nil)
	if err != nil {
		log.Fatal("http request fail ", err.Error())
	}

	q := req.URL.Query()
	q.Add("id", id)
	q.Add("msg", msg)
	req.URL.RawQuery = q.Encode()
	//req, _ := http.NewRequest("GET", "/ping", nil)

	router.ServeHTTP(w, req)
	fmt.Println(w.Body)
	//fmt.Println(w.Code)
}

