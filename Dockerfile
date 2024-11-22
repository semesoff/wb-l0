FROM ubuntu:latest
LABEL authors="semesoff"

ENTRYPOINT ["top", "-b"]