FROM golang:1.22.0
WORKDIR /app
COPY . .
WORKDIR /app/cmd/DiscordAI
RUN go mod download && go mod verify
RUN go build -o discordai
CMD [ "./discordai" ]