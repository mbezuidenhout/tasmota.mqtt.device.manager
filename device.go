package tasmota.mqtt.device.manager

type Device struct {
	name, ip string
}

func (d *Device) GetUptime() string {
	uptime := d.name
	return uptime
}