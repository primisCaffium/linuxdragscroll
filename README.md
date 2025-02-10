# Linux Drag Scroll

I wrote this utility to execute post-login and enable natural scrolling 
and drag scroll with my trackball, leaving the regular scroll wheel with 
non-natural scrolling.

Tested on Debian 12 using X11.  
Depends on `xinput` and `evtest` 

## Build
Modify the following values accordingly in `src/main.go`:
```
deviceIdXinput := "24"            // xinput device id. Find yours with `xinput list`
buttonIdHoldTriggersScroll := "3" // xinput button id to hold for drag scroll
deviceIdEvtest := "27"            // evtest device id. Find yours with `sudo evtest`
```

Then run
```shell
cd src
go get
go build -o linuxdragscroll
./linuxdragscroll
```

## Install
Open the "install" directory, copy the `.env-example` as `.env` and modify its content to your desired installation path.  
Run `./install.sh`
The "install" script will set it up as autostart on login.  
