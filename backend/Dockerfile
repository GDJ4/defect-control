FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o defect-tracker ./cmd/api

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /srv/app
COPY --from=builder /app/defect-tracker .

EXPOSE 8080
ENTRYPOINT ["./defect-tracker"]
