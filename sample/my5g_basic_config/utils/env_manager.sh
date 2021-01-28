#!/bin/bash

if [[ $# -ne 2 ]] || ([[ $1 != "up" ]] && [[ $1 != "down" ]]); then
  echo "Usage: $0 [up|down] [internet_iface]"
  exit 1
fi

HOSTNAMES=(
  "10.1.1.2 amf"
  "10.1.1.3 smf"
  "10.1.1.4 ausf"
  "10.1.1.5 nssf"
  "10.1.1.6 pcf"
  "10.1.1.7 udm"
  "10.1.1.8 udr"
  "10.1.1.9 upf"
  "10.1.1.11 db"
  "10.1.1.10 nrf"
)

IFACE=$2

if [[ $1 == "up" ]]; then
  echo "Creating network interfaces and namespaces..."
  # create network interfaces and add ip addresses
  # 5gc network (it's not needed but helps to organize/separate the networks)
  
  #to create bridge of core: br-5gc
  sudo ip link add br-5gc type  bridge # bridge for core

  # ip address of components
  sudo ip addr add 10.1.1.2/24  dev br-5gc # amf
  sudo ip addr add 10.1.1.3/24  dev br-5gc # smf
  sudo ip addr add 10.1.1.4/24  dev br-5gc # ausf
  sudo ip addr add 10.1.1.5/24  dev br-5gc # nssf
  sudo ip addr add 10.1.1.6/24  dev br-5gc # pcf
  sudo ip addr add 10.1.1.7/24  dev br-5gc # udm
  sudo ip addr add 10.1.1.8/24  dev br-5gc # udr
  sudo ip addr add 10.1.1.9/24  dev br-5gc # upf
  sudo ip addr add 10.1.1.10/24 dev br-5gc # nrf
  sudo ip addr add 10.1.1.11/24 dev br-5gc # mongodb

  # to up bridge br-5gc
  sudo ip link set br-5gc up 

  # Inteface added to handle N2 interface (it's not needed but helps to organize/separate the networks)
  
  # to create bridge to handle n2 interface
  sudo ip link add br-n2 type bridge

  # ip address of components at n2 interface
  sudo ip addr add 172.16.0.1/24 dev br-n2
  sudo ip addr add 172.16.0.2/24 dev br-n2

  # Setup network namespace for UPF
  sudo ip netns add UPFns

  sudo ip link add veth0 type veth peer name veth1
  sudo ip link set veth0 up
  sudo ip addr add 60.60.0.1 dev lo
  sudo ip addr add 10.200.200.1/24 dev veth0
  sudo ip addr add 10.200.200.2/24 dev veth0

  sudo ip link set veth1 netns UPFns
  sudo ip netns exec UPFns ip link set lo up
  sudo ip netns exec UPFns ip link set veth1 up
  sudo ip netns exec UPFns ip addr add 60.60.0.101 dev lo
  sudo ip netns exec UPFns ip addr add 10.200.200.101/24 dev veth1
  sudo ip netns exec UPFns ip addr add 10.200.200.102/24 dev veth1

  #veth2-veth3 beetween AP and N3IWF
  sudo ip link add dev veth2 type veth peer name veth3

  # veth2- ip address:192.168.127.1/24
  sudo ip addr add 192.168.127.1/24 dev veth2  
  sudo ip link set veth2 up

  # veth3 in APns - ip address: 192.168.127.2/24
  sudo ip link set veth3 netns APns
  sudo ip netns exec APns ip addr add 192.168.127.2/24 dev veth3
  sudo ip netns exec APns ip link set veth3 up
  
  # default route for wlan1
 #sudo ip netns exec UEns route add default gw 192.168.1.10 wlan1 #wlan 1 - interface cliente
  sudo ip netns exec UEns route add -net 192.168.127.0/24 gw 192.168.1.10 dev wlan1

  # static route 
  sudo ip route add 192.168.1.0/24 via 192.168.127.2 dev veth2 

  # settings for ipsec0
  sudo ip netns exec UEns ip link set lo up
  sudo ip netns exec UEns ip link add ipsec0 type vti local 192.168.1.1 remote 192.168.127.1 key 5
  sudo ip netns exec UEns ip link set ipsec0 up

  sudo ip link add name ipsec0 type vti local 192.168.127.1 remote 0.0.0.0  key 5
  sudo ip addr add 10.0.0.1/24 dev ipsec0
  sudo ip link set ipsec0 up

  sudo ip link add veth4 type veth peer name veth5
  sudo ip addr add 10.1.2.1/24 dev veth4
  sudo ip link set veth4 up

  sudo ip link set veth5 netns UPFns
  sudo ip netns exec UPFns ip addr add 10.1.2.2/24 dev veth5
  sudo ip netns exec UPFns ip link set veth5 up
  sudo ip netns exec UPFns ip route add default via 10.1.2.1

  sudo ip netns exec UPFns iptables -t nat -A POSTROUTING -o veth5 -j MASQUERADE
  sudo iptables -t nat -A POSTROUTING -s 10.1.2.2/24 -o ${IFACE} -j MASQUERADE
  sudo iptables -A FORWARD -i ${IFACE} -o veth4 -j ACCEPT
  sudo iptables -A FORWARD -o ${IFACE} -i veth4 -j ACCEPT

  echo "Network interfaces and namespaces created."

elif [[ $1 == "down" ]]; then
  echo "Removing network interfaces and namespaces.."
  sudo ip link set br-5gc down
  sudo ip link delete br-5gc
  sudo ip link delete br-n2
  sudo ip xfrm policy flush
  sudo ip xfrm state flush
  sudo ip link del veth2
  sudo ip link del veth4
  sudo ip link del ipsec0
  sudo ip link del veth0
  sudo ip netns exec UEns ip link del ipsec0
  sudo ip netns del UEns
  sudo ip netns del UPFns
  sudo ip netns del APns

  sudo rm /dev/mqueue/*
  for host in "${HOSTNAMES[@]}"; do
    sudo sed -i "/$host/d" /etc/hosts
  done
  echo "Network interfaces and namespaces removed."
fi
