FROM ubuntu:24.04

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Seoul
ENV CGO_ENABLED=1
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
ENV PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig

RUN apt-get update && apt-get install -y wget git build-essential pkg-config && \
    apt-get install -y libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev \
    gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl \
    gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio libglib2.0-dev

RUN if [ "$(uname -m)" = "x86_64" ]; then \
    wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz -O go.tar.gz; \
else \
    wget https://go.dev/dl/go1.24.1.linux-arm64.tar.gz -O go.tar.gz; \
fi && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

COPY . .

RUN go mod download

RUN mkdir -p /app/hls_output

# RUN GOARCH=$(go env GOARCH) go build -o api ./cmd/api/main.go
# RUN GOARCH=$(go env GOARCH) go build -o file_server ./cmd/file_server/main.go
