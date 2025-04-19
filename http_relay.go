package main

import (
    "bufio"
    "fmt"
    "net"
    "net/http"
    "net/url"
    "strconv"
)

func relayHttp(clientConn net.Conn, client *Client, relay Relay) error {
    proxyAddr := net.JoinHostPort(relay.RelayIP, relay.RelayPort)

    destAddr := net.JoinHostPort(client.destAddr.String(), strconv.Itoa(int(client.destPort)))
    destIP, err := resolveThroughProxy(client.destAddr.String(), client.destPort) // Proxy üzerinden çözümleme
    if err != nil {
        return fmt.Errorf("Failed to resolve destination address %s: %s", client.destAddr.String(), err)
    }
    resolvedDestAddr := net.JoinHostPort(destIP.String(), strconv.Itoa(int(client.destPort)))

    destUrl, err := url.Parse(resolvedDestAddr)
    if err != nil {
        return fmt.Errorf("Invalid destAddr %s: %s", resolvedDestAddr, err)
    }

    proxyConn, err := net.Dial("tcp", proxyAddr)
    if err != nil {
        return fmt.Errorf("Error dialing HTTP proxy %s: %s", proxyAddr, err)
    }

    connectReq := &http.Request{
        Method: "CONNECT",
        URL:    destUrl,
        Host:   resolvedDestAddr,
        Header: make(http.Header),
    }
    connectReq.Write(proxyConn)
    br := bufio.NewReader(proxyConn)
    resp, err := http.ReadResponse(br, connectReq)
    if err != nil {
        proxyConn.Close()
        return err
    }
    if resp.StatusCode != 200 {
        proxyConn.Close()
        return fmt.Errorf("Proxy returned non-200 status: %d", resp.StatusCode)
    }

    go copyAndClose(proxyConn, clientConn)
    go copyAndClose(clientConn, proxyConn)
    return nil
}


