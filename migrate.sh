#!/bin/bash

MIGRATION_PATH="./migrations"
DB_DSN=""

usage() {
  echo "Usage: $0 [--create name] [--migrate] [--rollback number|all] [--goto number] --db [db_dsn]"
  echo "Flags:"
  echo "  --create [name]    : Create a migration"
  echo "  --migrate          : Run migrations up"
  echo "  --rollback [num]   : Roll back [num] migrations"
  echo "  --goto [num]       : Goto specified [num] version of migrations"
  echo "  --rollback all     : Roll back all migrations"
  echo "  --db [db_dsn]      : Specify the database DSN"
  exit 1
}

create_migration() {
  if [[ -z "$1" ]]; then
    echo "Error: Missing migration name argument."
    echo "Usage: ./script.sh --create [migration_name]"
    exit 1
  fi

  migrate create -seq -ext=.sql -dir=$MIGRATION_PATH "$1"
}

run_migrate() {
  if [[ -z "$DB_DSN" ]]; then
    echo "Error: Missing --db argument for database DSN."
    echo "Usage: ./script.sh --migrate --db [db_dsn]"
    exit 1
  fi

  migrate -path=$MIGRATION_PATH -database="$DB_DSN" up
}

rollback_migrate() {
  if [[ -z "$DB_DSN" ]]; then
    echo "Error: Missing --db argument for database DSN."
    echo "Usage: ./script.sh --rollback [num] --db [db_dsn]"
    exit 1
  fi
  if [[ -z "$1" || ! "$1" =~ ^[1-9][0-9]*$ ]]; then
    echo "Error: Invalid rollback number. Provide a positive integer (greater than 0)."
    echo "Usage: ./script.sh --rollback [number] --db [db_dsn]"
    exit 1
  fi

  migrate -path=$MIGRATION_PATH -database="$DB_DSN" down "$1"
}


goto_migrate() {
  if [[ -z "$DB_DSN" ]]; then
    echo "Error: Missing --db argument for database DSN."
    echo "Usage: ./script.sh --rollback [num] --db [db_dsn]"
    exit 1
  fi
  if [[ -z "$1" || ! "$1" =~ ^[1-9][0-9]*$ ]]; then
    echo "Error: Invalid goto number. Provide a positive integer (greater than 0)."
    echo "Usage: ./script.sh --goto [number] --db [db_dsn]"
    exit 1
  fi

  migrate -path=$MIGRATION_PATH -database="$DB_DSN" goto "$1"
}

rollback_all_migrate() {
  if [[ -z "$DB_DSN" ]]; then
    echo "Error: Missing --db argument for database DSN."
    echo "Usage: ./script.sh --rollback all --db [db_dsn]"
    exit 1
  fi

  migrate -path=$MIGRATION_PATH -database="$DB_DSN" down
}

if [[ $# -eq 0 ]]; then
  usage
fi

# First Pass: Extract --db flag
args=("$@")  # Save all original arguments
i=0

while [[ $i -lt $# ]]; do
  if [[ "${args[$i]}" == "--db" ]]; then
    if [[ -n "${args[$((i + 1))]}" ]]; then
      DB_DSN="${args[$((i + 1))]}"
      # Remove the --db and its value from the arguments array
      unset 'args[i]'
      unset 'args[i+1]'
      args=("${args[@]}")  # Rebuild the array
    else
      echo "Error: Missing database DSN after --db."
      exit 1
    fi
  fi
  ((i++))
done

# Second Pass: Process remaining arguments
set -- "${args[@]}"  # Reset the positional parameters to the updated list

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --create)
      create_migration "$2"
      shift 2
      ;;
    --migrate)
      run_migrate
      shift
      ;;
    --rollback)
      if [[ "$2" == "all" ]]; then
        rollback_all_migrate
        shift 2
      elif  [[ -z "$2" ]]; then
        rollback_migrate "1"
        shift
      else
        rollback_migrate "$2"
        shift 2
      fi
      ;;
    --goto)
      goto_migrate "$2"
      shift 2
      ;;
    *)
      usage
      ;;
  esac
done
