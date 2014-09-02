package heroku

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Plugin struct {
}

type dyno struct {
	State string
}

func (this dyno) String() string {
	return this.State
}

func (this *Plugin) Authenticate() string {
	out, err := exec.Command("heroku", "auth:token").Output()

	if err != nil {
		log.Fatal(err)
	}

	withColon := ":" + string(out)
	return base64.StdEncoding.EncodeToString([]byte(withColon))
}

func (this *Plugin) Status(token string, app string) []int {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	herokuUrl := os.Getenv("HEROKU_API_URL")
	if herokuUrl == "" {
		herokuUrl = "https://api.heroku.com"
	}
	req, _ := http.NewRequest("GET", herokuUrl+"/apps/"+app+"/dynos", nil)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	dynoJson, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var dynos []dyno
	err = json.Unmarshal(dynoJson, &dynos)
	if err != nil {
		log.Fatal(err)
	}

	upCount := 0
	for _, dyno := range dynos {
		if dyno.State == "up" {
			upCount++
		}
	}

	return []int{upCount, len(dynos)}
}
