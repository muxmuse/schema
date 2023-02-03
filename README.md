# schema – a Package Manager for SQL Server

Version control, integrity checks and package publishing for SQL Server

WARNING: MIGHT DESTROY YOUR DATA, ABSOLUTELY NO WARRANTIES, WORK IN PROGRESS -- But I use it in production

## Installation

1. Install go

2. Install schema
    ``` bash
    go install github.com/muxmuse/schema@latest
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
schema install git@github.com:muxmuse/hello_world.schema.git v0.0.3
# Checking out https://github.com/muxmuse/hello_world.schema.git refs/tags/v0.0.3
# [schema] HELLO_WORLD v0.0.3
# [module] hello-module
# About to run migrations v0.0.0 -> v0.0.3
# - v0.0.0_v0.0.1.migrate.sql
# - v0.0.1_v0.0.3.migrate.sql
# Confirm (y/n): y
# [running]  v0.0.0_v0.0.1.migrate.sql
# [success]  v0.0.0_v0.0.1.migrate.sql
# [running]  v0.0.1_v0.0.3.migrate.sql
# [success]  v0.0.1_v0.0.3.migrate.sql
# [running] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.3/install.sql
# [success] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.3/install.sql
# 
# Successfully installed HELLO_WORLD v0.0.3
# 

schema list
# Installed schemas on ...
# - v0.0.3         HELLO_WORLD
# 
# Unmanaged schemas: dbo, ..., sys

schema show HELLO_WORLD

# name: HELLO_WORLD
# description: Demo package for schemapm
# gitTag: v0.0.3
# gitRepoUrl: git@github.com:muxmuse/hello_world.schema.git

schema uninstall HELLO_WORLD
# Checking out git@github.com:muxmuse/hello_world.schema.git refs/tags/v0.0.3
# [schema] HELLO_WORLD v0.0.3
# [module] hello-module
# [running] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.3/uninstall.sql
# [success] /home/mfa/.schemapm/schemas/HELLO_WORLD-refs-tags-v0.0.3/uninstall.sql
# About to run migrations v0.0.3 -> v0.0.0
# - v0.0.3_v0.0.1.migrate.sql
# - v0.0.1_v0.0.0.migrate.sql
# Confirm (y/n): y
# [running]  v0.0.3_v0.0.1.migrate.sql
# [success]  v0.0.3_v0.0.1.migrate.sql
# [running]  v0.0.1_v0.0.0.migrate.sql
# [success]  v0.0.1_v0.0.0.migrate.sql
# 
# Successfully removed HELLO_WORLD v0.0.3
```

## Dump table data (alpha)

Dump the data of supported datatypes using SQL-Server's built-in JSON serialization.
The data of all tables will be dumped in an order which attempts to respect foreign key constraints.

``` bash
schema dump-data-json > dump.sql
```

The data can be loaded with the exec command wich reads batches from stdin

``` bash
schema exec < dump.sql
```

## Create a package

``` bash
schema create MY_SCHEMA
```

Each schema is a folder containing 
- `schema.yaml` (see example package hello-world)
- files ending on `install.sql` to be executed on installation
- matching files ending on `uninstall.sql` to be executed on uninstallation
- subdirectories with the same structure (modules)
- migration files `v0.0.1_v0.0.2.migrate.sql`
- other files, which are ignored by schema.

All sql files consist of T-SQL statements, separated by `GO`.

## Package development

To re-install the package you are currently working on, use

``` bash
schema install ./local/path
```

This will run all uninstall scripts and all install scripts from your local directory. Errors in uninstall scripts are ignored.

## Important notes regarding the SQL files

- Always create contraints with explicit names, such that contraints can be dropped in future migration scripts
    ``` sql
    -- bad
    CREATE TABLE MY_SCHEMA.T_1(
      t2_id int not null foreign key references MY_SCHEMA.T_2([id]) ON DELETE CASCADE,
      ...
    )

    -- good
    CREATE TABLE MY_SCHEMA.T_1(
      t2_id int,
      ...
    )
    GO

    ALTER TABLE MY_SCHEMA.T_1 WITH CHECK ADD CONSTRAINT fk_t_2 NOT NULL foreign key([t2_id]) references MY_SCHEMA.T_2([id]) ON DELETE CASCADE
    GO
    ```

## Whishes
1. Dependency management
2. Create a schema from existing objects
    ``` bash
      schema create --from-installed 'dbo'
      # Collecting information:
      # - dbo.p_my_old_business_logic
      # - ...
      # Creating schema.yaml
      # done, created new local package dbo at version 0.0.1
      ```
3. Compare all objects in an installed schema with the repository
    ``` bash
    schema install file .
    # Collecting changes
    # - updating [my_schema].[p_my_proc_2] ...
    # - updating [my_schema].[p_my_proc_9] ...
    # my_schema unchanged.
    # Setting version to 0.0.13-wip
    ```

4. Install functions before installing procedures

5. Data integrity check

``` sql
DBCC CHECKCONSTRAINTS WITH ALL_CONSTRAINTS;

select
    [query] = 'ALTER TABLE [' + s.name + '].[' + t.name + '] WITH CHECK CHECK CONSTRAINT [' + fk.name + ']',
    s.name as [schema],
    t.name as [table],
    fk.name as [fkc],
    fk.is_not_trusted
    from sys.foreign_keys fk
    join sys.tables t on t.object_id = fk.parent_object_id
    join sys.schemas s on t.schema_id = s.schema_id
where fk.is_not_trusted = 1
```

## Notes

Consider [mssql-scripter](https://github.com/microsoft/mssql-scripter/blob/dev/doc/installation_guide.md#linux-installation) for downloading schemas initially
