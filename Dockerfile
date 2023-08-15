FROM golang:alpine AS BUILD
WORKDIR /src
COPY . .
RUN go build -o youtube-mp3-bot

FROM alpine:3.18
RUN apk --no-cache add yt-dlp ffmpeg
WORKDIR /app
COPY --from=BUILD /src/youtube-mp3-bot /app/
ENTRYPOINT ./youtube-mp3-bot