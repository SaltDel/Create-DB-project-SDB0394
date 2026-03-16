#!/usr/bin/env python3
"""
Simple Persistent Key-Value Database

Supported commands (read from STDIN):

SET <key> <value>
GET <key>
EXIT

Design:
- Uses an append-only log file called data.db
- On startup, the program replays the log to rebuild the in-memory index
- The in-memory index is implemented as a list (no dictionaries allowed)
- Last write wins semantics
"""

import sys
from typing import List, Tuple, Optional

DB_FILE: str = "data.db"


class Database:
    """Simple key-value database using a list as the index."""

    def __init__(self) -> None:
        # In-memory store: list of (key, value) pairs
        self.store: List[Tuple[str, str]] = []

        # Ensure the database file exists
        with open(DB_FILE, "a"):
            pass

        # Load existing data
        self.load()

    def load(self) -> None:
        """Replay the append-only log file to rebuild memory."""
        try:
            with open(DB_FILE, "r") as f:
                for line in f:
                    parts = line.strip().split()

                    if len(parts) == 3 and parts[0] == "SET":
                        key = parts[1]
                        value = parts[2]
                        self._set_memory(key, value)

        except OSError:
            # If file can't be read, start with empty store
            pass

    def _set_memory(self, key: str, value: str) -> None:
        """Update the in-memory index only."""
        for i in range(len(self.store)):
            if self.store[i][0] == key:
                self.store[i] = (key, value)
                return

        self.store.append((key, value))

    def set(self, key: str, value: str) -> None:
        """Store a key-value pair and persist to disk."""
        self._set_memory(key, value)

        try:
            with open(DB_FILE, "a") as f:
                f.write(f"SET {key} {value}\n")
                f.flush()
        except OSError:
            pass

    def get(self, key: str) -> Optional[str]:
        """Retrieve value for key or None if not found."""
        for k, v in self.store:
            if k == key:
                return v
        return None


def main() -> None:
    """Main command loop reading commands from STDIN."""
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

            if result is not None:
                sys.stdout.write(result + "\n")
            else:
                sys.stdout.write("\n")

            sys.stdout.flush()

        elif cmd == "EXIT":
            break


if __name__ == "__main__":
    main()
