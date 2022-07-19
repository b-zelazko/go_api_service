FROM golang:latest

RUN mkdir /build
WORKDIR /build
RUN cd /build && git clone https://github.com/b-zelazko/go_api_service.git
RUN cd /build/go_api_service && go build

EXPOSE 80/tcp

ENTRYPOINT [ "/build/go_api_service/go_api_service" ]