# SSH HTTP Proxy

## Summary

This utility provides a standalone binary implementing an SSH client 
and HTTP proxy with configurable DNS rewriting.   It reads a single 
configuration file and then:

1. Opens an ssh tunnel to a remote server (without starting a remote shell)
2. Starts an HTTP and HTTPS proxy on the local machine
3. Opens all requests received from the proxy via the ssh tunnel

Additionally the proxy can be configuerd to map specific hostnames to 
arbitrary IP addresses, allowing development servers to be reached using
production hostnames.

## Requirements

Once compiled, the program has no further dependencies; the end user doesn't
need to install ssh clients, or any other software to use the program.

## Use Case

### Simple Access To Dev Services

Suppose:

1. You have a server running a test web service, possibly in a Docker container
which is not (yet) directly accessible from the public Internet.
2. You have a relatively non-technical client who wishes to preview that service
from a Windows or Mac machine

Your options generally are to publish the site on a different domain (eg.
`staging.example.com`) and either restrict access by IP address, or an authentication
mechanism, which may impede full testing of the site.

Using this tool you can leave the site in-situ and provide a single zip file to
client containing an executable and a single configuration file.  The client
may then run the client and configure their browser to use a local proxy and
immediately and securely reach your test service.


### Tunneled Access To The Web

Without leveraging the DNS rewriting, sshproxy can be used in a similar fashion
to using ssh (or putty, on Windows) directly as a SOCKS proxy - All web traffic
directed to the proxy will be routed to the remote server.

This can be useful if you don't want to download & configure putty on Windows.

NOTE: sshproxy performs DNS resolution on the remote server, avoiding "DNS leaks"


## Build Instructions

