FROM mcr.microsoft.com/vscode/devcontainers/go:latest

RUN apt-get update && \
    apt-get install -y netcat-openbsd

RUN go install github.com/bootdotdev/bootdev@latest
RUN go install github.com/spf13/cobra-cli@latest
