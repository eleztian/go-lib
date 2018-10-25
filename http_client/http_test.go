package http_client

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type T struct {
}

func (t *T) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	w.Write(b)
	//w.Write([]byte(`{"Content":"Hello world"}`))
}

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := &T{}
		t.ServeHTTP(w, r)
	})
	go http.ListenAndServe(":8080", nil)
	time.Sleep(5 * time.Second)
}

func TestClient_ReqJson(t *testing.T) {
	c := New(true, Options{})
	d := struct {
		Content string
	}{"HELLO"}
	code, err := c.ReqJson("POST", "http://127.0.0.1:8080", nil, d, &d)
	if err != nil {
		t.Error(err)
		return
	}
	if code != 200 {
		t.Error("not ok")
		return
	}
	t.Log(d)
}
