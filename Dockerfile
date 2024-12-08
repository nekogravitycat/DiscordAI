FROM golang:1.23.4
WORKDIR /app
COPY . .
WORKDIR /app/cmd/DiscordAI
RUN go mod download && go mod verify
RUN go build -o discordai
CMD [ "./discordai" ]