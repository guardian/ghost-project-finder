#!/bin/bash -e

echo Installing ghost-project-finder to /usr/local/bin and the associated plist to /Library/LaunchDaemons

declare -x PLIST=com.gu.ghost-project-finder.plist

install -g 0 -o 0 -m 755 ghost-project-finder /usr/local/bin
cp ghost-project-finder.plist /Library/LaunchDaemons/$PLIST

cd /Library/LaunchDaemons
launchctl unload $PLIST || /bin/true 2>/dev/null #needed if it's been started before, but fails if it hasn't so silence the error
launchctl load $PLIST
launchctl start $PLIST

echo Installation completed successfully