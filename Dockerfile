FROM benschw/litefs

ADD app /opt/app

EXPOSE 8080

ENTRYPOINT ["/opt/app"]


