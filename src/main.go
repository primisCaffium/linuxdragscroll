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
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	"os/exec"
)

func main() {
	deviceIdXinput := "24"            // xinput device id. Find yours with `xinput list`
	buttonIdHoldTriggersScroll := "3" // xinput button id to hold for drag scroll
	deviceIdEvtest := "27"            // evtest device id. Find yours with `sudo evtest`

	enableDragScrollWhileHoldingButton(deviceIdXinput, buttonIdHoldTriggersScroll)
	handleNaturalScrollingState(deviceIdEvtest, deviceIdXinput)
}

func enableDragScrollWhileHoldingButton(deviceIdXinput, buttonIdHoldTriggersScroll string) {
	runCliOrDie(exec.Command("xinput", "set-prop", deviceIdXinput, "libinput Button Scrolling Button", buttonIdHoldTriggersScroll))
	runCliOrDie(exec.Command("xinput", "set-prop", deviceIdXinput, "libinput Scroll Method Enabled", "0", "0", "1"))
}

func handleNaturalScrollingState(deviceIdEvtest, deviceIdXinput string) {
	device, err := evdev.Open("/dev/input/event" + deviceIdEvtest)
	if err != nil {
		fmt.Println("Error opening evdev device:", err)
		return
	}
	defer func(device *evdev.InputDevice) {
		err := device.Release()
		if err != nil {
			panic(fmt.Sprintf("Error releasing evdev device: %s", err.Error()))
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
							//fmt.Println("Right button pressed - Enabling natural scrolling")
							setNaturalScroll(deviceIdXinput, true)
						} else {
							//fmt.Println("Right button released - Disabling natural scrolling")
							setNaturalScroll(deviceIdXinput, false)
						}
						prevRightButtonPressedState = curRightButtonPressedState
					}
				}
			}
		}
	}
}

func setNaturalScroll(deviceIdXinput string, enable bool) {
	value := "0"
	if enable {
		value = "1"
	}
	runCliOrDie(exec.Command("xinput", "set-prop", deviceIdXinput, "libinput Natural Scrolling Enabled", value))
}

func runCliOrDie(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}
