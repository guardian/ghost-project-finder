#!/bin/bash -e

echo Installing ghost-project-finder to /usr/local/bin and the associated plist to /Library/LaunchDaemons

declare -x DLPLIST=com.gu.ghost-project-finder-downloads.plist
declare -x RTPLIST=com.gu.ghost-project-finder-root.plist
install -g 0 -o 0 -m 755 ghost-project-finder /usr/local/bin
cp ghost-project-finder-downloads.plist /Library/LaunchDaemons/$DLPLIST
cp ghost-project-finder-root.plist /Library/LaunchDaemons/$RTPLIST
cd /Library/LaunchDaemons
launchctl unload $DLPLIST || /bin/true 2>/dev/null #needed if it's been started before, but fails if it hasn't so silence the error
launchctl load $DLPLIST
launchctl start $DLPLIST

launchctl unload $RTPLIST || /bin/true 2>/dev/null #needed if it's been started before, but fails if it hasn't so silence the error
launchctl load $RTPLIST
launchctl start $RTPLIST

echo Installation completed successfully