#!/bin/bash
SCRIPT_DIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

for badf in test_data/bad*.eml; do
    echo -n "$badf: "
    if ./sieve_date_offset_filter < "$badf"; then
        echo "Filter should have failed for $badf; but it reported OK."
        exit 1
    fi
    echo "OK"
done

for goodf in test_data/good*.eml; do
    echo -n "$goodf: "
    ./sieve_date_offset_filter < "$goodf"
    if [ $? != 0 ]; then
        echo "Filter should have succeeded on $goodf but it failed."
        exit 1
    fi
    echo "OK"
done

echo "All OK :)"