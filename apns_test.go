package apns_http2

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

const (
	PEM_FILE = "/cert/apns/liurenyou.pem"
	//PEM_FILE = "/cert/apns/error.pem"

	PEM_PWD = "xz"
)

func TestApnsSingle(t *testing.T) {
	cli, err := New(PEM_FILE)
	if err != nil {
		t.Fatal("new err: \t", err)
	}

	message := time.Now().Format("2006-01-02 15:04:05")

	token := []string{
		/*
			"be7904b41a5db16294385fbb070f9fba0cffba14ea1626f5d9511dbd19be195b",
			"9478797d966211a624d1919077b3983ae9d910032ac95ba7a630e66ccc55b9e3", //xz
			"deb17ef7d5a0f65a3216289bc61eab154f4e83e7a57747785644268827091533", //me
			"5374058223d9713d8c6f31c4285ef91fabf3e0be139838c7205e8961ed7bff12", //5s
			"be7904b41a5db16294385fbb070f9fba0cffba14ea1626f5d9511dbd19be195b",
			"509952b74243a9d0db95db346013e64158160ef79e4a16181baa4cb2f65359a5", //useful
			"e5bfc7e7ad57d9f6af5a3dfe4ba9ceb26cf1e2351718f0f5788632dfd561ec9e",
			"b3b79681a0e177962bfe9956aa2431006ebee4e4862d909332e5fe30135721b5",
			"538538df1c9401aa3a544a884916a8dc126d7737f6a20a00a1c246acbc783bce",
			"77239ba1a985b555d45139f09a2fd192aa94734b2404066dc8fceaa6e77fdf4c",
			"b087b02ec0074a2c5c4a73260348320759871fad022efdf4a0eab710c9e79f3a",
		*/
		"22db81cf89826d7ebddc0ec92993dc9f323a8c0ef39637a974a4383e47aa5276",
		"8fe0327359a613cb30f1d3d2feee438235c892762168daf72dc4de9d34a14c39",
		"11c5e86764ed452c3dec0a911fc49885a8448b92843d071d8ac617b4adbee6a2",
		"1bf3a59d3ec9671e79636d02cb761f9f3015f6a1ab70ffd39f100ae86c3fa7b5",
	}
	data := fmt.Sprintf(`{"aps":{"alert":"%s","sound":"default","badge":1}, "type":"1", "url":"https://www.bing.com"}`, message)
	fail_num := 0

	var wg sync.WaitGroup
	queue := NewQueue(10, cli)

	go func() {
		for res := range queue.Response {
			log.Print(res.Err)
			wg.Done()
		}
	}()

	for _, t := range token {
		id := fmt.Sprintf("%d", time.Now().Unix())
		wg.Add(1)
		queue.Push(id, t, data, 86400)
	}

	wg.Wait()
	queue.Close()
	log.Print("fail num is  ", fail_num, "\t", message)
}
