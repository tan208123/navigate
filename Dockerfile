FROM centos:6.9

COPY main /usr/bin

EXPOSE 12345

ENTRYPOINT ["main"]
