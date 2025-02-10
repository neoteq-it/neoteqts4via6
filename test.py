ipv4 = input("IPv4 Adresse")
id = input("ID")

ipv4 = ipv4.split(".")

ipv6 = "fd7a:115c:a1e0:b1a:0:" + f'{int(id):x}' + ":" + f'{int(ipv4[0]):02x}' + f'{int(ipv4[1]):02x}' + ':' + f'{int(ipv4[2]):02x}' + f'{int(ipv4[3]):02x}'


print(ipv6)