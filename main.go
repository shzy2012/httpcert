package main

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

/*https://mlog.club/article/244
HostPolicy 是要注册的域名，可填写多个域名：autocert.HostWhitelist(domain, domain2, domain3)，如果留空则是任何解析向该服务器ip 的域名。
Cache 是存放证书的目录
RenewBefore 是指定更新证书的时间，如果不填则是在过期前30天自动更新
*/
const (
	contactEmail = "shzy2012@gmail.com"
	domain       = "example.com"
)

func main() {

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain), //your domain here
		Cache:      autocert.DirCache("certs"),     //folder for storing certificates
		Email:      contactEmail,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world ssl!"))
	})

	go http.ListenAndServe(":http", certManager.HTTPHandler(nil)) // 支持 http-01
	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			NextProtos:     []string{http2.NextProtoTLS, "http/1.1"},
			MinVersion:     tls.VersionTLS12,
		},
		MaxHeaderBytes: 32 << 20,
	}

	server.ListenAndServeTLS("", "") //key and cert are comming from Let's Encrypt
}
