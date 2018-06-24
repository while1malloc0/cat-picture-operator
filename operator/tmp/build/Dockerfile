FROM alpine:3.6

# The operator user already existed for some reason
# that I don't feel like investigating
RUN adduser -D operator || true
USER operator

ADD tmp/_output/bin/operator /usr/local/bin/operator
