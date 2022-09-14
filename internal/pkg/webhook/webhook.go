package webhook

import (
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	Type   string `json:"type"`
}

var typeMapping = map[string]string{
	"text": "text/plan; charset=UTF-8",
	"html": "text/html; charset=UTF-8",
	"json": "application/json; charset=utf-8",
}

func Send(conf Config, data io.Reader) error {
	if "" == conf.Method {
		conf.Method = http.MethodPost
	}
	req, err := http.NewRequest(conf.Method, conf.Url, data)
	if nil != err {
		return err
	}
	if "" != conf.Type {
		req.Header.Set("Content-Type", conf.Type)
	}

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return fmt.Errorf("http request failed status code %d", resp.StatusCode)
	}

	//b, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(b))

	return nil
}
