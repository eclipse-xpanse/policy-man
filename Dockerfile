FROM frolvlad/alpine-glibc

RUN addgroup -S policy-man && adduser -S -G policy-man policy-man

COPY policy-man /usr/bin/policy-man
USER policy-man
WORKDIR /home/policy-man

ENTRYPOINT ["policy-man"]