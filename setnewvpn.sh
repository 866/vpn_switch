nmcli connection delete ua-vpn
nmcli connection import type wireguard file "/etc/wireguard/wg0.conf"
nmcli connection modify wg0 connection.id "ua-vpn"