#!/bin/bash
#
# This deploy script assumes the main deploy script has been run to rsync
# necessary files to a non-production location on the production server.  From
# there we stop services, move all files into the production location, migrate
# the database, and restart services.
set -eu

ESC=$(echo -e "\x1b")
CSI="${ESC}["
EXAMPLE="    ${CSI}1m"
WARN="${CSI}31;1m"
INFO="${CSI}34;1m"
RESET="${CSI}0m"

normal() {
  echo $@
}
info() {
  echo "${INFO}$@${RESET}"
}
process_start() {
  echo
  echo "--- ${INFO}$@${RESET} ---"
}
process_done() {
  echo "--- ${INFO}Done${RESET} ---"
}
warn() {
  echo "${WARN}$@${RESET}"
}
example() {
  echo "${EXAMPLE}$@${RESET}"
}

src=$(pwd)
scidir=/usr/local/sci

process_start "Stopping services"
systemctl stop httpd || true
systemctl stop sci-httpd || true
process_done

process_start "Copying production files to $scidir"
mkdir -p $scidir
./scripts/rsync_production_files.sh $src $scidir
chown -R root:root $scidir
chmod -R o-rwx $scidir
process_done

cd $scidir

process_start "Checking required config/service files"

conffile=/etc/sci.conf
if [[ ! -f $conffile ]]; then
  warn "No configuration detected"
  normal "You should copy (and edit), and secure the example configuration:"
  example "cp $src/example.conf $conffile"
  example "chmod 640 $conffile"
  exit 1
fi

syslogfile=/etc/rsyslog.d/sci.conf
if [[ ! -f $syslogfile ]]; then
  warn "No rsyslog configuration detected"
  normal "You should copy the example syslog file:"
  example "cp $src/rhel7/example.rsyslog.conf $syslogfile"
  exit 1
fi

unitfile=$scidir/sci-httpd.service
if [[ ! -f $unitfile ]]; then
  warn No systemd unit file detected
  normal "You should copy the example service unit file:"
  example "cp $src/rhel7/example.service $unitfile"
  exit 1
fi

process_done

process_start "Migrating the database"
./scripts/dbmigrate.sh
process_done

process_start "Doing a daemon reload and starting the service"
systemctl enable $unitfile
systemctl daemon-reload
# We ignore failures with httpd because in staging environments we don't actually need apache
systemctl start httpd || true
systemctl start sci-httpd
systemctl restart rsyslog
process_done
