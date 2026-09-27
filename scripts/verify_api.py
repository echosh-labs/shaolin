#!/usr/bin/env python3
import urllib.request
import json
import sys

def post(url, data, headers=None):
    if headers is None:
        headers = {}
    req = urllib.request.Request(
        url,
        data=json.dumps(data).encode(),
        headers={"Content-Type": "application/json", **headers}
    )
    with urllib.request.urlopen(req) as res:
        return json.loads(res.read().decode())

def get(url, headers=None):
    if headers is None:
        headers = {}
    req = urllib.request.Request(url, headers=headers)
    with urllib.request.urlopen(req) as res:
        return json.loads(res.read().decode())

def delete(url, headers=None):
    if headers is None:
        headers = {}
    req = urllib.request.Request(url, headers=headers, method="DELETE")
    with urllib.request.urlopen(req) as res:
        return res.status

def main():
    print("=== 1. Student Authentication ===")
    auth = post("http://localhost:8081/api/auth/login", {
        "email": "student@martialartsacademy.com",
        "password": "password123"
    })
    token = auth["token"]
    user = auth["user"]
    print(f"✅ Logged in: {user['first_name']} {user['last_name']} ({user['role']}, {user['current_rank']})")

    print("\n=== 2. Token Ledger Balance ===")
    bal = get("http://localhost:8081/api/tokens/balance", {"Authorization": f"Bearer {token}"})
    print(f"✅ Initial Tokens Remaining: {bal['tokens_remaining']}")

    print("\n=== 3. Classes Catalog ===")
    classes = get("http://localhost:8081/api/classes")
    print(f"✅ Loaded {len(classes)} classes. First: {classes[0]['name']}")

    print("\n=== 4. Class Occurrences ===")
    occs = get("http://localhost:8081/api/occurrences")
    # Pick a future occurrence strictly on or after current local date
    future_occ = next((o for o in occs if o["date"] >= "2026-08-24"), occs[-1])
    print(f"✅ Loaded {len(occs)} occurrences. Selected future slot: {future_occ['date']} ({future_occ['id']})")

    print("\n=== 5. Online Store Catalog ===")
    store = get("http://localhost:8081/api/store/products")
    print(f"✅ Loaded {len(store)} categories. First category: {store[0]['category']['name']} with {len(store[0]['products'])} items")

    print("\n=== 6. Guides & Manuals ===")
    guides = get("http://localhost:8081/api/guides")
    print(f"✅ Loaded {len(guides)} guides. First title: {guides[0]['title']}")

    print("\n=== 7. Forum Message Boards ===")
    boards = get("http://localhost:8081/api/forums/boards")
    print(f"✅ Loaded {len(boards)} boards. Names: {[b['name'] for b in boards]}")

    print("\n=== 8. Grading Curriculum Tracks ===")
    tracks = get("http://localhost:8081/api/grading/tracks")
    print(f"✅ Loaded {len(tracks)} tracks: {[t['name'] for t in tracks]}")

    print("\n=== 9. Admin Authentication ===")
    admin_auth = post("http://localhost:8081/api/auth/login", {
        "email": "admin@martialartsacademy.com",
        "password": "password123"
    })
    admin_user = admin_auth["user"]
    print(f"✅ Admin logged in: {admin_user['first_name']} {admin_user['last_name']} ({admin_user['role']})")

    print("\n=== 10. Class Booking Test ===")
    initial_tokens = bal['tokens_remaining']
    book_res = post("http://localhost:8081/api/bookings", {
        "occurrence_id": future_occ["id"],
        "attendance_mode": "in_person"
    }, headers={"Authorization": f"Bearer {token}"})
    book_id = book_res["id"]
    print(f"✅ Booking created: ID {book_id}")

    # Check updated balance (decreased by 1)
    bal2 = get("http://localhost:8081/api/tokens/balance", {"Authorization": f"Bearer {token}"})
    print(f"✅ Tokens Remaining after booking: {bal2['tokens_remaining']} (Expected: {initial_tokens - 1})")
    assert bal2['tokens_remaining'] == initial_tokens - 1, "Token was not deducted!"

    print("\n=== 11. Class Cancellation & Instant Refund Test ===")
    status = delete(f"http://localhost:8081/api/bookings/{book_id}", headers={"Authorization": f"Bearer {token}"})
    print(f"✅ Cancellation executed (HTTP {status})")

    # Check restored balance (increased back by 1)
    bal3 = get("http://localhost:8081/api/tokens/balance", {"Authorization": f"Bearer {token}"})
    print(f"✅ Tokens Remaining after refund: {bal3['tokens_remaining']} (Expected: {initial_tokens})")
    assert bal3['tokens_remaining'] == initial_tokens, "Token was not refunded!"

    print("\n🎉 ALL 11 END-TO-END DATABASE & API INTEGRATION TESTS PASSED PERFECTLY WITH FULL TOKEN REFUND INTEGRITY!")

if __name__ == "__main__":
    main()
