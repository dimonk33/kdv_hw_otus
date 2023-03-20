#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | "/c/Work/Utilities/Nmap/ncat.exe" -l localhost 4242 >/d/kuzmindv/temp/nc.out &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >/d/kuzmindv/temp/telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals /d/kuzmindv/temp/nc.out "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals /d/kuzmindv/temp/telnet.out "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"
