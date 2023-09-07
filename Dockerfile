FROM golang:1.21.1-bullseye as builder

WORKDIR /src

ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Final Docker image
FROM alpine:3.18 AS final-stage

RUN apk add --update --no-cache ca-certificates
# Create user validation
RUN addgroup -S standards && adduser -u 1234 -S standards -G standards
# must be numeric to work with Pod Security Policies:
# https://kubernetes.io/docs/concepts/policy/pod-security-policy/#users-and-groups
USER 1234
WORKDIR ${HOME}/app
COPY --from=builder /src/standards-insights .

EXPOSE 3000

ENTRYPOINT ["./standards-insights"]







