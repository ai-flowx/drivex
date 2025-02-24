package config

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

const (
	ProxyTypeHTTP   = "http"
	ProxyTypeSOCKS5 = "socks5"
)

func getHttpProxyTransport(proxyURL string, timeout int) (*http.Transport, error) {
	parsedProxyURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing proxy URL %s: %v", proxyURL, err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedProxyURL),
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: time.Duration(timeout) * time.Second,
		}).DialContext,
	}

	return transport, nil
}

func getSocks5Transport(proxyAddr string, timeout int) (*http.Transport, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("error creating SOCKS5 proxy at %s: %v", proxyAddr, err)
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dialer.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()
			return dialContext(ctx, network, addr)
		},
	}

	return transport, nil
}

func GetTypeProxyTransport(proxyType, proxyAddr string, timeout int) (*http.Transport, error) {
	switch proxyType {
	case ProxyTypeHTTP:
		return getHttpProxyTransport(proxyAddr, timeout)
	case ProxyTypeSOCKS5:
		return getSocks5Transport(proxyAddr, timeout)
	default:
		return nil, errors.New("unsupported proxy type: " + proxyType)
	}
}

func GetConfProxyTransport() (proxyType, proxyAddr string, transport *http.Transport, err error) {
	proxyType = strings.ToLower(GProxyConf.Type)

	timeout := GProxyConf.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	switch proxyType {
	case ProxyTypeHTTP:
		proxyAddr = GProxyConf.HTTPProxy
		transport, err = getHttpProxyTransport(proxyAddr, timeout)
	case ProxyTypeSOCKS5:
		if len(GProxyConf.Socks5Proxy) >= 7 && GProxyConf.Socks5Proxy[:7] == "socks5:" {
			proxyURL, err := url.Parse(GProxyConf.Socks5Proxy)
			if err != nil {
				return "", "", nil, fmt.Errorf("error parsing proxy URL: %v\n", err)
			}
			proxyAddr = proxyURL.Host
		} else {
			proxyAddr = GProxyConf.Socks5Proxy
		}
		transport, err = getSocks5Transport(proxyAddr, timeout)
	default:
		return "", "", nil, errors.New("unsupported proxy type: " + proxyType)
	}

	return proxyType, proxyAddr, transport, err
}
