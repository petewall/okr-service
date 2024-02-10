FROM golang AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go /app/
ADD cmd /app/cmd
ADD internal /app/internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /okr-service

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /okr-service /okr-service

EXPOSE 8080

USER nonroot:nonroot

VOLUME [ "/data" ]

ENTRYPOINT ["/okr-service"]
