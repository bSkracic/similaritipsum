#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$DB_USER" --dbname "$DB_NAME" <<-EOSQL


    CREATE TABLE IF NOT EXISTS word_entries
    (
        id serial PRIMARY KEY,
        word text,
        skewer text
    );
EOSQL