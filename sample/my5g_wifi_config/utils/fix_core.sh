#!/bin/bash

# enable forwarding
sudo sysctl -w net.ipv4.ip_forward=1

# stopping ufw service
sudo /etc/init.d/ufw stop

# compiling gtp5g module
cd ~/gtp5g
sudo make && sudo make install

# adding rule in iptables
sudo iptables -t nat -A POSTROUTING -o $(ip route | grep default | cut -d' ' -f5) -j MASQUERADE

