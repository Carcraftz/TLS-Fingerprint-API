# TLS-Fingerprint-API
A server that proxies requests and uses CycleTLS to modify your clienthello and prevent your requests from being fingerprinted.

#How to use:

Deploy this server somewhere. Localhost is preferrable to reduce latency. 


Modify your code to make requests to the server INSTEAD of the endpoint you want to request. Ex: If running on localhost, make requests to http://127.0.0.1:8082. Make sure to also remove any code that uses a proxy in the request.


Add the request header "poptls-url", and set it equal to the endpoint you want to request. For example, if you want to request https://httpbin.org, you would add the header "poptls-url" = "https://httpbin.org"


Optional: Add the request header "poptls-proxy" and set it equal to the URL for the proxy you want to use (format: http://user:pass@host:port or http://host:port). This will make the server use your proxy for the request.

