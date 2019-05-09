package redditvid

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Reddit Video scrapper
func RedditVideoBot() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("REDDIT_VID_BOT_SECRET_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "Quăng cho con bot này một link reddit có chứa video, nó sẽ lục lọi và trả về link video (nếu có thể). Cách sử dụng: `/reddit url_có_chứa_video`")
	})

	b.Handle("/reddit", func(m *tb.Message) {
		wantedUrl := parseJson(getJson(m.Payload + ".json"))

		//log.Print(m.Payload)
		//log.Print(wantedUrl)
		if wantedUrl != "" {
			b.Send(m.Sender, wantedUrl)
		}
	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		if q.Text != "" {
			// Declaration of G Translate url (our beloved sister)
			googleTranslateUrl := "https://translate.google.com.vn/translate_tts?ie=UTF-8&tl=vi&client=tw-ob&q="

			// List of result for inline query
			results := make(tb.Results, 1)

			// The one we need
			result := &tb.AudioResult{
				Title: q.Text,
				URL:   googleTranslateUrl + url.QueryEscape(q.Text),
			}

			results[0] = result
			results[0].SetResultID(strconv.Itoa(0))

			err := b.Answer(q, &tb.QueryResponse{
				Results:   results,
				CacheTime: 60, // in sec
			})

			if err != nil {
				fmt.Println(err)
			}
		}
	})

	b.Start()
}

func parseJson(json []byte) string {
	if json == nil {
		return ""
	}

	// Check if error was received
	output, err := jsonparser.GetString(json, "message")
	if err != nil {
		log.Print("Getting error message failed, everything seems right")
	} else {
		return output
	}

	output, err = jsonparser.GetString(json, "[0]", "data", "children", "[0]", "data", "secure_media", "reddit_video", "fallback_url")
	if err != nil {
		log.Print(string(json[:]))
		log.Print(err)
	} else {
		return output
	}

	return ""
}

func getJson(url string) []byte {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil
	}

	req.Header.Set("User-Agent", "telegram-bot:reddit-video:v1.0.0")

	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)
	} else {

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Print(err)
		} else {
			return body
		}
	}

	return nil
}

func main() {
	go RedditVideoBot()
}
