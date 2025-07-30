#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Usage: database-setup.sh <database_user>"
  exit 1
fi

for arg in "$@"; do
    echo "Argument: $arg"
done

instance_name=cloudsql-instance-database
database_name=database
database_user=$1

gcloud sql connect ${instance_name} --database ${database_name} --user ${database_user} <<EOF

\i ./database/init.sql

\q

EOF

