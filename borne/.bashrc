#!/bin/sh
# .bashrc

# StartX and start chromium in full screen kiosk mode
if [ -z "$DISPLAY" ] && [ $(tty) = /dev/tty1 ]; then
    startx
fi

