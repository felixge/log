#!/usr/bin/env bash
set -eu

readonly TEST_BIN='log.test'
readonly CPU_PROFILE='pprof.out'

function print_usage {
    echo "Usage: ${0} [options] <benchmark>"
    echo "  -h|--help: print this help menu and exits"
    echo "  -p|--pdf <file>: create pprof pdf file"
}

if [[ $# -lt 1 ]]; then
  print_usage
  exit 1
fi

pdf=''
while [[ ${#} -gt 0 ]]; do
    case $1 in
        "-h"|"--help") print_usage; exit 0;;
        "-p"|"--pdf") pdf="$2"; shift ;;
        *) break;;
    esac
    shift
done
benchmark="$@"

testargs="-bench ${benchmark}"
if [[ -n "${pdf}" ]]; then
  testargs="${testargs} -x -cpuprofile=${CPU_PROFILE}"
fi

go test ${testargs} .
if [[ -n "${pdf}" ]]; then
  go tool pprof --pdf "${TEST_BIN}" "${CPU_PROFILE}" > "${pdf}"
fi
