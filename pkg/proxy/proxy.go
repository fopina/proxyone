package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/fopina/proxyone/pkg/config"
)

type Proxy struct {
	config *config.Config
}

func NewProxy(cfg *config.Config) *Proxy {
	return &Proxy{config: cfg}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
		return
	}

	p.handleHTTP(w, r)
}

func (p *Proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Host
	if !p.config.ShouldUseProxy(host) {
		// Direct connection for HTTPS
		conn, err := net.Dial("tcp", host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer conn.Close()

		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
			return
		}

		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer clientConn.Close()

		// Send 200 Connection established
		clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

		// Start tunneling
		go io.Copy(conn, clientConn)
		io.Copy(clientConn, conn)
		return
	}

	// Use upstream proxy for HTTPS
	p.handleConnectViaProxy(w, r)
}

func (p *Proxy) handleConnectViaProxy(w http.ResponseWriter, r *http.Request) {
	upstreamConn, err := net.Dial("tcp", p.config.UpstreamProxy.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer upstreamConn.Close()

	// Send CONNECT request to upstream proxy
	connectReq := "CONNECT " + r.URL.Host + " HTTP/1.1\r\n" +
		"Host: " + r.URL.Host + "\r\n" +
		"Proxy-Connection: keep-alive\r\n\r\n"
	_, err = upstreamConn.Write([]byte(connectReq))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Read response from upstream proxy
	buf := make([]byte, 1024)
	n, err := upstreamConn.Read(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if !strings.Contains(string(buf[:n]), "200") {
		http.Error(w, "Upstream proxy connection failed", http.StatusServiceUnavailable)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	// Send 200 Connection established
	clientConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

	// Start tunneling
	go io.Copy(upstreamConn, clientConn)
	io.Copy(clientConn, upstreamConn)
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Host
	if host == "" {
		host = r.Host
	}

	if !p.config.ShouldUseProxy(host) {
		// Direct request
		resp, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Copy headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
		return
	}

	// Use upstream proxy
	p.handleHTTPViaProxy(w, r)
}

func (p *Proxy) handleHTTPViaProxy(w http.ResponseWriter, r *http.Request) {
	// Create a new request for the upstream proxy
	upstreamURL := &url.URL{
		Scheme: p.config.UpstreamProxy.Scheme,
		Host:   p.config.UpstreamProxy.Host,
		Path:   r.URL.String(), // Full URL as path for proxy
	}

	req, err := http.NewRequest(r.Method, upstreamURL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add proxy authorization if needed
	if p.config.UpstreamProxy.User != nil {
		password, _ := p.config.UpstreamProxy.User.Password()
		req.SetBasicAuth(p.config.UpstreamProxy.User.Username(), password)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func (p *Proxy) Start() error {
	log.Printf("Starting proxy server on %s", p.config.ListenAddr)
	return http.ListenAndServe(p.config.ListenAddr, p)
}
