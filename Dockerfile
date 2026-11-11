FROM oven/bun:1 AS web
WORKDIR /w
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile
COPY web/ ./
RUN bun run build

FROM golang:1.26-alpine AS server
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=server /out/server /server
COPY --from=web /w/dist /web/dist
ENV WEB_DIR=/web/dist PORT=8080
EXPOSE 8080
ENTRYPOINT ["/server"]
