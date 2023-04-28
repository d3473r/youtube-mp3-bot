# youtube-mp3-bot

![Build Status](https://img.shields.io/github/actions/workflow/status/d3473r/youtube-mp3-bot/docker-image.yml)
![Latest Version](https://ghcr-badge.egpl.dev/d3473r/youtube-mp3-bot/latest_tag?trim=major&label=latest)
![Docker Image Size (tag)](https://ghcr-badge.egpl.dev/d3473r/youtube-mp3-bot/size?tag=main)

## Usage

- `git clone https://github.com/d3473r/youtube-mp3-bot.git`
- `cd youtube-mp3-bot`
- `cp .env.example .env`
- Enter your Telegram Bot API Token in the `.env` file 
- `docker-compose up -d`
- Send `www.youtube.com`, `m.youtube.com` or `youtu.be` links to the Bot and it will respond with a mp3 of that video

## Caveats

- [If the mp3 of the video is >50Mb the Bot cannot send the mp3](https://core.telegram.org/bots/faq#how-do-i-upload-a-large-file), You can get the file from the download folder though.
