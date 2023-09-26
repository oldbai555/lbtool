package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	InnerIp            = "$inner_ip"
	NetInterfacePrefix = "$iface"
)

func GetMacAddrs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var macAddrs []string
	for _, netInterface := range interfaces {
		if macAddr := netInterface.HardwareAddr.String(); macAddr != "" {
			macAddrs = append(macAddrs, macAddr)
		}
	}
	return macAddrs, nil
}

func GetLocalIpsWithoutLoopback() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && IsLocalIpV4(ipNet.IP) {
			ips = append(ips, ipNet.IP.String())
		}
	}
	return ips, nil
}

func EvalVarToParseIp(ipVar string) (string, error) {
	if ipVar == InnerIp {
		innerIps, err := GetLocalIpsWithoutLoopback()
		if err != nil {
			return "", err
		}
		if len(innerIps) == 0 {
			return "", errors.New("not found inner ip")
		}
		if len(innerIps) > 1 {
			return innerIps[0], nil
		}
	}

	if strings.HasPrefix(ipVar, NetInterfacePrefix) {
		iface := ipVar[len(NetInterfacePrefix):]
		ip, err := GetIpV4ByIFace(iface)
		if err != nil {
			return "", err
		}
		return ip, nil
	}

	return ipVar, nil
}

func GetIpV4ByIFace(name string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.Name != name {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ip4 := ipNet.IP.To4(); ip4 != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("not found face %s", name))
}

func IsLocalIpV4(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}
	return false
}
