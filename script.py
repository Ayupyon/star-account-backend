#!/usr/bin/env python3

import os
import sys


commands = {
    "init": "docker run --name star-postgres -p 15432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres",
    "createdb": "docker exec -it star-postgres createdb --username=root --owner=root star_account",
    "dropdb": "docker exec -it star-postgres dropdb star_account",
    "migrateup": "migrate -path db/migration -database \"postgresql://root:secret@localhost:15432/star_account?sslmode=disable\" -verbose up",
    "migratedown": "migrate -path db/migration -database \"postgresql://root:secret@localhost:15432/star_account?sslmode=disable\" -verbose down",
    "sqlc": "sqlc generate",
    "server": "go run main.go",
}

if __name__ == "__main__":
    succeed = True
    try:
        os.system(commands[sys.argv[1]]) == 0
    except:
        succeed = False

    if not succeed:
        print("command execution failed.")