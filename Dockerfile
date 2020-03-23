FROM golang:1.13.8
ADD . /go/src/scaha_micro_member
WORKDIR /go/src/scaha_micro_member
RUN go get scaha_micro_member
RUN go install
ENTRYPOINT ["/go/bin/scaha_micro_member"]