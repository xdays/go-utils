# v2ray-sdk

V2ray Go SDK

# Features

* Add inbond user
* Remove inbond user
* Get Stats by inbond or user


# Example

Setup v2ray

```
v2ray -config config/server_api.json
v2ray -config config/client_api.json
```

Server should reject client connection since it doesn't have client uuid

```
curl -Iv -x 127.0.0.1:10084 https://ip.cn
*   Trying 127.0.0.1:10084...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 10084 (#0)
* allocate connect buffer!
* Establish HTTP proxy tunnel to ip.cn:443
> CONNECT ip.cn:443 HTTP/1.1
> Host: ip.cn:443
> User-Agent: curl/7.65.1
> Proxy-Connection: Keep-Alive
>
< HTTP/1.1 200 Connection established
HTTP/1.1 200 Connection established
<

* Proxy replied 200 to CONNECT request
* CONNECT phase completed!
* ALPN, offering http/1.1
* CONNECT phase completed!
* CONNECT phase completed!
* Server aborted the SSL handshake
* Closing connection 0
curl: (35) Server aborted the SSL handshake
```

Add a new user(uuid)

```
go run example/main.go -i a994b3c1-c7cc-4868-8072-c93e491bba0b -a
```

Now test again

```
curl -Iv -x 127.0.0.1:10084 https://ip.cn
*   Trying 127.0.0.1:10084...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 10084 (#0)
* allocate connect buffer!
* Establish HTTP proxy tunnel to ip.cn:443
> CONNECT ip.cn:443 HTTP/1.1
> Host: ip.cn:443
> User-Agent: curl/7.65.1
> Proxy-Connection: Keep-Alive
>
< HTTP/1.1 200 Connection established
HTTP/1.1 200 Connection established
<

* Proxy replied 200 to CONNECT request
* CONNECT phase completed!
* ALPN, offering http/1.1
* CONNECT phase completed!
* CONNECT phase completed!
* TLS 1.2 connection using TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
* Server certificate: sni.cloudflaressl.com
* Server certificate: CloudFlare Inc ECC CA-2
* Server certificate: Baltimore CyberTrust Root
> HEAD / HTTP/1.1
> Host: ip.cn
> User-Agent: curl/7.65.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
HTTP/1.1 200 OK
< Date: Sun, 18 Aug 2019 17:23:11 GMT
Date: Sun, 18 Aug 2019 17:23:11 GMT
< Content-Type: application/json; charset=UTF-8
Content-Type: application/json; charset=UTF-8
< Connection: keep-alive
Connection: keep-alive
< Set-Cookie: __cfduid=d4b3cfa6dac0e15c1ddbaaf8e751084ec1566148991; expires=Mon, 17-Aug-20 17:23:11 GMT; path=/; domain=.ip.cn; HttpOnly
Set-Cookie: __cfduid=d4b3cfa6dac0e15c1ddbaaf8e751084ec1566148991; expires=Mon, 17-Aug-20 17:23:11 GMT; path=/; domain=.ip.cn; HttpOnly
< Expect-CT: max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"
Expect-CT: max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"
< Server: cloudflare
Server: cloudflare
< CF-RAY: 508592fcbbd97a80-LAX
CF-RAY: 508592fcbbd97a80-LAX

<
* Connection #0 to host 127.0.0.1 left intact
```

It works which indicates our new user is added into server successfully.

Besides, you can get user traffic usage

```
go run example/main.go -i a994b3c1-c7cc-4868-8072-c93e491bba0b -u
stat: <
  name: "user>>>a994b3c1-c7cc-4868-8072-c93e491bba0b@x.xx>>>traffic>>>downlink"
  value: 2954
>
 <nil>
```
