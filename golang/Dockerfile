FROM golang:alpine
COPY run.go .
RUN go build -o /tree/app .

FROM alpine
EXPOSE 80
WORKDIR /tree
COPY --from=0 /tree/app .
CMD [ "./app" ]