# Example
This is a sample process to be ran for testing in conjunction with layover.

## Usage
Source example.env for layover
> source example.env

Run the example http server provided here
> go run layover/example

This will serve up http requests on 127.0.0.1:8080

Run layover
> layover proxy

This will proxy requests from 127.0.0.1:8081

Make an http request to the proxied port
> http get http://127.0.0.1:8081/test?msg=hello

You will then receive your message 'hello' back.