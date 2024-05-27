FROM ubuntu:latest
LABEL authors="rdelper"

ENTRYPOINT ["top", "-b"]