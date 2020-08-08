# SchmaPM -- A Package Manager for SQL Server

``` yaml
# ~/.schema

connections:
  mydb:
    host:
    db:
    user:
    password:

```

``` bash
schema use mydb
# > Using connection 'mydb'

# Install from git repository
schema install https://gitea.my-company.com/my-user/my-package.schema/src/tag/0.0.1
# Install from local directory
schema install ms-scs-rest.schema
schema uninstall ms-scs-rest.schema

schema list
# Installed packages on mydb
# - my-package 0.0.1
# - my-package$authorization 0.0.1

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
