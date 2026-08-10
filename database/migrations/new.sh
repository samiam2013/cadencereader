#!/bin/bash

set -euo pipefail

read -rp "migration name: " title

migrate create -ext sql -seq $title;

