# Build on golang:1.9
FROM golang:1.9

# Create a diretory to save our files
RUN mkdir -p /go/src/service-end
WORKDIR /go/src/service-end

# Copy source code
COPY . /go/src/service-end

# Download dependencies
RUN go get -u github.com/spf13/pflag
#RUN go get -u github.com/sysu-saad-project/service-end/core/service
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/gorilla/websocket
RUN go get -u github.com/urfave/negroni
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/go-xorm/xorm
RUN go-wrapper download
RUN go-wrapper install

# Setting ENV Value
ENV PORT 8080

# Expose 8080 port to communicate with host

# Containner Begin
CMD ["go-wrapper", "run"]