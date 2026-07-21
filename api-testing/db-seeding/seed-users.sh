#!/bin/bash
# seed-users.sh — Populate the database with users.
# Usage: ./seed-users.sh
# The server must be running (default: http://localhost:3000).

API_URL="${API_URL:-http://localhost:3000}"

NAMES=(
  "Marta Cabrera"
  "Benoit Blanc"
  "Harlan Thrombey"
  "Linda Drysdale"
  "Richard Drysdale"
  "Ransom Drysdale"
  "Walt Thrombey"
  "Donna Thrombey"
  "Jacob Thrombey"
  "Joni Thrombey"
  "Meg Thrombey"
  "Fran"
  "Lieutenant Elliott"
  "Trooper Wagner"
)

echo "Seeding ${#NAMES[@]} users ..."

for name in "${NAMES[@]}"; do
  curl -s -X POST "$API_URL/session" \
    -H "Content-Type: application/json" \
    -d "{\"userName\":\"$name\"}" > /dev/null
  echo "  ✓ $name"
done

echo "Done."
