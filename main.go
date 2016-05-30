// Copyright (c) 2016, Gareth Watts
// All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Bowery/prompt"
	"github.com/gwatts/goproxy"
	"golang.org/x/crypto/ssh"
)

var (
	cfgFilename = flag.String("config", "", "Filename for configuration file.  Defaults to sshproxy.yml")
	noHostCheck = flag.Bool("no-host-check", false, "Disables checking of the ssh host key")
	saveHostSig = flag.Bool("save-host-sig", false, "Saves the ssh host signature on first connection without prompting")
)

const (
	defaultProxyListen    = "127.0.0.1:8123"
	defaultConfigFilename = "sshproxy.yaml"
)

func main() {
	flag.Parse()
	log.Println("sshproxy version", version, "https://github.com/gwatts/sshproxy")
	cfg := loadConfig()

	key, err := ssh.ParsePrivateKey([]byte(cfg.SSHKey))
	if err != nil {
		log.Fatal("Failed to parse ssh key:", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: cfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			if *noHostCheck {
				log.Println("Warning: skipped ssh host signature validation")
				return nil
			}
			serverSig := hostSigFromKey(key)
			if cfgSig := cfg.hostSig(); cfgSig == nil {
				if *saveHostSig {
					cfg.updateHostSig(serverSig)
				} else {
					fmt.Println("")
					fmt.Printf("No host-sig entry for server %q found in configuration file.\n", hostname)
					fmt.Println("Remote server address is", remote)
					fmt.Println("Remote server signature is", serverSig)
					fmt.Println("")
					if ok, err := prompt.Ask("Accept and save this signature to the configuration file"); err != nil {
						log.Fatal("Error prompting user", err)
					} else if ok {
						cfg.updateHostSig(serverSig)
					} else {
						log.Fatal("Host key rejected by user")
					}
				}
			} else if !serverSig.isEqual(cfgSig) {
				log.Println("Host key does not match configured host-sig signature!")
				log.Println("Host signature", serverSig)
				log.Println("Config signature", cfgSig)
				log.Fatal("Abandoning connection")
			}
			log.Println("SSH server has correct signature", serverSig)
			return nil
		},
	}

	if sig := cfg.hostSig(); sig != nil {
		sshConfig.HostKeyAlgorithms = []string{sig.sigType}
	}

	sshClient, err := ssh.Dial("tcp", cfg.SSHHost, sshConfig)
	if err != nil {
		log.Fatalf("Failed to connect to SSH server at %q with username %q: %v",
			cfg.SSHHost, cfg.SSHUser, err)
	}
	log.Printf("Connected to remote SSH server at %q with username %q",
		cfg.SSHHost, cfg.SSHUser)

	proxy := goproxy.NewProxyHttpServer()
	proxy.NonproxyHandler = http.HandlerFunc(nonProxy)
	proxy.Tr.Dial = func(network, addr string) (c net.Conn, err error) {
		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			if newHost := cfg.mapHost(host); newHost != "" {
				addr = newHost + ":" + port
			}
		}
		c, err = sshClient.Dial(network, addr)
		log.Printf("DIAL host=%s  target=%s error=%v", host, addr, err)
		return c, err
	}

	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if r != nil && r.Request != nil {
			log.Printf("[%d] REQUEST client=%s response=%d  %s %s",
				ctx.Session, r.Request.RemoteAddr, r.StatusCode, r.Request.Method, r.Request.URL)
		}
		return r
	})

	log.Println("Starting proxy on", cfg.ProxyListen)
	log.Fatal(http.ListenAndServe(cfg.ProxyListen, proxy))
}
