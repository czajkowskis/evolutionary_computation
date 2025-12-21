#!/bin/bash
# Helper script to update lab imports (for reference - manual updates are being done)

LAB_DIR=$1
LAB_NUM=$2

if [ -z "$LAB_DIR" ] || [ -z "$LAB_NUM" ]; then
    echo "Usage: $0 <lab_dir> <lab_num>"
    exit 1
fi

echo "Updating $LAB_DIR (Lab $LAB_NUM)"

