#!/bin/bash

# Setup the network interfaces and namespaces required for ue
ip link set lo up
ip link add ipsec0 type vti 172.16.0.2 remote 172.16.0.1 key 5
ip link set dev ipsec0 up


#UENS="UEns"
#EXEC_UENS="ip netns exec ${UENS}"
#ip link add veth1 type veth peer name veth2
#ip netns add ${UENS}
#ip link add veth2 type veth peer name veth3
#ip addr add 192.168.127.1/24 dev veth2
#ip link set veth2 up
#ip link set veth3 netns ${UENS}
#${EXEC_UENS} ip addr add 192.168.127.2/24 dev veth3
#${EXEC_UENS} ip link set lo up
#${EXEC_UENS} ip link set veth3 up
#${EXEC_UENS} ip link add ipsec0 type vti local 192.168.127.2 remote 192.168.127.1 key 5
#${EXEC_UENS} ip link set ipsec0 up

#sudo ip link add name ipsec0 type vti local 192.168.127.1 remote 0.0.0.0 key 5
#sudo ip addr add 10.0.0.1/24 dev ipsec0    
