#!/bin/sh -e

table="stylusup_migrations"

dbmate -d migrations --migrations-table "$table" -u "$1" up
