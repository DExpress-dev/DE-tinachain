package host

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"time"

	"github.com/kataras/iris/core/netutil"
)

// ProxyHandler returns a new ReverseProxy that rewrites
// URLs to the scheme, host, and base path provided in target. If the
// target's path is "/base" and the incoming request was for "/dir",
// the target request will be for /base/dir.
//
// Relative to httputil.NewSingleHostReverseProxy with some additions.
//
// Look `ProxyHandlerRemote` too.
func ProxyHandler(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		modifyProxiedRequest(req, target)
		req.Host = target.Host
		req.URL.Path = path.Join(target.Path, req.URL.Path)
	}

	p := &httputil.ReverseProxy{Director: director}

	if netutil.IsLoopbackHost(target.Host) {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // lint:ignore
		}
		p.Transport = transport
	}

	return p
}

func modifyProxiedRequest(req *http.Request, target *url.URL) {
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host

	if target.RawQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
	}

	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "")
	}
}

// ProxyHandlerRemote returns a new ReverseProxy that rewrites
// URLs to the scheme, host, and path provided in target.
// Case 1: req.Host == target.Host
// behavior same as ProxyHandler
// Case 2: req.Host != target.Host
// the target request will be forwarded to the target's url
// insecureSkipVerify indicates enable ssl certificate verification or not.
//
// Look `ProxyHandler` too.
func ProxyHandlerRemote(target *url.URL, insecureSkipVerify bool) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		modifyProxiedRequest(req, target)

		if req.Host != target.Host {
			req.URL.Path = target.Path
		} else {
			req.URL.Path = path.Join(target.Path, req.URL.Path)
		}

		req.Host = target.Host
	}
	p := &httputil.ReverseProxy{Director: director}

	if netutil.IsLoopbackHost(target.Host) {
		insecureSkipVerify = true
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify}, // lint:ignore
	}
	p.Transport = transport
	return p
}

// NewProxy returns a new host (server supervisor) which
// proxies all requests to the target.
// It uses the httputil.NewSingleHostReverseProxy.
//
// Usage:
// target, _ := url.Parse("https://mydomain.com")
// proxy := NewProxy("mydomain.com:80", target)
// proxy.ListenAndServe() // use of `proxy.Shutdown` to close the proxy server.
func NewProxy(hostAddr string, target *url.URL) *Supervisor {
	proxyHandler := ProxyHandler(target)
	proxy := New(&http.Server{
		Addr:    hostAddr,
		Handler: proxyHandler,
	})

	return proxy
}

// NewProxyRemote returns a new host (server supervisor) which
// proxies all requests to the target.
// It uses the httputil.NewSingleHostReverseProxy.
//
// Usage:
// target, _ := url.Parse("https://anotherdomain.com/abc")
// proxy := NewProxyRemote("mydomain.com", target, false)
// proxy.ListenAndServe() // use of `proxy.Shutdown` to close the proxy server.
func NewProxyRemote(hostAddr string, target *url.URL, insecureSkipVerify bool) *Supervisor {
	proxyHandler := ProxyHandlerRemote(target, insecureSkipVerify)
	proxy := New(&http.Server{
		Addr:    hostAddr,
		Handler: proxyHandler,
	})

	return proxy
}

// NewRedirection returns a new host (server supervisor) which
// redirects all requests to the target.
// Usage:
// target, _ := url.Parse("https://mydomain.com")
// r := NewRedirection(":80", target, 307)
// r.ListenAndServe() // use of `r.Shutdown` to close this server.
func NewRedirection(hostAddr string, target *url.URL, redirectStatus int) *Supervisor {
	redirectSrv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		Addr:         hostAddr,
		Handler:      RedirectHandler(target, redirectStatus),
	}

	return New(redirectSrv)
}

// RedirectHandler returns a simple redirect handler.
// See `NewProxy` or `ProxyHandler` for more features.
func RedirectHandler(target *url.URL, redirectStatus int) http.Handler {
	targetURI := target.String()
	if redirectStatus <= 300 {
		// here we should use StatusPermanentRedirect but
		// that may result on unexpected behavior
		// for end-developers who might change their minds
		// after a while, so keep status temporary.
		// Note thatwe could also use StatusFound
		// as we do on the `Context#Redirect`.
		// It will also help us to prevent any post data issues.
		redirectStatus = http.StatusTemporaryRedirect
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectTo := path.Join(targetURI, r.URL.Path)
		if len(r.URL.RawQuery) > 0 {
			redirectTo += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, redirectTo, redirectStatus)
	})
}
