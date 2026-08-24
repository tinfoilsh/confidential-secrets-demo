FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -o /demo .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /demo /demo
EXPOSE 8080
ENTRYPOINT ["/demo"]
