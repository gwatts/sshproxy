// Copyright (c) 2016, Gareth Watts
// All rights reserved.

package main

import (
	"fmt"
	"net/http"
)

var nonProxyResponse = `
<html>
<head><title>sshproxy</title></head>
<body>
	<p>
		This is a <a href="https://github.com/gwatts/sshproxy">sshproxy web proxy</a>.
	</p>
	<p>
		To use it you must configure your web browser's HTTP and HTTPS proxy settings
		to point to this service. Help on how to do that can be found for:
	</p>
	<ul>
	<li><a href="http://customers.trustedproxies.com/knowledgebase.php?action=displayarticle&id=37">Firefox</a></li>
	<li><a href="http://customers.trustedproxies.com/knowledgebase.php?action=displayarticle&id=10">Google Chrome</a></li>
	<li><a href="http://customers.trustedproxies.com/knowledgebase.php?action=displayarticle&id=38">Internet Explorer</a></li>
	<li><a href="http://customers.trustedproxies.com/knowledgebase.php?action=displayarticle&id=15">Safari (Mac)</a></li>
	<li><a href="http://customers.trustedproxies.com/knowledgebase.php?action=displayarticle&id=45">Safari (Windows)</a></li>
	</ul>
</body>
</html>
`

func nonProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, nonProxyResponse)
}
