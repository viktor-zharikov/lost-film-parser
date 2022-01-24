package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Series struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Cat     string   `xml:"category"`
	Pubdate string   `xml:"pubDate"`
	Link    string   `xml:"link"`
}

type Chanel struct {
	XMLName xml.Name `xml:"channel"`
	Series  []Series `xml:"item"`
}
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Chanel  Chanel   `xml:"channel"`
}

func download(y string) []byte {
	resp, err := http.Get(y)
	if err != nil {
		fmt.Println("fail load ", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", y, err)
	}

	return data
}

func main() {
	var client http.Client
	down := make(map[string]string)
	urls := make(map[string]interface{})
  list := "" //link to list with cinema. use only FULL ENGLISH name from lostfilm
	uid := "" //uid user from lostfilm
	usess := "" //usess from lostfilm
	API_KEY := "" // api key tg bot
	CHAT_ID := "" // chat id with tg bot

	cookie := &http.Cookie{
		Name:   "uid",
		Value:  uid,
		MaxAge: 300,
	}
	cookie2 := &http.Cookie{
		Name:   "usess",
		Value:  usess,
		MaxAge: 300,
	}

	file_downloaded, err := os.ReadFile("downloaded.txt")
	if err != nil {
		fmt.Println("fail open file with already downloaded", err)
		os.Create("downloaded.txt")
	}

	json.Unmarshal(file_downloaded, &urls)

	lines := strings.Split(string(download(list)), "\n")

	xmls := download("http://insearch.site/rssdd.xml")

	re := regexp.MustCompile(`.*\((.*)\)\..*`)
	var c rss
	xml.Unmarshal(xmls, &c)
	for _, x := range c.Chanel.Series {

		if x.Cat == "[1080p]" {

			down[re.ReplaceAllString(x.Title, "$1")] = x.Link
		}

	}

	for _, x := range lines {

		if down[x] != "" && urls[down[x]] != true {

			req, err := http.NewRequest("GET", down[x], nil)
			if err != nil {
				fmt.Println("Got error", err.Error())
			}
			req.AddCookie(cookie)
			req.AddCookie(cookie2)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error occured. Error is:", err.Error())
			}

			defer resp.Body.Close()

			re := regexp.MustCompile(`.*\=\"(.*)\"`)
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Read body error:", err)
			}
			fname := re.ReplaceAllString(resp.Header["Content-Disposition"][0], "$1")
			os.WriteFile(fname, data, 0766)

			body := strings.NewReader("chat_id=" + CHAT_ID + "&text=" + x + "\n" + fname)

			fname = "https://api.telegram.org/bot" + API_KEY + "/sendMessage"

			req, err = http.NewRequest("POST", fname, body)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err = client.Do(req)

			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()

			urls[down[x]] = true

			file, _ := json.MarshalIndent(urls, "", " ")

			_ = os.WriteFile("downloaded.txt", file, 0644)

		}
	}
}
