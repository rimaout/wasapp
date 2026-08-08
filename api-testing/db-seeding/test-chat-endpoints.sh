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

# --- Messages ---
echo ""
echo "=== Messages ==="

# Send a few messages to populate the chat for subsequent tests
MSG_A1_CODE=$(curl -s -o /tmp/msg-a1.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages" \
	-H "Authorization: $TOKA" -F "text=First message from A" --max-time 5)
if [[ "$MSG_A1_CODE" == "201" ]]; then
	echo "  ✓ Seed MSG_A1 (201)"; PASS=$((PASS + 1))
	MSG_A1=$(python3 -c "import json; print(json.load(open('/tmp/msg-a1.json'))['id'])" 2>/dev/null)
else
	echo "  ✗ Seed MSG_A1 (expected 201, got $MSG_A1_CODE)"; FAIL=$((FAIL + 1)); MSG_A1=""
fi

MSG_A2_CODE=$(curl -s -o /tmp/msg-a2.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages" \
	-H "Authorization: $TOKA" -F "text=Second message from A" --max-time 5)
if [[ "$MSG_A2_CODE" == "201" ]]; then
	echo "  ✓ Seed MSG_A2 (201)"; PASS=$((PASS + 1))
	MSG_A2=$(python3 -c "import json; print(json.load(open('/tmp/msg-a2.json'))['id'])" 2>/dev/null)
else
	echo "  ✗ Seed MSG_A2 (expected 201, got $MSG_A2_CODE)"; FAIL=$((FAIL + 1)); MSG_A2=""
fi

cp service/api/assets/default-group-avatar.png /tmp/test-msg-image.png
MSG_IMG_CODE=$(curl -s -o /tmp/msg-img.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages" \
	-H "Authorization: $TOKA" -F "text=Image message" -F "imageFile=@/tmp/test-msg-image.png;type=image/png" --max-time 5)
if [[ "$MSG_IMG_CODE" == "201" ]]; then
	echo "  ✓ Seed MSG_IMG (201)"; PASS=$((PASS + 1))
	MSG_IMG=$(python3 -c "import json; print(json.load(open('/tmp/msg-img.json'))['id'])" 2>/dev/null)
	IMG_ID=$(python3 -c "import json; print(json.load(open('/tmp/msg-img.json'))['content']['msgImageId'])" 2>/dev/null)
else
	echo "  ✗ Seed MSG_IMG (expected 201, got $MSG_IMG_CODE)"; FAIL=$((FAIL + 1)); MSG_IMG=""; IMG_ID=""
fi

MSG_B_REPLY_CODE=$(curl -s -o /tmp/msg-b-reply.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages" \
	-H "Authorization: $TOKB" -F "text=Message from B" --max-time 5)
if [[ "$MSG_B_REPLY_CODE" == "201" ]]; then
	echo "  ✓ Seed MSG_B (201)"; PASS=$((PASS + 1))
	MSG_B=$(python3 -c "import json; print(json.load(open('/tmp/msg-b-reply.json'))['id'])" 2>/dev/null)
else
	echo "  ✗ Seed MSG_B (expected 201, got $MSG_B_REPLY_CODE)"; FAIL=$((FAIL + 1)); MSG_B=""
fi

# Get the init message ID
INIT_MSG=$(curl -s "$API_URL/chats/$DCID/messages" -H "Authorization: $TOKA" --max-time 5 | \
	python3 -c "import sys,json; [print(m['id']) for m in json.load(sys.stdin)['messages'] if m['isInitMessage']]" 2>/dev/null)
[[ -z "$INIT_MSG" ]] && { echo "  ✗ Could not find init message ID"; FAIL=$((FAIL + 1)); }

# Send message (non-member)
MSG_NONMEMBER_CODE=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages" \
	-H "Authorization: $TOKC" -F "text=Outsider" --max-time 5)
if [[ "$MSG_NONMEMBER_CODE" == "403" ]]; then
	echo "  ✓ Send non-member (403)"; PASS=$((PASS + 1))
else
	echo "  ✗ Send non-member (expected 403, got $MSG_NONMEMBER_CODE)"; FAIL=$((FAIL + 1))
fi

# Send empty content (400)
MSG_EMPTY_CODE=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages" \
	-H "Authorization: $TOKA" -F "text= " --max-time 5)
if [[ "$MSG_EMPTY_CODE" == "400" ]]; then
	echo "  ✓ Send empty content (400)"; PASS=$((PASS + 1))
else
	echo "  ✗ Send empty content (expected 400, got $MSG_EMPTY_CODE)"; FAIL=$((FAIL + 1))
fi

# Get conversation (ordered, includes init + ordinary)
check "Get messages (ordered)" GET "$API_URL/chats/$DCID/messages" "$TOKA" "" "200" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); msgs=d[\"messages\"]; assert len(msgs)>=5, f\"expected >=5, got {len(msgs)}\"; times=[m[\"sendTime\"] for m in msgs]; assert times==sorted(times,reverse=True), f\"not sorted: {times}\""'
check "Get messages non-member" GET "$API_URL/chats/$DCID/messages" "$TOKC" "" "403"
check "Get messages chat not found" GET "$API_URL/chats/00000000-0000-0000-0000-000000000000/messages" "$TOKA" "" "404"

# Delete own message
MSG_DEL_CODE=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" -X DELETE "$API_URL/chats/$DCID/messages/$MSG_A1" \
	-H "Authorization: $TOKA" --max-time 5)
if [[ "$MSG_DEL_CODE" == "200" ]]; then
	is_del=$(python3 -c "import json; print(json.load(open('/tmp/test-resp.json'))['isDeleted'])" 2>/dev/null || true)
	content=$(python3 -c "import json; print(json.load(open('/tmp/test-resp.json')).get('content'))" 2>/dev/null || true)
	if [[ "$is_del" == "True" && "$content" == "None" ]]; then
		echo "  ✓ Delete own message (200)"; PASS=$((PASS + 1))
	else
		echo "  ✗ Delete own message (isDeleted=$is_del, content=$content)"; FAIL=$((FAIL + 1))
	fi
else
	echo "  ✗ Delete own message (expected 200, got $MSG_DEL_CODE)"; FAIL=$((FAIL + 1))
fi

# Delete non-owner (B tries to delete A's message)
check "Delete non-owner" DELETE "$API_URL/chats/$DCID/messages/$MSG_A2" "$TOKB" "" "403"

# Delete already deleted
check "Delete already deleted" DELETE "$API_URL/chats/$DCID/messages/$MSG_A1" "$TOKA" "" "400"

# Delete init message
check "Delete init message" DELETE "$API_URL/chats/$DCID/messages/$INIT_MSG" "$TOKA" "" "400"

# Delete non-existent
check "Delete non-existent" DELETE "$API_URL/chats/$DCID/messages/00000000-0000-0000-0000-000000000000" "$TOKA" "" "404"

# Reply to message (text)
MSG_REPLY_CODE=$(curl -s -o /tmp/msg-reply.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages/$MSG_A2/reply" \
	-H "Authorization: $TOKA" -F "text=This is a reply" --max-time 5)
if [[ "$MSG_REPLY_CODE" == "201" ]]; then
	rr=$(python3 -c "import json; m=json.load(open('/tmp/msg-reply.json')); print(m['replyTo'])" 2>/dev/null)
	if [[ "$rr" == "$MSG_A2" ]]; then
		echo "  ✓ Reply text (201, replyTo=$rr)"; PASS=$((PASS + 1))
	else
		echo "  ✗ Reply text (wrong replyTo: $rr, expected $MSG_A2)"; FAIL=$((FAIL + 1))
	fi
else
	echo "  ✗ Reply text (expected 201, got $MSG_REPLY_CODE)"; FAIL=$((FAIL + 1))
fi

# Reply with text+image
MSG_REPLY2_CODE=$(curl -s -o /tmp/msg-reply2.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages/$MSG_A2/reply" \
	-H "Authorization: $TOKA" -F "text=Reply with image" -F "imageFile=@/tmp/test-msg-image.png;type=image/png" --max-time 5)
if [[ "$MSG_REPLY2_CODE" == "201" ]]; then
	rr2=$(python3 -c "import json; m=json.load(open('/tmp/msg-reply2.json')); print(m['replyTo'])" 2>/dev/null)
	c2=$(python3 -c "import json; m=json.load(open('/tmp/msg-reply2.json')); c=m['content']; print(c['text'], c.get('msgImageId','')[:8])" 2>/dev/null)
	if [[ "$rr2" == "$MSG_A2" ]]; then
		echo "  ✓ Reply text+image (201, replyTo=$rr2)"; PASS=$((PASS + 1))
	else
		echo "  ✗ Reply text+image (wrong replyTo: $rr2)"; FAIL=$((FAIL + 1))
	fi
else
	echo "  ✗ Reply text+image (expected 201, got $MSG_REPLY2_CODE)"; FAIL=$((FAIL + 1))
fi

# Reply to non-existent message
MSG_REPLY_404_CODE=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages/00000000-0000-0000-0000-000000000000/reply" \
	-H "Authorization: $TOKA" -F "text=Reply to nothing" --max-time 5)
if [[ "$MSG_REPLY_404_CODE" == "404" ]]; then
	echo "  ✓ Reply to non-existent (404)"; PASS=$((PASS + 1))
else
	echo "  ✗ Reply to non-existent (expected 404, got $MSG_REPLY_404_CODE)"; FAIL=$((FAIL + 1))
fi

# Reply non-member
MSG_REPLY_403_CODE=$(curl -s -o /tmp/test-resp.json -w "%{http_code}" -X POST "$API_URL/chats/$DCID/messages/$MSG_A2/reply" \
	-H "Authorization: $TOKC" -F "text=Outsider reply" --max-time 5)
if [[ "$MSG_REPLY_403_CODE" == "403" ]]; then
	echo "  ✓ Reply non-member (403)"; PASS=$((PASS + 1))
else
	echo "  ✗ Reply non-member (expected 403, got $MSG_REPLY_403_CODE)"; FAIL=$((FAIL + 1))
fi

# Serve message image
if [[ -n "$IMG_ID" ]]; then
	IMG_SERVE_CODE=$(curl -s -o /tmp/served-img.png -w "%{http_code}" "$API_URL/chats/$DCID/images/$IMG_ID" \
		-H "Authorization: $TOKA" --max-time 5)
	if [[ "$IMG_SERVE_CODE" == "200" ]]; then
		if diff /tmp/test-msg-image.png /tmp/served-img.png > /dev/null 2>&1; then
			echo "  ✓ Serve message image (200, binary match)"; PASS=$((PASS + 1))
		else
			echo "  ✗ Serve message image (binary mismatch)"; FAIL=$((FAIL + 1))
		fi
	else
		echo "  ✗ Serve message image (expected 200, got $IMG_SERVE_CODE)"; FAIL=$((FAIL + 1))
	fi
else
	echo "  ✗ Serve message image (no image ID)"; FAIL=$((FAIL + 1))
fi

if [[ -n "$IMG_ID" ]]; then
	check "Serve image non-member" GET "$API_URL/chats/$DCID/images/$IMG_ID" "$TOKC" "" "403"
	check "Serve image chat not found" GET "$API_URL/chats/00000000-0000-0000-0000-000000000000/images/$IMG_ID" "$TOKA" "" "404"
else
	echo "  ✗ Serve image tests skipped (no image ID)"; FAIL=$((FAIL + 1))
	echo "  ✗ Serve image tests skipped (no image ID)"; FAIL=$((FAIL + 1))
fi
check "Serve image not in chat" GET "$API_URL/chats/$DCID/images/00000000-0000-0000-0000-000000000000" "$TOKA" "" "404"

# --- Forwarding ---
echo ""
echo "=== Forwarding ==="

check "Forward text (201)" POST "$API_URL/chats/$DCID/messages/$MSG_A2/forwards" "$TOKA" \
	"{\"forwardTo\":\"$GCID\"}" "201" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); ff=d[\"forwardedFrom\"]; assert ff[\"chatId\"]==\"'"$DCID"'\"; assert ff[\"text\"] is not None"'
FWD_MSG_ID=$(python3 -c "import json; print(json.load(open('/tmp/test-resp.json'))['id'])" 2>/dev/null || true)
[[ -z "$FWD_MSG_ID" ]] && { echo "  ✗ Could not capture forwarded message ID"; FAIL=$((FAIL + 1)); }

if [[ -n "$IMG_ID" ]]; then
	check "Forward image (201)" POST "$API_URL/chats/$DCID/messages/$MSG_IMG/forwards" "$TOKA" \
		"{\"forwardTo\":\"$GCID\"}" "201" \
		'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); ff=d[\"forwardedFrom\"]; assert ff[\"msgImageId\"] is not None"'

	# Verify image is accessible in destination chat via visibility table
	IMG_VIS_CODE=$(curl -s -o /tmp/served-fwd.png -w "%{http_code}" "$API_URL/chats/$GCID/images/$IMG_ID" \
		-H "Authorization: $TOKB" --max-time 5)
	if [[ "$IMG_VIS_CODE" == "200" ]] && diff /tmp/test-msg-image.png /tmp/served-fwd.png > /dev/null 2>&1; then
		echo "  ✓ Forward image visibility (200)"; PASS=$((PASS + 1))
	else
		echo "  ✗ Forward image visibility (expected 200+match, got $IMG_VIS_CODE)"; FAIL=$((FAIL + 1))
	fi
else
	echo "  ✗ Forward image test skipped (no image ID)"; FAIL=$((FAIL + 1))
	echo "  ✗ Forward image visibility skipped (no image ID)"; FAIL=$((FAIL + 1))
fi

if [[ -n "$FWD_MSG_ID" ]]; then
	check "Forward forwarded (400)" POST "$API_URL/chats/$GCID/messages/$FWD_MSG_ID/forwards" "$TOKA" \
		"{\"forwardTo\":\"$DCID\"}" "400"
else
	echo "  ✗ Forward forwarded skipped (no FWD_MSG_ID)"; FAIL=$((FAIL + 1))
fi

check "Forward same chat (400)" POST "$API_URL/chats/$DCID/messages/$MSG_B/forwards" "$TOKA" \
	"{\"forwardTo\":\"$DCID\"}" "400"
check "Forward non-member origin (403)" POST "$API_URL/chats/$DCID/messages/$MSG_B/forwards" "$TOKC" \
	"{\"forwardTo\":\"$GCID\"}" "403"

BCID_CODE=$(curl -s -o /tmp/bc-chat.json -w "%{http_code}" -X POST "$API_URL/user/$UIDB/chats" \
	-H "Authorization: $TOKC" --max-time 5)
if [[ "$BCID_CODE" == "201" ]]; then
	BCID=$(python3 -c "import json; print(json.load(open('/tmp/bc-chat.json'))['id'])" 2>/dev/null || true)
	if [[ -n "$BCID" ]]; then
		check "Forward non-member dest (403)" POST "$API_URL/chats/$DCID/messages/$MSG_B/forwards" "$TOKA" \
			"{\"forwardTo\":\"$BCID\"}" "403"
	else
		echo "  ✗ Forward non-member dest skipped (failed to parse BCID)"; FAIL=$((FAIL + 1))
	fi
else
	echo "  ✗ Forward non-member dest skipped (failed to create B-C chat, got $BCID_CODE)"; FAIL=$((FAIL + 1))
fi

check "Forward msg not found (404)" POST "$API_URL/chats/$DCID/messages/00000000-0000-0000-0000-000000000000/forwards" "$TOKA" \
	"{\"forwardTo\":\"$GCID\"}" "404"
check "Forward dest not found (404)" POST "$API_URL/chats/$DCID/messages/$MSG_B/forwards" "$TOKA" \
	"{\"forwardTo\":\"00000000-0000-0000-0000-000000000000\"}" "404"
check "Forward origin not found (404)" POST "$API_URL/chats/00000000-0000-0000-0000-000000000000/messages/$MSG_B/forwards" "$TOKA" \
	"{\"forwardTo\":\"$GCID\"}" "404"

# --- Reactions ---
echo ""
echo "=== Reactions ==="

check "Add reaction (A 👍)" POST "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKA" \
	'{"emojiId":0}' "201" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); rl=d[\"reactionsList\"]; assert len(rl)==1, f\"expected 1, got {len(rl)}\"; assert rl[0][\"emojiId\"]==0"'
check "Replace reaction (A 🤔)" POST "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKA" \
	'{"emojiId":5}' "201" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); rl=d[\"reactionsList\"]; assert len(rl)==1, f\"expected 1, got {len(rl)}\"; assert rl[0][\"emojiId\"]==5"'
check "Second user reacts (B ❤️)" POST "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKB" \
	'{"emojiId":1}' "201" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); rl=d[\"reactionsList\"]; assert len(rl)==2, f\"expected 2, got {len(rl)}\""'
check "Invalid emojiId < 0" POST "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKA" \
	'{"emojiId":-1}' "400"
check "Invalid emojiId > 9" POST "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKA" \
	'{"emojiId":10}' "400"
check "Add reaction non-member" POST "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKC" \
	'{"emojiId":0}' "403"
check "Add reaction msg not found" POST "$API_URL/chats/$DCID/messages/00000000-0000-0000-0000-000000000000/reactions" "$TOKA" \
	'{"emojiId":0}' "404"

check "Remove reaction (A)" DELETE "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKA" "" "200" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); rl=d[\"reactionsList\"]; assert len(rl)==1, f\"expected 1 remaining, got {len(rl)}\""'
check "Remove non-existent reaction" DELETE "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKA" "" "404"
check "Remove reaction non-member" DELETE "$API_URL/chats/$DCID/messages/$MSG_A2/reactions" "$TOKC" "" "403"

# --- My Chats ---
echo ""
echo "=== My Chats ==="

check "Get my chats (A)" GET "$API_URL/chats" "$TOKA" "" "200" \
	'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); chats=d[\"chatsPreviewList\"]; assert len(chats)>=2, f\"expected >=2, got {len(chats)}\"; dm=[c for c in chats if not c[\"isGroupChat\"]][0]; assert dm[\"displayName\"]==\"Harlan Thrombey\"; assert \"lastMessage\" in dm"'

# New user with no chats
NEW_USER_CODE=$(curl -s -o /tmp/new-user.json -w "%{http_code}" -X POST "$API_URL/session" \
	-H "Content-Type: application/json" -d '{"userName":"Test EmptyChats"}' --max-time 5)
if [[ "$NEW_USER_CODE" == "200" || "$NEW_USER_CODE" == "201" ]]; then
	NEW_TOKEN=$(python3 -c "import json; print(json.load(open('/tmp/new-user.json'))['token'])" 2>/dev/null)
	check "Get my chats (empty)" GET "$API_URL/chats" "$NEW_TOKEN" "" "200" \
		'python3 -c "import json; d=json.load(open(\"/tmp/test-resp.json\")); assert d[\"chatsPreviewList\"]==[], f\"expected empty, got {d}\""'
else
	echo "  ✗ Create empty-chats user (expected 200/201, got $NEW_USER_CODE)"; FAIL=$((FAIL + 1))
fi

check "Get my chats no auth" GET "$API_URL/chats" "" "" "401"

# --- Summary ---
echo ""
echo "============================================"
echo " RESULTS: $PASS passed, $FAIL failed"
echo "============================================"

# --- Cleanup ---
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
rm -f /tmp/decaf.db /tmp/test-resp.json /tmp/test-avatar.png /tmp/test-avatar-resp.png /tmp/wasapp-server.log /tmp/msg-a1.json /tmp/msg-a2.json /tmp/msg-img.json /tmp/msg-b-reply.json /tmp/msg-reply.json /tmp/msg-reply2.json /tmp/served-img.png /tmp/served-fwd.png /tmp/new-user.json /tmp/test-msg-image.png /tmp/bc-chat.json
rm -rf uploads/

if [[ $FAIL -gt 0 ]]; then
	exit 1
fi
exit 0
