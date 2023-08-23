#!/bin/sh
# .bashrc

export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm

# StartX and start chromium in full screen kiosk mode
# if [ -z "$DISPLAY" ] && [ $(tty) = /dev/tty1 ]; then
#     startx
# fi

