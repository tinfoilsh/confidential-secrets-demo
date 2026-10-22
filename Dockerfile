FROM golang:1.27.3@sha256:f20d83a83e18dbd2a4714b8ffd4b30b287b6243962dd75608e86abcc63ba4efc AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -o /demo .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /demo /demo
EXPOSE 8080
ENTRYPOINT ["/demo"]
