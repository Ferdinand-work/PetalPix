FROM alpine:latest
WORKDIR /app
COPY linux-app config.yml /app/
RUN chmod +x linux-app
EXPOSE 9095
CMD ["./linux-app"]