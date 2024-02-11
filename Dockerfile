FROM golang:1.21.7
WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN go build -o discordai
CMD [ "./discordai" ]