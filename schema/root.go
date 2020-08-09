package schema

import (
	"fmt"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/muxmuse/schema/mfa"
)

type TSchema struct {
	name string
	url string
	owner string
	version string
}

var DB *sql.DB

func getConnectedDatabase(con TConnectionConfig) (*sql.DB) {
	db, err := sql.Open(
		"sqlserver",
		mfa.Format("url", "sqlserver://{{.User}}:{{.Password}}@{{.Url}}", con))

	mfa.CatchFatal(err)
	mfa.CatchFatal(db.Ping())
	
	return db
}

func listSchemas(db *sql.DB) []TSchema {
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
		        'INFORMATION_SCHEMA'
		    )
	`)
	mfa.CatchFatal(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()

	var schemas []TSchema
	for rows.Next() {
			var schema TSchema
			rows.Scan(&schema.name, &schema.owner)
			schemas = append(schemas, schema)
	}
	mfa.CatchFatal(rows.Err())

	return schemas
}


func List(db *sql.DB) {
	schemas := listSchemas(db)
	
	fmt.Println("Installed packages:")

	var err error
	var result []TSchema
	for _, schema := range schemas {
		err = DB.QueryRow("select [" + schema.name + "].[SCHEMA___$$$]()").Scan(&schema.version)
		if err != nil {
		} else {
			result = append(result, schema)
			fmt.Println("-", schema.name, "@", schema.version)
		}
	}
}


func init() {
	var con TConnectionConfig = getSelectedConnectionConfig(getConfig())
	DB = getConnectedDatabase(con)
	fmt.Println("Using", con.Name)

	// var version string
	// mfa.CatchFatal(DB.QueryRow("select [CALDERA].[SCHEMA___$$$]()").Scan(&version))
	// log.Print("version", version)
}
