package powerMeter

import (
	"time"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ina219"
)

const (
)

type INA219 struct {
	Device *ina219.Dev
	LastMeasurement time.Time
}

func NewINA219PowerMeter(address int, senseResistorMilliOhm int, maxCurrentMilliAmp int) (*INA219, error) {
	if _, err := driverreg.Init(); err != nil {
		return nil, err
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	var opts = ina219.Opts{Address: address, SenseResistor: senseResistorMilliOhm * physic.MilliOhm, MaxCurrent: maxCurrentMilliAmp * physic.MilliAmpere,}
	var device = ina219.New(b, opts)
	var powerMeter = INA219{Device: device, LastMeasurement: time.Now()}
	return &powerMeter, nil
}

func (receiver *INA219) Reset() {
}

func (receiver *INA219) sense() (ina219.PowerMonitor) {
	var p, err = receiver.Device.Sense()
	if err != nil {
		return nil
	}
	return p
}

func (receiver *INA219) GetEnergy() float64 {
	p := receiver.sense()

	now := time.Now()
	dt := receiver.LastMeasurement.Sub(now)

	dE := p.Power * dt.Seconds()

	receiver.LastMeasurement = now
	return float64(dE)
}
func (receiver *INA219) GetPower() float64 {
	p := receiver.sense()
	return float64(p.Power)
}
func (receiver *INA219) GetCurrent() float64 {
	p := receiver.sense()
	return float64(p.Current)
}
func (receiver *INA219) GetVoltage() float64 {
	p := receiver.sense()
	return float64(p.Voltage)
}
func (receiver *INA219) GetRMSCurrent() float64 {
	return receiver.GetCurrent()
}
func (receiver *INA219) GetRMSVoltage() float64 {
	return receiver.GetVoltage()
}
