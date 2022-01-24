package main

import (
	"bufio"
	"bytes"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

var YOUTUBE_DL = "yt-dlp"
var FFMPEG = "ffmpeg"
var PYTHON = "python"

func main() {
	if !commandExists(YOUTUBE_DL) || !commandExists(FFMPEG) {
		log.Fatal("youtube-dl, ffmpeg and python need to be installed")
		return
	}

	err := os.MkdirAll("download", os.ModePerm)

	token := os.Getenv("TELEGRAM_API_TOKEN")

	if token == "" {
		err = godotenv.Load()
		token = os.Getenv("TELEGRAM_API_TOKEN")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle(tb.OnText, func(message *tb.Message) {
		handleMessage(bot, message)
	})

	log.Println("Starting bot")
	bot.Start()
}

func handleMessage(bot *tb.Bot, message *tb.Message) {
	uri, err := url.Parse(message.Text)
	if err != nil {
		panic(err)
	}
	if uri.Host == "www.youtube.com" || uri.Host == "m.youtube.com" || uri.Host == "youtu.be" {
		log.Printf("Downloading: %s for %s: %d", uri, message.Sender.FirstName, message.Sender.ID)

		/*
			re := regexp.MustCompile("/(?:youtube.com/(?:[^/]+/.+/|(?:v|e(?:mbed)?)/|.*[?&]v=)|youtu.be/)([^\"&?/\\s]{11})")
			match := re.FindStringSubmatch(uri.String())
			if len(match) > 0 {
				log.Printf("Saving video id: %s", match[1])
			}
		*/

		cmd := exec.Command(YOUTUBE_DL, "--extract-audio", "--audio-format", "mp3", "-o", "download/%(title)s.%(ext)s", uri.String())

		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		scanner := bufio.NewScanner(stdout)
		scanner.Split(ScanLinesWithCarriageReturn)

		filename := ""
		filenameFound := false
		var response *tb.Message
		firstResponse := true
		lastUpdate := time.Now()
		for scanner.Scan() {
			line := scanner.Text()

			re := regexp.MustCompile(`\[download]  (.*)`)
			match := re.FindStringSubmatch(line)

			if len(match) > 0 {
				if firstResponse {
					response, _ = bot.Reply(message, match[1])
					firstResponse = false
				} else {
					if time.Now().After(lastUpdate.Add(time.Second)) {
						response, _ = bot.Edit(response, match[1])
						lastUpdate = time.Now()
					}
				}
			}

			if !filenameFound {
				re := regexp.MustCompile(`\[ExtractAudio] Destination: (.*).mp3`)
				match := re.FindStringSubmatch(line)

				if len(match) > 0 {
					filename = match[1] + ".mp3"
					filenameFound = true
				}
			}
		}

		cmd.Wait()

		if _, err := os.Stat(filename); err == nil {
			mp3 := &tb.Audio{File: tb.FromDisk(filename), FileName: filepath.Base(filename)}
			bot.Delete(response)
			bot.Reply(message, mp3)
		}
	} else {
		bot.Reply(message, "Host: "+uri.Host+" is not a youtube site")
	}
}

func ScanLinesWithCarriageReturn(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
