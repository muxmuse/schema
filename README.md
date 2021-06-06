# schema – a Package Manager for SQL Server

** WORK IN PROGRESS, EVERYTHING MIGHT CHANGE, NOTHING MIGHT WORK, MIGHT DESTROY YOUR DATA, ABSOLUTELY NO WARRANTIES **

## Installation

1. Install go

2. Install schema
    ``` bash
    go get https://github.com/muxmuse/schema
    ```
3. Configure your database connections
    ``` yaml
    # ~/.schemapm/config.yaml

    connections:
    - name: database_name_1
      url: host:port/database_name_1
      sendStringParametersAsUnicode: true
      prepareSQL: 2
      log: 8
      database: database_name_1
      user: db_user_1
      password: db_password_1
      selected: false # the first selected connection is used by schema
    - name: database_name_2
      url: host:port/database_name_2
      sendStringParametersAsUnicode: true
      prepareSQL: 2
      log: 16
      database: database_name_2
      user: db_user_2
      password: db_password_2
      selected: true
    ```

## Usage

List all available contexts (connections) and show which one is selected
``` bash
schema context
#   conn_1
# > conn_2
#   conn_3
```

Select a context
``` bash
schema context conn_3
#   conn_1
#   conn_2
# > conn_3
```

List installed schemas on the database
``` bash
schema list
# - v0.0.2       MY_SCHEMA_1
# - v0.3.2-beta  MY_SCHEMA_2
```

Install a schema
``` bash
schema install git@github.com:muxmuse/hello_world.schema.git v0.0.2
# Checking out https://github.com/muxmuse/hello_world.schema.git refs/tags/v0.0.2
# [schema] HELLO_WORLD v0.0.2
# [module] hello-module
# About to run migrations v0.0.0 -> v0.0.2
# - v0.0.0_v0.0.1.migrate.sql
# - v0.0.1_v0.0.2.migrate.sql
# Confirm (y/n): y
# [running]  v0.0.0_v0.0.1.migrate.sql
# [success]  v0.0.0_v0.0.1.migrate.sql
# [running]  v0.0.1_v0.0.2.migrate.sql
# [success]  v0.0.1_v0.0.2.migrate.sql
# [running] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.2/install.sql
# [success] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.2/install.sql
# 
# Successfully installed HELLO_WORLD v0.0.2
# 

schema list
# Installed schemas on ...
# - v0.0.2         HELLO_WORLD
# 
# Unmanaged schemas: dbo, ..., sys

schema show HELLO_WORLD

# name: HELLO_WORLD
# description: Demo package for schemapm
# gitTag: v0.0.2
# gitRepoUrl: git@github.com:muxmuse/hello_world.schema.git

schema uninstall HELLO_WORLD
# Checking out git@github.com:muxmuse/hello_world.schema.git refs/tags/v0.0.2
# [schema] HELLO_WORLD v0.0.2
# [module] hello-module
# [running] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.2/uninstall.sql
# [success] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.2/uninstall.sql
# About to run migrations v0.0.2 -> v0.0.0
# - v0.0.2_v0.0.1.migrate.sql
# - v0.0.1_v0.0.0.migrate.sql
# Confirm (y/n): y
# [running]  v0.0.2_v0.0.1.migrate.sql
# [success]  v0.0.2_v0.0.1.migrate.sql
# [running]  v0.0.1_v0.0.0.migrate.sql
# [success]  v0.0.1_v0.0.0.migrate.sql
# 
# Successfully removed HELLO_WORLD v0.0.2
```

## Create a package
Each schema is a folder containing 
- `schema.yaml` (see example package hello-world)
- files ending on `install.sql` to be executed on installation. At least one such file must exist and install the version function
- matching files ending on `uninstall.sql` to be executed on uninstallation
- subdirectories with the same structure (modules)
- migration files `v0.0.1_v0.0.2.migrate.sql`
- other files, which are ignored by schema.

All sql files consist of T-SQL statements, separated by `GO`.

## Whishes
1. Dependency management
2. Create a schema from scratch `schema new` or so
3. Create a schema from existing objects
    ``` bash
      schema new --from-installed 'dbo'
      # Collecting information:
      # - dbo.p_my_old_business_logic
      # - ...
      # Creating schema.yaml
      # done, created new local package dbo at version 0.0.1
      ```
4. Compare all objects in an installed schema with the repository
    ``` bash
    schema install file .
    # Collecting changes
    # - updating [my_schema].[p_my_proc_2] ...
    # - updating [my_schema].[p_my_proc_9] ...
    # Rest unchanged.
    # Setting version to 0.0.13-wip
    ```


5. Install functions before installing procedures

6. Dependency management

## Notes

Consider [mssql-scripter](https://github.com/microsoft/mssql-scripter/blob/dev/doc/installation_guide.md#linux-installation) for downloading schemas initially
