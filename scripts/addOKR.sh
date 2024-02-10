#!/bin/bash

QUARTER=$1
CATEGORY=$2
TYPE=$3
DESCRIPTION=$4
GOAL=$5

OKR_SERVICE="${OKR_SERVICE:-http://localhost:8080}"
set -x
curl -X PUT "${OKR_SERVICE}/api/okr" --data \
    "$(jq -n \
    --arg quarter "${QUARTER}" \
    --arg category "${CATEGORY}" \
    --arg type "${TYPE}" \
    --arg description "${DESCRIPTION}" \
    --arg goal "${GOAL}" \
    '{"quarter": $quarter, "category": $category, "type": $type, "description": $description, "goal": $goal | tonumber}')"
