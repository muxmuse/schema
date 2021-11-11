package schema

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/muxmuse/schema/mfa"

	"gopkg.in/yaml.v2"
)

func listSchemas(db *sql.DB) ([]TSchema, []TSchema) {
	 stmt, err := db.Prepare(`
		select 
		    [schema_name] = [schema].[name], 
		    [owner_name] = [user].[name],
		    [modifiedAt] = max([object].modify_date), 
    		[hash] = CONVERT(varchar(max), CAST(HASHBYTES('SHA2_512', STRING_AGG(def, ',')) AS varbinary(max)), 1)
		from sys.schemas [schema]
		join sys.sysusers [user]
		  on [schema].[principal_id] = [user].[uid]
		join (
		    select 
		    		o.[name],
		        o.schema_id,
		        modify_date,
		        [type],
		        def = case 
		            when o.TYPE in ('C','D','P','FN','R','RF','TR','IF','TF','V') 
		            then OBJECT_DEFINITION(o.object_id)
		            -- when o.TYPE in ('P', 'RF', 'V', 'TR', 'FN', 'IF', 'TF', 'R') 
		            -- then sp_helptext(o.object_id)
		            when o.TYPE in ('T', 'U') 
		            then (select
		                *
		            from INFORMATION_SCHEMA.COLUMNS 
		            where TABLE_SCHEMA = s.name and TABLE_NAME = o.name
		            FOR XML AUTO)
		        end
		    from sys.objects o
		    join sys.schemas s on s.schema_id = o.schema_id
		) [object]
		on [schema].[schema_id] = [object].schema_id
		and [object].[name] <> 'SCHEMA_INFO'
		where [user].[name] not in (
		    'db_accessadmin',
		    'db_backupoperator',
		    'db_datareader',
		    'db_datawriter',
		    'db_ddladmin',
		    'db_denydatareader',
		    'db_denydatawriter',
		    'db_owner',
		    'db_securityadmin',
		    'INFORMATION_SCHEMA')
		group by [schema].[name], [user].[name]`)

	mfa.CatchFatal(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()

	var schemas []TSchema
	for rows.Next() {
			var schema TSchema
			rows.Scan(&schema.Name, &schema.dbOwner, &schema.modifiedAt, &schema.hash)
			schemas = append(schemas, schema)
	}
	mfa.CatchFatal(rows.Err())

	var dbname string
	mfa.CatchFatal(DB.QueryRow("select DB_NAME()").Scan(&dbname))
	if SelectedConnectionConfig.Name != dbname {
		log.Fatal("Actual database is " + dbname)
	}

	var marshalledSchema string
	var managedSchemas []TSchema
	var otherSchemas []TSchema
	for _, schema := range schemas {
		err = DB.QueryRow("select [" + schema.Name + "].[SCHEMA_INFO]()").Scan(&marshalledSchema)

		if err != nil {
			otherSchemas = append(otherSchemas, schema)
		} else {
			err = yaml.Unmarshal([]byte(marshalledSchema), &schema)

			if(err != nil) {
				fmt.Println("[WARNING]", schema.Name, " cannot be unmarshalled from yaml. Parsing SCHEMA_INFO() as Version")
				schema.GitTag = marshalledSchema
			}

			managedSchemas = append(managedSchemas, schema)
		}
	}

	return managedSchemas, otherSchemas
}


func List() {
	managedSchemas, otherSchemas := listSchemas(DB)

	fmt.Println()
	fmt.Println("Installed schemas on " + SelectedConnectionConfig.Name)
	for _, schema := range managedSchemas {
		fmt.Println("-", schema.GitTag, "\t", schema.Name)
		if len(schema.InstalledHash) > 0 && (schema.InstalledHash != schema.hash || schema.modifiedAt != schema.InstalledAt) {
			fmt.Println("  WARNING: Integrity compromised")
			fmt.Println("  Last modified at", schema.modifiedAt)
			if schema.InstalledHash != schema.hash {
				fmt.Println("  Contents has changed: hashes differ.")
			}
			fmt.Println();
		}
	}

	fmt.Println()
	fmt.Print("Unmanaged schemas: ")
	fmt.Println(mfa.Format("list", `{{ range $index, $element := .}}{{if $index}}, {{end}}{{$element.Name}}{{end}}`, otherSchemas))
}


func Show(schemaName string) {
	var marshalledSchema string
	err := DB.QueryRow("select [" + schemaName + "].[SCHEMA_INFO]()").Scan(&marshalledSchema)
	if err != nil {
		fmt.Println("[failure]", "Failed to get SCHEMA_INFO() for " + schemaName)
		fmt.Println("[failure]", err)
	} else {
		fmt.Println()
		fmt.Println(marshalledSchema)		
	}
}