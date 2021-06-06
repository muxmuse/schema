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
		    [owner_name] = [user].[name]
		from 
		    sys.schemas [schema]
		    join
		    sys.sysusers [user]
		    on [schema].[principal_id] = [user].[uid]
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
		        'INFORMATION_SCHEMA')`)

	mfa.CatchFatal(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()

	var schemas []TSchema
	for rows.Next() {
			var schema TSchema
			rows.Scan(&schema.Name, &schema.dbOwner)
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