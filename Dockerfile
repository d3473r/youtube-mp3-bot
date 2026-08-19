FROM golang:alpine AS build
WORKDIR /src
COPY . .
RUN go build -o youtube-mp3-bot

FROM alpine:3.24.1
RUN apk --no-cache add yt-dlp ffmpeg
WORKDIR /app
COPY --from=build /src/youtube-mp3-bot /app/
ENTRYPOINT ["./youtube-mp3-bot"]
