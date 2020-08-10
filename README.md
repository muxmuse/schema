# SchmaPM -- A Package Manager for SQL Server

## Installation

1. Install go

2. Install schema
    ``` bash
    go get https://github.com/muxmuse/schema
    ```

## Create a package
Each schema is a folder containing 
- `schema.yaml` (see below)
- files ending on `uninstall.sql` to be executed on installation
- other files ending on `install.sql` to be executed on uninstallation

All sql files feature T-SQL statements, separated by `GO`.

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

## Configure database connections
``` yaml
# ~/.schemapm/config.yaml

connections:
- name: database_name_1
  url: db05:1433/database_name_1?sendStringParametersAsUnicode=true&prepareSQL=2&database=database_name_1
  user: db_user_1
  password: db_password_1
  selected: false # the first selected connection is used by schema

- name: database_name_2
  url: db05:1433/database_name_2?sendStringParametersAsUnicode=true&prepareSQL=2&log=16&database=database_name_2
  user: db_user_2
  password: db_password_2
  selected: true
```

## Use schema
```bash
schema install git https://github.com/my_user/my_repo.schema

# nesting in schemas/ is convenction
schema install file ./schemas/my_directory

schema uninstall file ./schemas/my_directory
```



## Whishes
``` bash
schema use mydb
# > Using connection 'mydb'

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
