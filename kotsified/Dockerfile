FROM golang AS builder

WORKDIR /go/src/app

RUN echo 'kotsified:x:1000:1000:,,,:/:' > /etc/passwd \
    && echo 'kotsified:x:1000:' > /etc/group

COPY . .

RUN CGO_ENABLED=0 go build -o kotsified

FROM scratch AS base

COPY --from=builder /go/src/app/kotsified /kotsified
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd

FROM scratch

COPY --from=base . .

USER kotsified:kotsified
    
CMD [ "/kotsified" ]
