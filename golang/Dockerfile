FROM golang:alpine
COPY run.go .
RUN go build -o /tree/app .

FROM alpine
EXPOSE 80
WORKDIR /tree
COPY global_tree_search_trees_1_2.csv .
COPY --from=0 /tree/app .
CMD [ "./app" ]