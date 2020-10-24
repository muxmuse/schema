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
      url: host:port/database_name_1?sendStringParametersAsUnicode=true&prepareSQL=2&log=8&database=database_name_1
      user: db_user_1
      password: db_password_1
      selected: false # the first selected connection is used by schema
    - name: database_name_2
      url: host:port/database_name_2?sendStringParametersAsUnicode=true&prepareSQL=2&log=16&database=database_name_2
      user: db_user_2
      password: db_password_2
      selected: true
    ```

## Use schema with an existing package

In the folder containing `schema.yaml`

``` bash
schema use database_name_1
# database_name_1 true
# database_name_2 false
schema list
# - MY_SCHEMA_1 @ 0.0.2
# - MY_SCHEMA_2 @ 0.3.2-beta
```
If your schema contains migrations scripts named like `.state.initial_0.0.1.sql`, then you can run them manually to apply changes to the data in the schema.

Afterwards:
``` bash
# run all install scripts
schema install file .

# run all uninstall scripts
schema uninstall file .
```

## Create a package
Each schema is a folder containing 
- `schema.yaml`
    ``` yaml
    # schema.yaml

    name: MY_SCHEMA
    description: A test schema
    dependencies:
    - url: ssh://git@gitea.mycompany.com/my_gitea_user/my_repo.git
      getter: git
      name: my_repo
    - url: https://github.com/my_user/my_repo
      getter: git
      name: my_repo
    - url: schemas/local_dependency
    ```
- files ending on `install.sql` to be executed on installation. At least one such file must exist and install the version function
    ``` mssql
    CREATE FUNCTION [MY_SCHEMA].[SCHEMA_INFO]()
    RETURNS varchar(100)
    AS
    BEGIN
        RETURN '0.1.17'
    END
    GO
    ```
- other files ending on `uninstall.sql` to be executed on uninstallation
- subdirectories with the same structure (sub-schemas)
- migration files name ending on something like `state.v1_v2.sql` to be executed manually by the user to migrate between version. At least on of those files must exist, ending on `.state.initial_(...)` that creates the schema.
- other files, which are ignored by schema.

All sql files feature T-SQL statements, separated by `GO`.

## Whishes

1. Run migration scripts using schema
    ``` bash
    schema migrate my-package --to 0.0.4
    # Migrating package my-package on connection mydb
    # Installed version is 0.0.1
    #  
    # running migration scripts
    # - 0.0.2.upgrade-to.sql
    # - 0.0.3.upgrade-to.sql
    # - 0.0.4.upgrade-to.sql
    # 
    # done, my-packge is now on 0.0.4

    schema migrate my-package --to 0.0.2
    # Migrating package my-package on connection mydb
    # Installed version is 0.0.1
    #  
    # running migration scripts
    # - 0.0.3.downgrade-to.sql
    # - 0.0.2.downgrade-to.sql
    # 
    # done, my-packge is now on 0.0.2
    ```
2. Create database schema and version information automatically, without specifying in SQL script

3. Compare all objects in an installed schema with the repository
    ``` bash
    schema install file .
    # Collecting changes
    # - updating [my_schema].[p_my_proc_2] ...
    # - updating [my_schema].[p_my_proc_9] ...
    # Rest unchanged.
    # Setting version to 0.0.13-wip
    ```

4. Create a new schema by pulling existing objects from a database
    ``` bash
    schema create --from-installed 'dbo'
    # Collecting information:
    # - dbo.p_my_old_business_logic
    # - ...
    # Creating schema.yaml
    # done, created new local package dbo at version 0.0.1
    ```

5. Install functions before installing procedures

## Notes

Consider [mssql-scripter](https://github.com/microsoft/mssql-scripter/blob/dev/doc/installation_guide.md#linux-installation) for downloading schemas initially