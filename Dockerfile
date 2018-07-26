FROM golang:alpine
WORKDIR /tree
COPY backend .
RUN go build -o app .

FROM alpine
EXPOSE 80
WORKDIR /tree
COPY --from=0 /tree/app .
COPY data/data.json /data/data.json
CMD [ "./app", "--data=/data/data.json"]