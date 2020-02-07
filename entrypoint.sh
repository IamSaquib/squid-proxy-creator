#!/bin/bash
set -e
SQUID_USER=proxy
SQUID_LOG_DIR=/var/log/squid
SQUID_CACHE_DIR=/var/spool/squid
SQUID_VERSION=3.5.27   
create_log_dir() {
  mkdir -p ${SQUID_LOG_DIR}
  chmod -R 755 ${SQUID_LOG_DIR}
  chown -R ${SQUID_USER}:${SQUID_USER} ${SQUID_LOG_DIR}
}

create_cache_dir() {
  mkdir -p ${SQUID_CACHE_DIR}
  chown -R ${SQUID_USER}:${SQUID_USER} ${SQUID_CACHE_DIR}
}

create_log_dir
create_cache_dir

# default behaviour is to launch squid
echo "Initializing cache..."
$(which squid) -N -f /etc/squid/squid.conf -z
sleep 5

echo "Starting squid..."
exec $(which squid) -f /etc/squid/squid.conf -NYCd 1
