/*
This file is part of primisCaffium/linuxdragscroll.

primisCaffium/linuxdragscroll is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software Foundation, either
version 3 of the License, or (at your option) any later version.

primisCaffium/linuxdragscroll is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Foobar. If not,
see <https://www.gnu.org/licenses/>.

Author: Primis Caffium
Description:

	Allows you to "drag scroll" in "natural scrolling" mode while having
	your global "natural scrolling" off, so your scroll wheel is unaffected.
	Evtest is used to detect button hold and then change the natural scrolling behavior.
	This is currently untested on wayland.

	To run this as your regular user, you need to add your user to the "input" group,
	otherwise you'd have to run this as root, and it's not recommended.
	```
		Sudo usermod -aG input <your_username>
	```

	Modify the constants to match your device and desired button.
*/
package main

import (
	"bytes"
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Tool struct {
	deviceIdXinput             string
	buttonIdHoldTriggersScroll string
	evtestDevicePath           string
}

func main() {
	deviceIdXinput := "24"            // xinput device id. Find yours with `xinput list`
	buttonIdHoldTriggersScroll := "3" // xinput button id to hold for drag scroll
	deviceIdEvtest := "27"            // evtest device id. Find yours with `sudo evtest`

	t := Tool{
		deviceIdXinput:             deviceIdXinput,
		buttonIdHoldTriggersScroll: buttonIdHoldTriggersScroll,
		evtestDevicePath:           "/dev/input/event" + deviceIdEvtest,
	}
	t.Start()
}

func (o *Tool) Start() {
	for {
		o.waitForDevice()
		o.enableDragScrollWhileHoldingButton()
		o.handleNaturalScrollingState()
		time.Sleep(3 * time.Second)
	}
}

func (o *Tool) enableDragScrollWhileHoldingButton() {
	o.runCliOrDie(exec.Command("xinput", "set-prop", o.deviceIdXinput, "libinput Button Scrolling Button", o.buttonIdHoldTriggersScroll))
	o.runCliOrDie(exec.Command("xinput", "set-prop", o.deviceIdXinput, "libinput Scroll Method Enabled", "0", "0", "1"))
}

func (o *Tool) waitForDevice() {
	i := 0
	for {
		_, err := os.Stat(o.evtestDevicePath)
		if err == nil {
			fmt.Printf("evtest device %s is available.\n", o.evtestDevicePath)

			for {
				if math.Mod(float64(i), 10) == 0 {
					fmt.Printf("Waiting for xinput device.\n")
				}

				xinputReady := o.isXinputDeviceAvailable()
				if xinputReady {
					fmt.Printf("xinput device %s is available.\n", o.deviceIdXinput)
					break
				}
				time.Sleep(1 * time.Second)
				i++
			}

			break
		}

		if math.Mod(float64(i), 10) == 0 {
			fmt.Printf("Device %s is not available yet. Waiting...\n", o.evtestDevicePath)
		}
		time.Sleep(1 * time.Second)
		i++
	}
}

func (o *Tool) isXinputDeviceAvailable() bool {
	cmd := exec.Command("xinput", "list")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running xinput command: %s\n", err)
		return false
	}

	output := out.String()
	return strings.Contains(output, "id="+o.deviceIdXinput)
}

func (o *Tool) openEvtestDevice() *evdev.InputDevice {
	device, err := evdev.Open(o.evtestDevicePath)
	if err != nil {
		fmt.Println("Error opening evdev device:", err)
		return nil
	}
	return device
}

func (o *Tool) handleNaturalScrollingState() {
	device := o.openEvtestDevice()
	if device == nil {
		fmt.Println(fmt.Sprintf("Error opening evdev device."))
		return
	}
	defer func(device *evdev.InputDevice) {
		if device == nil {
			return
		}
		err := device.Release()
		if err != nil {
			fmt.Println(fmt.Sprintf("Error releasing evdev device: %s", err.Error()))
		}
	}(device)

	var prevRightButtonPressedState byte
	fmt.Println("Listening for mouse events...")

	for {
		events, err := device.Read() // blocking
		if err != nil {
			fmt.Println("Error reading input event:", err)
			return
		}

		for _, ev := range events {

			if ev.Type == evdev.EV_KEY {
				if ev.Code == evdev.BTN_RIGHT {

					curRightButtonPressedState := byte(ev.Value)

					if curRightButtonPressedState != prevRightButtonPressedState {
						if curRightButtonPressedState > 0 {
							o.setNaturalScroll(true)
						} else {
							o.setNaturalScroll(false)
						}
						prevRightButtonPressedState = curRightButtonPressedState
					}
				}
			}
		}
	}
}

func (o *Tool) setNaturalScroll(enable bool) {
	value := "0"
	if enable {
		value = "1"
	}
	o.runCliOrDie(exec.Command("xinput", "set-prop", o.deviceIdXinput, "libinput Natural Scrolling Enabled", value))
}

func (o *Tool) runCliOrDie(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}
