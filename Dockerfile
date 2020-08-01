FROM alpine
ADD web-analytics /web-analytics
ENTRYPOINT [ "/web-analytics" ]