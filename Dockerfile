FROM golang:1.20-bullseye as server
RUN mkdir /app
WORKDIR /app
COPY ./server/. .
RUN CGO_ENABLED=0 go build -o /main .
RUN chmod +x /main


FROM node:18.12-bullseye as web
RUN mkdir /app
WORKDIR /app
COPY . .
RUN yarn
RUN yarn build

RUN rm -rf build/delay

FROM gcr.io/distroless/static-debian11
COPY --from=server /main /main
COPY --from=web /app/build /build
EXPOSE 8080
ENTRYPOINT [ "/main" ]