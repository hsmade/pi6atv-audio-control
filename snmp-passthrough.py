"""
Connects to the local web server and queries the API for metrics.
This is then converted in to a format that the SNMP agent expects
"""
import snmp_passpersist as snmp
import requests

PP = snmp.PassPersist(".1.3.6.1.4.1.8072.2.255")
ports = {
    '0': "DSP 1",
    '1': "DSP 2",
    '2': "DSP 3",
    '3': "DSP 4",
    '4': "DSP 5",
    '5': "program",
    '6': "reset",
    '7': "carrier 7.02",
    '10': "carrier 7.20",
    '11': "carrier 7.38",
    '12': "carrier 7.56",
    '13': "carrier 7.74",
    '14': "carrier 7.92",
    '15': "carrier nicam 5.85",
    '16': "carrier nicam 6.552",
}


def update():
    PP.add_str("1.255", "error")
    try:
        data = requests.get("http://localhost/control/").json()
    except:
        PP.add_int("0.255", 1)
    else:
        PP.add_int("0.255", 0)
        for port, name in ports.items():
            PP.add_int("0.{}".format(port), (0, 1)[data[port]])
            PP.add_str("1.{}".format(port), name)


def main():
    PP.start(update, 1)
    PP.debug = True


if __name__ == "__main__":
    main()
