FROM ubuntu:latest
LABEL authors="mu"

ENTRYPOINT ["top", "-b"]