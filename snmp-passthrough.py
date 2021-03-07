"""
Connects to the local web server and queries the API for metrics.
This is then converted in to a format that the SNMP agent expects
"""
import snmp_passpersist as snmp
import json

PP = snmp.PassPersist(".1.3.6.1.4.1.8072.2.255")


def update():
    with open("/dev/shm/sensors.json", "r") as sensors_file:
        data = json.load(sensors_file)
    index = 0
    for sensor in data:
        index += 1
        if sensor.get("type") in ["status", "reverse_status"]:
            PP.add_int("0.{}".format(index), 1 if sensor.get("value") else 0, sensor.get("name"))
            PP.add_str("1.{}".format(index), sensor.get("name"))
            continue
        if sensor.get("type") in ["flow"]:
            PP.add_int("0.{}".format(index), sensor.get("value", -128) * 1000, sensor.get("name"))
            PP.add_str("1.{}".format(index), sensor.get("name"))
            continue
        if sensor.get("type") in ["rpm", "temperature", "pa_power"]:
            PP.add_int("0.{}".format(index), sensor.get("value", -128), sensor.get("name"))
            PP.add_str("1.{}".format(index), sensor.get("name"))
            continue
        if sensor.get("type") == "temp_fan":
            PP.add_str("1.{}".format(index), sensor.get("name"))
            PP.add_int("0.{}.0".format(index), sensor.get("value", {}).get("temp", -128) * 1000, "{}:{}".format(sensor.get("name"), "cpu-temp"))
            PP.add_str("1.{}.0".format(index), "cpu-temp")
            PP.add_int("0.{}.1".format(index), sensor.get("value", {}).get("fan", -1), "{}:{}".format(sensor.get("name"), "fan-status"))
            PP.add_str("1.{}.1".format(index), "fan-status")
            continue
        if sensor.get("type") == "dht22":
            PP.add_str("1.{}".format(index), sensor.get("name"))
            PP.add_int("0.{}.0".format(index), sensor.get("value", {}).get("temperature", -128) * 1000, "{}:{}".format(sensor.get("name"), "temperature"))
            PP.add_str("1.{}.0".format(index), "temperature")
            PP.add_int("0.{}.1".format(index), sensor.get("value", {}).get("humidity", -128) * 1000, "{}:{}".format(sensor.get("name"), "humidity"))
            PP.add_str("1.{}.1".format(index), "humidity")
            continue
        if sensor.get("type") == "power":
            PP.add_str("1.{}".format(index), sensor.get("name"))
            PP.add_int("0.{}.0".format(index), sensor.get("value", {}).get("power", -1) * 1000, "{}:{}".format(sensor.get("name"), "power"))
            PP.add_str("1.{}.0".format(index), "power")
            PP.add_int("0.{}.1".format(index), sensor.get("value", {}).get("current", -1) * 1000, "{}:{}".format(sensor.get("name"), "current"))
            PP.add_str("1.{}.1".format(index), "current")
            PP.add_int("0.{}.2".format(index), sensor.get("value", {}).get("voltage", -1) * 1000, "{}:{}".format(sensor.get("name"), "voltage"))
            PP.add_str("1.{}.2".format(index), "voltage")
            continue


def main():
    PP.start(update, 30)
    PP.debug = True


if __name__ == "__main__":
    main()
