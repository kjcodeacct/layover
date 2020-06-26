# This is an example Dockerfile primarily for testing but if you would like to build your own
#   release of layover, please see the README.md section 'Installing'
# $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o layover
# $ docker build . -t layover:test
# $ docker run -e LAYOVER_PROXYPORT=8090 --publish 8080:8080 --name test1 layover:test

FROM scratch

# Copy our static executable.
COPY layover /

# Run the layover binary.
ENTRYPOINT ["/layover"]