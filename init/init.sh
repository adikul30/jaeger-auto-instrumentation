#!/bin/bash
set -x
# Forward TCP traffic on port 80 to port 8000 on the eth0 interface.
# iptables -t nat -A PREROUTING -p tcp -i eth0 --dport 80 -j REDIRECT --to-port 8000


#iptables -t nat -P PREROUTING ACCEPT
#iptables -t nat -P INPUT ACCEPT
#iptables -t nat -P OUTPUT ACCEPT
#iptables -t nat -P POSTROUTING ACCEPT
iptables -t nat -N ISTIO_INBOUND
iptables -t nat -N ISTIO_IN_REDIRECT
iptables -t nat -N ISTIO_OUTPUT
iptables -t nat -N ISTIO_REDIRECT
iptables -t nat -A PREROUTING -p tcp -i eth0 -j ISTIO_INBOUND
iptables -t nat -A OUTPUT -p tcp -j ISTIO_OUTPUT
iptables -t nat -A ISTIO_INBOUND -p tcp -i eth0 --dport 80 -j ISTIO_IN_REDIRECT
iptables -t nat -A ISTIO_IN_REDIRECT -p tcp -i eth0 -j REDIRECT --to-port 8000
iptables -t nat -A ISTIO_OUTPUT ! -d 127.0.0.1/32 -o lo -j ISTIO_REDIRECT
iptables -t nat -A ISTIO_OUTPUT -m owner --uid-owner 1337 -j RETURN
iptables -t nat -A ISTIO_OUTPUT -m owner --gid-owner 1337 -j RETURN
iptables -t nat -A ISTIO_OUTPUT -d 127.0.0.1/32 -j RETURN
iptables -t nat -A ISTIO_OUTPUT -j ISTIO_REDIRECT
iptables -t nat -A ISTIO_REDIRECT -p tcp -j REDIRECT --to-port 8000

# List all iptables rules.
iptables -t nat --list
