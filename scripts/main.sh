#!/bin/sh
# The same extension as a shell script (Linux and macOS; scripts/main.ps1 is the Windows twin).
# Manifest for this variant: runtime = sh, run = scripts/main.sh (see README).
# Protocol: https://rhun.app/extensions/v1/
printf 'hello 1\n'
id=0
request() { id=$((id + 1)); printf '%s %s\n' "$id" "$1"; [ -n "${2-}" ] && printf '%s' "$2"; }
while IFS= read -r line; do
    set -- $line
    case "$1 $2" in
        "0 run")
            case "$3" in
                hello.say) request 'toast 24' 'Hello from a shell script' ;;
                hello.upper)
                    # text comes back as "ok ID VERSION LEN" + LEN bytes; a shell has to read
                    # exactly LEN bytes, which `head -c` does
                    request "text $4"
                    IFS= read -r reply; set -- $reply
                    [ "$1" = ok ] || continue
                    version=$3; len=$4
                    text=$(head -c "$len"; printf x); text=${text%x}
                    upper=$(printf '%s' "$text" | tr '[:lower:]' '[:upper:]')
                    lines=$(printf '%s' "$text" | awk 'END { print NR + (length($0) ? 0 : 0) }')
                    last=$(printf '%s' "$text" | awk 'END { print length($0) + 1 }')
                    request "replace $4 $version 1 1 $lines $last ${#upper}" "$upper"
                    ;;
            esac ;;
        "0 bye") exit 0 ;;
    esac
done
