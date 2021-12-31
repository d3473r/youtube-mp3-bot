FROM golang:alpine AS BUILD
WORKDIR /src
COPY . .
RUN go build -o youtube-mp3-bot

FROM alpine
RUN apk --no-cache add youtube-dl ffmpeg
WORKDIR /app
COPY --from=BUILD /src/youtube-mp3-bot /app/
ENTRYPOINT ./youtube-mp3-bot