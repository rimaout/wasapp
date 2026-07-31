#!/bin/bash
# test-chat-endpoints.sh — Validate all chat API endpoints
# Usage: ./test-chat-endpoints.sh
# Run from project root:  ./api-testing/db-seeding/test-chat-endpoints.sh

set -euo pipefail
cd "$(dirname "$0")/../.."
API_URL="${API_URL:-http://localhost:3000}"
PASS=0
FAIL=0

UUID_RE='^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'

# --- Helper ---
check() {
	local label="$1" method="$2" url="$3" auth="$4" body="$5" expect="$6" extra="${7:-}"
	local code=0

	case "$method" in
		DELETE)
			code=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" -X DELETE "$url" -H "Authorization: $auth" --max-time 5)
			;;
		GET)
			code=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" "$url" -H "Authorization: $auth" --max-time 5)
			;;
		*)
			code=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" -X "$method" "$url" \
				-H "Authorization: $auth" \
				-H "Content-Type: application/json" \
				-d "$body" --max-time 5)
			;;
	esac

	local ok=1
	[[ "$code" == "$expect" ]] || ok=0
	if [[ -n "$extra" ]] && [[ $ok -eq 1 ]]; then
		eval "$extra" 2>/dev/null || ok=0
	fi

	if [[ $ok -eq 1 ]]; then
		echo "  ✓ $label ($code)"
		PASS=$((PASS + 1))
	else
		echo "  ✗ $label (expected $expect, got $code)"
		FAIL=$((FAIL + 1))
	fi
}

# --- Start server ---
echo "Starting server..."
lsof -ti:3000 2>/dev/null | xargs kill 2>/dev/null || true
rm -f /tmp/decaf.db
rm -rf uploads/
go run ./cmd/webapi/ > /tmp/wasapp-server.log 2>&1 &
SERVER_PID=$!
sleep 3

# --- Seed users ---
echo "Seeding users..."
./api-testing/db-seeding/seed-users.sh > /dev/null 2>&1

# --- Login all 4 users ---
echo "Logging in..."
RESA=$(curl -s -X POST "$API_URL/session" -H "Content-Type: application/json" -d '{"userName":"Marta Cabrera"}' --max-time 5)
RESB=$(curl -s -X POST "$API_URL/session" -H "Content-Type: application/json" -d '{"userName":"Harlan Thrombey"}' --max-time 5)
RESC=$(curl -s -X POST "$API_URL/session" -H "Content-Type: application/json" -d '{"userName":"Ransom Drysdale"}' --max-time 5)
RESD=$(curl -s -X POST "$API_URL/session" -H "Content-Type: application/json" -d '{"userName":"Walt Thrombey"}' --max-time 5)
TOKA=$(echo "$RESA" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])" 2>/dev/null)
UIDA=$(echo "$RESA" | python3 -c "import sys,json; print(json.load(sys.stdin)['userId'])" 2>/dev/null)
TOKB=$(echo "$RESB" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])" 2>/dev/null)
UIDB=$(echo "$RESB" | python3 -c "import sys,json; print(json.load(sys.stdin)['userId'])" 2>/dev/null)
TOKC=$(echo "$RESC" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])" 2>/dev/null)
UIDC=$(echo "$RESC" | python3 -c "import sys,json; print(json.load(sys.stdin)['userId'])" 2>/dev/null)
TOKD=$(echo "$RESD" | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])" 2>/dev/null)
UIDD=$(echo "$RESD" | python3 -c "import sys,json; print(json.load(sys.stdin)['userId'])" 2>/dev/null)

# --- Direct Chat ---
echo ""
echo "=== Direct Chat ==="
check "Create A→B" POST "$API_URL/user/$UIDB/chats" "$TOKA" "" "201" \
	'python3 -c "import sys,json; d=json.load(open(\"/tmp/test-resp.json\")); assert d[\"displayName\"]==\"Harlan Thrombey\", f\"wrong displayName: {d}\"; assert d[\"isGroupChat\"]==False; assert \"$UUID_RE\" and True"'
DCID=$(grep -o '"id":"[^"]*"' /tmp/test-resp.json | head -1 | cut -d'"' -f4)
DCID="${DCID//[$'\t\r\n ']/}"

check "Create duplicate" POST "$API_URL/user/$UIDB/chats" "$TOKA" "" "409"
check "Create self" POST "$API_URL/user/$UIDA/chats" "$TOKA" "" "400"
check "Create unknown user" POST "$API_URL/user/00000000-0000-0000-0000-000000000000/chats" "$TOKA" "" "404"

check "Get members (as A)" GET "$API_URL/chats/$DCID/members" "$TOKA" "" "200" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); assert len(d)==2, f\"expected 2 members, got {len(d)}\""'
check "Get members (as C, not member)" GET "$API_URL/chats/$DCID/members" "$TOKC" "" "403"
check "Get members (not found)" GET "$API_URL/chats/00000000-0000-0000-0000-000000000000/members" "$TOKA" "" "404"

