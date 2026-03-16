#!/usr/bin/env python3
"""
Simple persistent key-value store.

Commands:
SET <key> <value>
GET <key>
EXIT

Uses append-only persistence in data.db
"""

import sys

DB_FILE = "data.db"


class Database:
    def __init__(self):
        # in-memory store (list of tuples)
        self.store = []

        # ensure database file exists immediately
        with open(DB_FILE, "a"):
            pass

        self.load()

    def load(self):
        """Replay the log file to rebuild the index."""
        with open(DB_FILE, "r") as f:
            for line in f:
                parts = line.strip().split()

                if len(parts) == 3 and parts[0] == "SET":
                    key = parts[1]
                    value = parts[2]
                    self._set_memory(key, value)

    def _set_memory(self, key, value):
        """Update only the in-memory store."""
        for i in range(len(self.store)):
            if self.store[i][0] == key:
                self.store[i] = (key, value)
                return

        self.store.append((key, value))

    def set(self, key, value):
        """Store key/value and persist."""
        self._set_memory(key, value)

        with open(DB_FILE, "a") as f:
            f.write(f"SET {key} {value}\n")
            f.flush()

    def get(self, key):
        """Return value or NULL."""
        for k, v in self.store:
            if k == key:
                return v

        return "NULL"


def main():
    db = Database()

    while True:
        line = sys.stdin.readline()

        if not line:
            break

        parts = line.strip().split()

        if not parts:
            continue

        cmd = parts[0].upper()

        if cmd == "SET" and len(parts) == 3:
            db.set(parts[1], parts[2])

        elif cmd == "GET" and len(parts) == 2:
            result = db.get(parts[1])
            sys.stdout.write(result + "\n")
            sys.stdout.flush()

        elif cmd == "EXIT":
            break


if __name__ == "__main__":
    main()
