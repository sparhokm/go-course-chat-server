#!/bin/bash

goose -dir "./migrations" postgres "${MIGRATION_DSN}" up -v