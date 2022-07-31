FROM golang:1.16
RUN mkdir /app 
ADD . /app/
ARG service_name
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY="https://goproxy.cn,direct"
WORKDIR /app/microservice/$service_name
RUN make 
CMD ["./main"]
