FROM alpine:latest
WORKDIR /app
COPY linux-app /app/
RUN chmod +x linux-app
EXPOSE 9095
CMD ["./linux-app"]