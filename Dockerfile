FROM golang:1.21-bullseye AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o myapp main.go

## Deploy
FROM alpine:latest

COPY --from=build /app/myapp /bin
COPY .env.prod /bin

EXPOSE 3000

ENTRYPOINT [ "/bin/myapp", "/bin/.env.prod" ]
# docker run -d -p 3000:3000 --name noteams-backend -v ~/.config/gcloud:/root/.config/gcloud noteam-backend-1-backend