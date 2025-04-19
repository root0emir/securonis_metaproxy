# Securonis Metaproxy

Securonis Metaproxy is a tool that redirects network traffic to a proxy server (e.g., SOCKS5 or HTTP) using iptables rules. It is designed for Linux-based systems and has been tested on Debian.

## Features

- **Transparent Proxying:** Redirects non-proxy-aware applications through a proxy server.
- **Tor Integration:** Can route traffic through the Tor network, even for applications that do not natively support SOCKS proxies.

## Usage

### Setting Up iptables Rules

1. Create a new chain for Metaproxy:
   ```
   $ iptables -t nat -N METAPROXY
   ```

2. Exclude reserved IP ranges:
   ```
   $ iptables -t nat -A METAPROXY -d 10.0.0.0/8 -j RETURN
   $ iptables -t nat -A METAPROXY -d 192.168.0.0/16 -j RETURN
   # Add other ranges as needed...
   ```

3. Redirect traffic to Metaproxy's default port (8675):
   ```
   $ iptables -t nat -A METAPROXY -p tcp -j REDIRECT --to-ports 8675
   ```

4. Start Metaproxy with the configuration file:
   ```
   $ securonis_metaproxy -c securonis_metaproxy.conf
   ```

### Testing the Configuration

Run the following command to test the configuration file:
```
$ securonis_metaproxy -c securonis_metaproxy.conf -t
```

### Using Unix Domain Sockets

Metaproxy supports relaying TCP connections to a SOCKS proxy over Unix domain sockets. For example, if Tor provides a SOCKS proxy via a Unix socket, you can configure Metaproxy as follows:

```
{
    "Relays": [
        {
            "destinationport": "*",
            "relaytype": "SOCKS5",
            "relayip": "unix:/var/lib/tor/socket"
        }
    ]
}
```

## Caveats

1. **Alpha Software:** Metaproxy is still in development and may not be stable or efficient.
2. **Debug Mode Risks:** Avoid using debug mode, as it may expose sensitive connection details.
3. **IPv6 Limitations:** Metaproxy has limited support for IPv6.

## Limitations

1. **iptables Dependency:** Requires iptables rules to redirect traffic.
3. **IPv6 Compatibility:** Not fully compatible with IPv6.

