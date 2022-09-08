FROM golang:1.18
RUN mkdir /app 
ADD . /app/
ARG service_name
RUN go env -w GOPROXY="https://goproxy.cn,direct"
WORKDIR /app/microservice/$service_name
RUN make 
CMD ["./main"]
