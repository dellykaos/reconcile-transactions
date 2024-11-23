FROM alpine:3.20

ENV APPUSER  appuser
ENV HOMEDIR  /app

RUN adduser -S -D -H -h ${HOMEDIR} ${APPUSER}

# ref https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

USER ${APPUSER}
WORKDIR ${HOMEDIR}

COPY /deploy/api .
EXPOSE 8080

CMD [ "./api" ]