# --- Group Chat ---
echo ""
echo "=== Group Chat ==="
check "Create Group A+B+C" POST "$API_URL/chats" "$TOKA" \
	"{\"groupName\":\"Thrombey Legacy\",\"membersList\":[\"$UIDB\",\"$UIDC\"]}" "201" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); assert d[\"displayName\"]==\"Thrombey Legacy\"; assert d[\"isGroupChat\"]==True"'
GCID=$(grep -o '"id":"[^"]*"' /tmp/test-resp.json | head -1 | cut -d'"' -f4)
GCID="${GCID//[$'\t\r\n ']/}"
[[ -z "$GCID" ]] && { echo "ERROR: Failed to extract group chat ID"; FAIL=$((FAIL + 1)); }

check "Create bad name" POST "$API_URL/chats" "$TOKA" \
	"{\"groupName\":\"!!\",\"membersList\":[\"$UIDB\"]}" "400"
check "Create empty list" POST "$API_URL/chats" "$TOKA" \
	"{\"groupName\":\"Empty\",\"membersList\":[]}" "400"
check "Create duplicate member" POST "$API_URL/chats" "$TOKA" \
	"{\"groupName\":\"Dup\",\"membersList\":[\"$UIDB\",\"$UIDB\"]}" "409"
check "Create unknown member" POST "$API_URL/chats" "$TOKA" \
	"{\"groupName\":\"Ghost\",\"membersList\":[\"00000000-0000-0000-0000-000000000000\"]}" "404"

check "Get group members (as A)" GET "$API_URL/chats/$GCID/members" "$TOKA" "" "200" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); assert len(d)==3, f\"expected 3, got {len(d)}\""'
check "Get group members (as D, not member)" GET "$API_URL/chats/$GCID/members" "$TOKD" "" "403"

# --- Group Operations ---
echo ""
echo "=== Group Operations ==="
check "Add D to group" POST "$API_URL/chats/$GCID/members" "$TOKA" \
	"{\"userId\":\"$UIDD\"}" "201"
check "Add D again (already in)" POST "$API_URL/chats/$GCID/members" "$TOKA" \
	"{\"userId\":\"$UIDD\"}" "409"
check "Add unknown user" POST "$API_URL/chats/$GCID/members" "$TOKA" \
	"{\"userId\":\"00000000-0000-0000-0000-000000000000\"}" "404"
check "Add to direct chat (not group)" POST "$API_URL/chats/$DCID/members" "$TOKA" \
	"{\"userId\":\"$UIDD\"}" "409"

check "Leave group (as D)" DELETE "$API_URL/chats/$GCID/members/me" "$TOKD" "" "204"
check "Leave group (D already left)" DELETE "$API_URL/chats/$GCID/members/me" "$TOKD" "" "403"
check "Leave direct chat (not group)" DELETE "$API_URL/chats/$DCID/members/me" "$TOKA" "" "409"
check "Leave group (not found)" DELETE "$API_URL/chats/00000000-0000-0000-0000-000000000000/members/me" "$TOKA" "" "404"

# --- Group Name ---
echo ""
echo "=== Group Name ==="
check "Set group name" PATCH "$API_URL/chats/$GCID/name" "$TOKA" \
	"{\"groupName\":\"Thrombey Dynasty\"}" "200" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); assert d[\"groupName\"]==\"Thrombey Dynasty\""'
check "Set name (bad name)" PATCH "$API_URL/chats/$GCID/name" "$TOKA" \
	"{\"groupName\":\"!!\"}" "400"
check "Set name (non-member)" PATCH "$API_URL/chats/$GCID/name" "$TOKD" \
	"{\"groupName\":\"Hijacked\"}" "403"
check "Set name (not a group)" PATCH "$API_URL/chats/$DCID/name" "$TOKA" \
	"{\"groupName\":\"FakeGroup\"}" "409"

# --- Avatar ---
echo ""
echo "=== Avatar ==="
cp service/api/assets/default-group-avatar.png /tmp/test-avatar.png

# Avatar upload uses direct curl (not the check helper) because check() only handles JSON bodies
AVATAR_CODE=$(curl -s -o /tmp/test-avatar-resp.png -w "%{http_code}" -X PUT "$API_URL/chats/$GCID/avatar" \
	-H "Authorization: $TOKA" \
	-F "binaryImage=@/tmp/test-avatar.png;type=image/png" --max-time 5)
if [[ "$AVATAR_CODE" == "204" ]]; then
	echo "  ✓ Set group avatar (204)"; PASS=$((PASS + 1))
else
	echo "  ✗ Set group avatar (expected 204, got $AVATAR_CODE)"; FAIL=$((FAIL + 1))
fi

check "Get group avatar" GET "$API_URL/chats/$GCID/avatar" "$TOKA" "" "200" \
	'file /tmp/test-resp.json | grep -qi "PNG image"'
check "Get direct chat avatar" GET "$API_URL/chats/$DCID/avatar" "$TOKA" "" "200" \
	'file /tmp/test-resp.json | grep -qi "PNG image"'
check "Get avatar (non-member)" GET "$API_URL/chats/$GCID/avatar" "$TOKD" "" "403"
check "Get avatar (not found)" GET "$API_URL/chats/00000000-0000-0000-0000-000000000000/avatar" "$TOKA" "" "404"

# --- Summary ---
echo ""
echo "============================================"
echo " RESULTS: $PASS passed, $FAIL failed"
echo "============================================"

# --- Cleanup ---
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
rm -f /tmp/decaf.db /tmp/test-resp.json /tmp/test-avatar.png /tmp/test-avatar-resp.png /tmp/wasapp-server.log
rm -rf uploads/

if [[ $FAIL -gt 0 ]]; then
	exit 1
fi
exit 0
