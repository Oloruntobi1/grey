FROM alpine:latest as build
RUN apk add --no-cache ca-certificates

FROM scratch
COPY main /main
COPY dev.env dev.env
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

ENV MY_ENV=development
EXPOSE 9191
EXPOSE 6060
USER nobody

ENTRYPOINT ["/main"]