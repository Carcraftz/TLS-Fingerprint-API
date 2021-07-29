# TLS-Fingerprint-API

A server that proxies requests and uses my fork of CycleTLS & fhttp (fork of net/http) to prevent your requests from being fingerprinted. Built on open source software, this repo is a simple yet effective solution to censorship. It uses CycleTLS to spoof tls fingerprints, and fhttp to enable mimicry of chrome http/2 connection settings, header order, pseudo header order, and enable push.

## Support

I decided to make this after being tired of similar software being gatekept in the community, no one should have to pay over 3k for this. If you like my work, any support would be greatly appreciated ❤️
https://paypal.me/carcraftz?locale.x=en_US

## How to use:

Deploy this server somewhere. Localhost is preferrable to reduce latency.

Modify your code to make requests to the server INSTEAD of the endpoint you want to request. Ex: If running on localhost, make requests to http://127.0.0.1:8082. Make sure to also remove any code that uses a proxy in the request.

Add the request header "poptls-url", and set it equal to the endpoint you want to request. For example, if you want to request https://httpbin.org/get, you would add the header "poptls-url" = "https://httpbin.org/get"

Optional: Add the request header "poptls-proxy" and set it equal to the URL for the proxy you want to use (format: http://user:pass@host:port or http://host:port). This will make the server use your proxy for the request.

Optional: Add the request header "poptls-allowredirect" and set it to true or false to enable/disable redirects. Redirects are enabled by default.

## Run on a different Port:

By default the program runs on port 8082. You can specify another port by passing a flag --port=PORTNUM

## Examples:

### Node.js

To call this in node.js, lets say with node-fetch, you could do

````
fetch("http://localhost:8082",{
headers:{
"poptls-url":"https://httpbin.org/get",
"poptls-proxy":"https://user:pass@ip:port", //optional
"poptls-allowredirect:"true" //optional (TRUE by default)
}
})```
````
