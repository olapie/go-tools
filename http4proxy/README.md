openssl pkcs12 -in cert.p12 -out cert.pem -clcerts -nokeys
openssl pkcs12 -in cert.p12 -out key.pem -nocerts -nodes


http proxy at 8080 is able to handle https request.
For https request, client will send a CONNECT command to the proxy server, then the proxy server 
forward all tcp traffic between client and destination server.
For http request, proxy server forward request/response between source and destination. 


export http_proxy="http://user:pwd@127.0.0.1:1234"
export https_proxy="http://user:pwd@127.0.0.1:1234"

http_proxy is the proxy server to transfer http requests between source and destination
https_proxy is the proxy server to transfer https requests between source and destination
https_proxy代理不是说代理服务器提供 https://xxx，而是说客户端将所有https的请求交给该代理转发，实际上就是处理CONNECT命令，实现http tunneling. 
所以源代码中的httpsProxy是没有作用的，客户端会首先尝试发送一个http connect命令，会产生错误 tls: first record does not look like a tls handshake 

所以我的理解是，http代理服务器与客户端之间不需要再加一层TLS，那样反而会有问题，这也是为什么iPhone/Android只提供设置http proxy就可以了
因为无需区分http,https请求，但是macos会区分，http发给代理1，https发给代理2
http, https是source和destination之间的事情，与代理服务器无关，不是说https代理服务器需要额外的证书


如果代理服务器需要作为中间人查看https流量，则需要额外的证书对destination返回的数据进行重新签名，并返回给source，当然需要source信任代理的证书，参考charles,mitmproxy等软件

在CONNECT中可以鉴权

Refer to https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT 
CONNECT command is designed for http tunneling, proxy tcp traffic
