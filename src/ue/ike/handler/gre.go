package handler

import (
	"fmt"
	"free5gc/src/ue/logger"
	"free5gc/src/ue/context"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
)

var greLog *logrus.Entry

func init() {
	greLog = logger.GRELog
}

func GetGREInfo(name string ) netlink.Link{

	links, err := netlink.LinkList()
	if err != nil {
		greLog.Fatal(err)
	}
	var linkGRE netlink.Link

	for _, link := range links {
		if link.Attrs() != nil {
			linkGRE = link
			break
		}
	}
	if linkGRE == nil {
		greLog.Fatal("GRE " + name + " interface not found in OS.")
	}

	return linkGRE
}

func SetupNewGRETunnelWithN3IWF(ctx context.UeN3iwf) (*netlink.Gretun, error){

	newGRETunnel := &netlink.Gretun{
		LinkAttrs: netlink.LinkAttrs{
			Name: ctx.GREName,
			MTU: 1400,
		},
		Local: net.IP(ctx.GREBindAddress),
		Remote: net.IP(ctx.IPSecGatewayAddress),
	}

	greLog.Infoln(fmt.Sprintf( "Add Gre Tun Interface -> %s", ctx.GREName ) )
	if err := netlink.LinkAdd(newGRETunnel); err != nil {
		greLog.Fatal(err)
		return newGRETunnel, err
	}

	ctx.GRETunnel = newGRETunnel
	return newGRETunnel, nil

}

func SetGREIpAddr(ipAddr net.IP, mask net.IPMask, link netlink.Link){
	linkGREAddr := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   ipAddr,
			Mask: mask,
		},
	}
	if err := netlink.AddrAdd(link, linkGREAddr); err != nil {
		greLog.Fatal(err)
	}
	//Set GRE interface
	if err := netlink.LinkSetUp(link); err != nil {
		greLog.Fatal(err)
	}
}

func AddGRERoute(GRETunnel netlink.Link, dstIP net.IP, netmask net.IPMask){

	upRoute := &netlink.Route{
		LinkIndex: GRETunnel.Attrs().Index,
			Dst: &net.IPNet{
				IP:  dstIP,
				Mask: netmask,
			},
		}

	if err := netlink.RouteAdd(upRoute); err != nil {
		greLog.Fatal(err)
	}
}