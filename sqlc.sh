#!/bin/bash

# Inspired by https://github.com/sqlc-dev/sqlc/issues/2469,
# this script generates an sqlc config file which generates separate packages for each query file.
# This allows for a cleaner separation between the different table queries.

base_out_dir="./repositories/database/implementation"
queries_dir="./repositories/database/queries"
migrations_dir="./cmd/migrations/schema"
sqlc_config_file="sqlc.yaml"

# the sqlc config file is overwritten with the required starting data
cat > $sqlc_config_file <<-END
version: "2"
sql:
END

# queries_config generates identical configurations for the different query files,
# with different output packages.
queries_config=$(cat <<-END
  - engine: "postgresql"
    schema: ${migrations_dir}
    queries: {query_file}
    rules:
        - sqlc/db-prepare
        - delete-without-where
        - update-without-where
    database:
        uri: 'postgresql://\${API_DB_USER}:\${API_DB_PASS}@\${API_DB_HOST}:\${API_DB_PORT}/\${API_DB_NAME}'
    gen:
      go:
        out: "{out_dir}"
        package: "{package_name}"
        sql_package: "pgx/v5"
        omit_unused_structs: true
        emit_db_tags: true
        emit_interface: true
        emit_empty_slices: true
        emit_methods_with_db_argument: true
        initialisms: []

        overrides:
          - db_type: "uuid"
            nullable: false
            go_type:
              import: "github.com/google/uuid"
              type: "UUID"

          - db_type: "uuid"
            nullable: true
            go_type:
              import: "github.com/google/uuid"
              type: "NullUUID"

          - db_type: "timestamp"
            nullable: false
            go_type:
              import: "time"
              type: "Time"

          - db_type: "timestamptz"
            nullable: false
            go_type:
              import: "time"
              type: "Time"
END
)

mkdir -p "$base_out_dir"

for file in "$queries_dir"/*; do
    # the leading path and extensions are trimmed out, leaving only the base file name
    file_name=$(basename "$file")
    file_name="${file_name%.*}"

    query_file="$file"
    package_name="$file_name"
    out_dir="$base_out_dir/$file_name"

    cumulated_data=$(echo "$queries_config" | \
        sed "s|{query_file}|$query_file|" | \
        sed "s|{package_name}|$package_name|" | \
        sed "s|{out_dir}|$out_dir|")

    echo "$cumulated_data" >> "$sqlc_config_file"
    echo  >> $sqlc_config_file   # adds a newline after each query configuration
done

cat >> $sqlc_config_file <<-END
rules:
    - name: delete-without-where
      message: "DELETE statements must have a WHERE clause"
      rule: |
        (query.sql.startsWith("DELETE") || query.sql.startsWith("delete")) &&
        !(query.sql.contains('WHERE') || query.sql.contains('where'))

    - name: update-without-where
      message: "UPDATE statements must have a WHERE clause"
      rule: |
        (query.sql.startsWith("UPDATE") || query.sql.startsWith("update")) &&
        !(query.sql.contains('WHERE') || query.sql.contains('where'))
END