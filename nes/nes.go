package nes

import (
	"errors"
	"fmt"
	"github.com/nwidger/m65go2"
	"github.com/nwidger/rp2ago3"
	"github.com/nwidger/rp2cgo2"
	"time"
)

type NES struct {
	cpu         *rp2ago3.RP2A03
	ppu         *rp2cgo2.RP2C02
	controllers *Controllers
	clock       m65go2.Clocker
	rom         ROM
}

func NewNES(filename string) (nes *NES, err error) {
	var rate time.Duration
	var cpuDivisor uint64

	rom, err := NewROM(filename)

	if err != nil {
		err = errors.New(fmt.Sprintf("Error loading ROM: %v", err))
		return
	}

	switch rom.Region() {
	case NTSC:
		rate = rp2ago3.NTSC_CLOCK_RATE
		cpuDivisor = rp2ago3.NTSC_CPU_CLOCK_DIVISOR
	case PAL:
		rate = rp2ago3.PAL_CLOCK_RATE
		cpuDivisor = rp2ago3.PAL_CPU_CLOCK_DIVISOR
	}

	clock := m65go2.NewClock(rate)
	cpu := rp2ago3.NewRP2A03(clock, cpuDivisor)
	ppu := rp2cgo2.NewRP2C02(clock, cpu.InterruptLine(m65go2.Nmi), rom.Mirroring())
	ctrls := NewControllers()

	cpu.Memory.AddMappings(ppu, rp2ago3.CPU)
	cpu.Memory.AddMappings(rom, rp2ago3.CPU)
	cpu.Memory.AddMappings(ctrls, rp2ago3.CPU)

	ppu.Memory.AddMappings(rom, rp2ago3.PPU)

	nes = &NES{cpu: cpu, ppu: ppu, clock: clock, rom: rom}
	return
}

func (nes *NES) Reset() {
	nes.cpu.Reset()
	nes.ppu.Reset()
}

func (nes *NES) Run() (err error) {
	nes.Reset()
	nes.clock.Start()

	go nes.cpu.Run()
	go nes.ppu.Run()

	for {
		time.Sleep(9999 * time.Second)
	}

	return
}
