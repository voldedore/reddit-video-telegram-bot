package redditvid

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	b.Handle(tb.OnText, func(m *tb.Message) {
		b.Send(m.Sender, "Sorry, this feature is currently not supported.")
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