1. Install Go (https://golang.org/dl/) and setup GOPATH
2. Install Glide (https://github.com/Masterminds/glide)
3. Clone this repo and run `glide install` to downoad the dependencies
4. run `go build .` to build the utility

Or just [download a compiled release instead](http://github.com/gwatts/sshproxy/releases/)

## Configuration

By default the program expects to find a configuration file called `sshproxy.yaml` either
in the current directory, or in the same directory as the program.

Example:

```yaml
proxy-listen: 127.0.0.1:8123
ssh-host: myserver.example.com:22
ssh-user: proxyonly
dns-map:
  '*mytestsite.example.com': 172.17.0.2
  'www.example.com': 172.17.0.2
ssh-key: |
  -----BEGIN RSA PRIVATE KEY-----
  MIIEpQIBAAKCAQEA0WOcrD3qsTGTHTR+f4sa9+IJT/umhXHToXHaY80gdQg4vGpi
  FfO3X1ptfb1ZyxpKalvJmBJno0lhAIK5RxdrW0OILU/r1dAIh8kKSyF+RsFpMMcC
  0YT10yciZ+DACpHXXi4kngVIRutyZXoFclbE7uRNog+KheVhmTBdlpCBGSb+gR/K
  uq6f4ti0tnfdEG6cx98Rlzzvh/Q3F4GMOjJohlB7/upWws+K0vdovctj+00L70b8
  H2f93i4Xi2DvX7BTlviozwpfRPkjgGTyfW63UgjSeblNEpSRbMS4meK+5NwGi/H5
  WT5N4nWbEIGXLQM75zcQbAN/iBKVU8Sujfjf7wIDAQABAoIBAQCO1APq+dE9TTOs
  mEIxfhHHRMhVZrMQE7ToS2FM8n9RVWpeG7MMhlJvGJ/XRXIauLRKGJJKyUMofsVM
  M99uPutcNZSOVBXqox0uglQjK5WXbhbyzs19XdTRU8CEWyqkCxd9hrwzjibfOXuu
  /Kz6cXWj/td11GQJMY6BkBxGuAtXu6XBAwQrY2kAuPT5OztewYqvDw1OGj2JH9tr
  F0vi0rVYjkWJjvohON5v+hy8S207CktJLmdOUUc3rYpPc8A/0XwpVCq2nPAvfyMN
  ptx1Ybc9xB7jgaI/0Ovp37Ui0UI8kkT3gkvmgKi0gRhVj9BaK4wkNd33BUGsFnj9
  i2EphQsBAoGBAOk/mt9Vcpyovh4Tm3a+5QTqan4UCFEm0kCfPbUHSEO82So9hjTd
  LIimz4tZYTPwzJeUrUNN2whdCL2fbeVo26mrZXnaXWcuauoaKz3c1PMFsAq4HLSk
  ALLuHdqBE30GqWCjFb+TBWCzhVSKRSD1WY8fwoEQu1GRhoRli8EZeEmNAoGBAOXQ
  OHmDyzJ7VF8B3/JbFYrCXspSDSc/bfxFZcsjWLAs8zUIMPwOBrLGT72PzFriaylL
  D3UI7z2nU4TRwXi3391xb+J6lmPzd5QI5Px7mDp/URpbuMhvYlMSSN23SLkVmuiO
  efv3/pDjM0V6dibQJm6RM61V8wKd5jIq2mDIZSprAoGAGnuXSP45qiHanC2bvCrG
  c/1to+0AWL5wptetuO0fvlklyw77OutV0BoofGjkiXIwuJEv7vFbCiMOCAGfB6oV
  LrmAJwqtCjcR+oyIFlkJcKJXr4/h6nyoe6hfiVyYatyjxI4fvQWjWaxoWgXs/WX+
  CisP+Xl92zALtuKUsJMEvk0CgYEAqHxv6ybk4q3ovX7yYQzGTmUSeeKOIiguyrVW
  XAgeDYvnAwpuX10pLAiYjbHPcRJu3mdZfcR/IgR7BvWBkq+8QO3ZyYF2oPDuyml6
  +GDkyn5tR5XXc5u1ypGtOmAVwRxF5hoO9Nxslmz8OgP+e5Y/lvB9oqdQ8qoxCrbA
  RBSnluECgYEA5Q9Wv9yZUoAvFwnw+GmF0JbyYG45YtNBylJfn5LZ1SClL14O4dUD
  637tB2OrcZwOunYLE8d3gi8cNWBv8TC9nLHyZWJinYXOetgmQYfCodW3NhSDoXXt
  pIqeUyeFghhjtwLtbxC6R5iKjhB5PzwGhCAbRAseIr8Pd/TFwOw3psY=
  -----END RSA PRIVATE KEY-----
```

* **proxy-listen** - Ths address for the HTTP proxy to listen on (at which Firefox must be directed) - defaults to 127.0.0.1:8123
* **ssh-host** - The address of the SSH server to connect to; must include the port number
* **ssh-user** - The username on the SSH server to connect to; the user does not need to have shell access, but must have port-forwarding privileges
* **host-sig** - The expected signature of the remote server - sshproxy will check that the signature presented by the server matches this
to prevent man-in-the-middle attacks.  If it's omitted, sshproxy will prompt on first connection and then save it in the configuraiton file
* **dns-map** - A map of hostname patterns to map to an alternative IP address.  The patterns may include asterisk and period operators per 
the Match function documented here https://golang.org/pkg/path/#Match
* **ssh-key** - The ssh key to use to connect to the SSH server.  Must not be encrypted with a passphrase.  **NOTE** Each line of the key
file must be indented.

## Configuring The Server

Assuming your server is running Linux, you can add a dedicated user for sshproxy to connect to
or add an entry to the authorized_keys file for an existing user.

Steps to create a new user, with no password, no shell access with port-forwarding
abilities only:

### Create a new user

Create a user with no password.  Ideally they woudln't have a functioning shell either,
but this'll do for a start.

```bash
$ adduser --disable-password proxyonly
```

### Create an ssh key

The user's ssh key must not have a passphrase.  It can be generated from any machine; the
private portion will go into sshproxy.yaml, the public portion will go in the server's
authorized_keys file.

```bash
$ ssh-keygen -N "" -f proxykey
```

This will generate two files - `proxykey` (the private portion to be included in
`sshproxy.yaml`) and `proxykey.pub` (to be used in the `authorized_keys` file)

### Create the authorized_keys file

Create `~proxyonly/.ssh/authorized_keys` on the server and include
the contents of proxykey.pub generated above:

```
no-agent-forwarding,no-pty,no-user-rc,no-X11-forwarding,command="/bin/false" ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDRY5ysPeqxMZMdNH5/ixr34glP+6aFcdOhcdpjzSB1CDi8amIV87dfWm19vVnLGkpqW8mYEmejSWEAgrlHF2tbQ4gtT+vV0AiHyQpLIX5GwWkwxwLRhPXTJyJn4MAKkddeLiSeBUhG63JlegVyVsTu5E2iD4qF5WGZMF2WkIEZJv6BH8q6rp/i2LS2d90QbpzH3xGXPO+H9DcXgYw6MmiGUHv+6lbCz4rS92i9y2P7TQvvRvwfZ/3eLheLYO9fsFOW+KjPCl9E+SOAZPJ9brdSCNJ5uU0SlJFsxLiZ4r7k3AaL8flZPk3idZsQgZctAzvnNxBsA3+IEpVTxK6N+N/v user@example.com
```

Make sure the .ssh and authorized_keys file are owned by the proxyuser user `chown -R proxyuser:proxyuser ~proxyuser/.ssh`



## Running the program

Give the client a compiled copy of sshproxy and a suitable `sshproxy.yaml` file.

On Mac they will have to right click on the program and select "Open" and accept
the warning to run the program to avoid GateKeeper restrictions on unsigned programs.

On Windows you may need to run the program from a command prompt for similar reasons.

The client can then configure their web browser to use the running proxy.  If they just
open http://localhost:8123 they will be shown a message with links to web sites telling
them how to configure their browser if they're unfamiliar with the process.
