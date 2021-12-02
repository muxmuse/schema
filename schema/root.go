package schema

import (
	"fmt"
	"log"

	// "crypto/tls"
	// "net/http"
	// "gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	// githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/muxmuse/schema/mfa"

	// "github.com/gookit/color"

	"io/ioutil"

  "strings"
  
  // "errors"
)

var DB *sql.DB
var SqlServerVersion TSqlServerVersion

type TSqlServerVersion struct {
	edition string
	versionStr string
	version [4]uint
	level string
}

func getConnectedDatabase(con TConnectionConfig) (*sql.DB) {
	db, err := sql.Open(
		"sqlserver",
		mfa.Format("url", "sqlserver://{{.User}}:{{.Password}}@{{.Url}}?sendStringParametersAsUnicode={{.SendStringParametersAsUnicode}}&prepareSQL={{.PrepareSQL}}&log={{.Log}}&database={{.Database}}", con))

	mfa.CatchFatal(err)

	mfa.CatchFatal(db.Ping())
	mfa.CatchFatal(err)

	var dbname string
	mfa.CatchFatal(db.QueryRow("select DB_NAME()").Scan(&dbname))
	if SelectedConnectionConfig.Name != dbname {
		log.Print("Actual database is " + dbname + " trying use...")
		_, err = db.Exec("use " + con.Name)
		mfa.CatchFatal(err)
	}

	err = db.QueryRow(`
		select 
	    [edition] = SERVERPROPERTY('Edition'),
	    [version] = SERVERPROPERTY ('productversion'),
	    [level] = SERVERPROPERTY('ProductLevel')
	`).Scan(&SqlServerVersion.edition, &SqlServerVersion.versionStr, &SqlServerVersion.level)
	mfa.CatchFatal(err)

	_, err = fmt.Sscanf(
		SqlServerVersion.versionStr, 
		"%d.%d.%d.%d",
		&(SqlServerVersion.version)[0],
		&(SqlServerVersion.version)[1],
		&(SqlServerVersion.version)[2],
		&(SqlServerVersion.version)[3])

	fmt.Printf("Connected to %s\n%s %s\n\n", con.Url, SqlServerVersion.edition, SqlServerVersion.versionStr)

	mfa.CatchFatal(err)
	
	return db
}

/*
func schemaLocallyAvailable(schema *TSchema) bool {
	var schemaRoot string
	
	switch {
	case len(schema.localDir) > 0:
		schemaRoot = schema.localDir
	case schema.Getter == "" || schema.Getter == "file":
		switch {
		case len(schema.Url) > 0:
			schemaRoot = filepath.Join(WorkingDirectory, "schemas", schema.Url)
		case len(schema.Name) > 0:
			schemaRoot = filepath.Join(WorkingDirectory, "schemas", schema.Name)
		}
	}
	
	if _, err := os.Stat(schemaRoot); err == nil {
	  return true

	} else if os.IsNotExist(err) {
	  return false

	} else {
	  // Schrodinger: file may or may not exist. See err for details.
	  // Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
	  errors.New("Failed to determine if '" + schemaRoot + "' exists.")
	}

	return false
}
*/

// Find schema locally or download it
/*
func Get(schema *TSchema) {
	if schemaLocallyAvailable(schema) {
		return
	}

	switch schema.Getter {
	case "git":
		// Create temporary directory to clone the repository (will be moved)
		tmpPath := filepath.Join(WorkingDirectory, "schemas")
		mfa.CatchFatal(os.MkdirAll(tmpPath, os.ModePerm))
		tmpDir, err := ioutil.TempDir(tmpPath, "getting")
		mfa.CatchFatal(err)

		_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		    URL: schema.Url,
		})
		mfa.CatchFatal(err)

		// Read schema definition from downloaded files
		yamlFile, err := ioutil.ReadFile(filepath.Join(tmpDir, "schema.yaml"))
		mfa.CatchFatal(err)

		err = yaml.Unmarshal(yamlFile, &schema)
		mfa.CatchFatal(err)
		
		// Move file to schemas/{name}
		schema.localDir = filepath.Join(WorkingDirectory, "schemas", schema.Name)
		if schemaLocallyAvailable(schema) {
			defer os.RemoveAll(tmpDir)
		} else {
			mfa.CatchFatal(os.Rename(tmpDir, schema.localDir))
		}

	case "sqlserver":
		// TODO [mfa] read schema contents from database

	default:
		// Treat as local directory
		schema.localDir = schema.Url
		yamlFile, err := ioutil.ReadFile(filepath.Join(schema.localDir, "schema.yaml"))
		mfa.CatchFatal(err)

		err = yaml.Unmarshal(yamlFile, &schema)
		mfa.CatchFatal(err)
	}

	fmt.Println("[pulled] " + schema.Name + " at " + schema.localDir)
}
*/



//func Uninstall(schema *TSchema) {
	// Get(schema)
	// Scan(schema)
	/*

	rows, err := DB.Query("select [name], [type] from sys.objects where schema_id = SCHEMA_ID('" + "') and [type] in ('P', 'FN')")
	mfa.CatchFatal(err)
	defer rows.Close()

	var (
		name string
		objectType string)

	for rows.Next() {
		err := rows.Scan(&name, &objectType)
		switch objectType {
		case "P":
			mfa.CatchFatal(DB.Exec("DROP PROCEDURE " + name))
		case "FN":
			mfa.CatchFatal(DB.Exec("DROP FUNCTION " + name))
		}
	}

	// SELECT * FROM sys.objects WHERE schema_id = SCHEMA_ID('...') where type in ('FN', 'P')
	// column name:  , ...
	// column type: { FN (SQL_SCALAR_FUNCTION), U (USER_TABLE), PK (PRIMARY_KEY_CONSTRAINT), P (SQL_STORED_PROCEDURE) }
	*/

	// Remove subpackages
	// for i := 0; i < len(schema.Packages); i++ {
	// 	Uninstall(&schema.Packages[i])
	// 	fmt.Println(schema.Packages[i].Name, schema.Packages[i].localDir)
	// }

  /*
	// Remove package
	for _, script := range schema.UninstallScripts {
		fmt.Print("[running] ", script, " ...")
		err := execBatchesFromFile(filepath.Join(schema.localDir, script))
		if err != nil {
			color.Yellow.Println("[ERROR] ", err)
		}
		fmt.Println("done")
	}
	*/

	// Remove schema-info from database
	// DB.Exec(fmt.Sprintf("DROP FUNCTION [%s].[SCHEMA_INFO]()", schema.Name))
	// DB.Exec(fmt.Sprintf("DROP SCHEMA [%s]", schema.Name))
//}

type SQLError interface {
	SQLErrorNumber() int32
	SQLErrorState() uint8
	SQLErrorClass() uint8
	SQLErrorMessage() string
	SQLErrorServerName() string
	SQLErrorProcName() string
	SQLErrorLineNo() int32
}


func execBatchesFromFile(path string) (error, string, SQLError) {
	content, err := ioutil.ReadFile(path)
	mfa.CatchFatal(err)
	
	for _, batch := range strings.Split(string(content), "GO\n") {
		_, err = DB.Exec(batch)
		if err != nil {
			if sqlError, ok := err.(SQLError); ok {
				return err, batch, sqlError
			} else {
				return err, batch, nil
			}
		}
	}

	return nil, "", nil
}


/*
func Pull() {
	for _, schema := range schemas {
		err = DB.QueryRow("select [" + schema.Name + "].[SCHEMA_INFO]()").Scan(&schema.Version)
		if err != nil {
		} else {
			result = append(result, schema)
			fmt.Println("-", schema.Name, "@", schema.Version)
		}
	}

	exec sp_helptext 'dbo.proc_akquisestamm_detail'

	SELECT * FROM sys.all_objects
WHERE ([type] = 'P' OR [type] = 'X' OR [type] = 'PC')
ORDER BY [name];
go


for each object function
	EXEC sp_helptext N'..'
}
*/

func Connect() {
	DB = getConnectedDatabase(SelectedConnectionConfig)
}


// func init() {
	// Disable TLS certificate checks for 
	/*
	customClient := &http.Client {
		Transport: &http.Transport {
  		TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
  	},
  }

  client.InstallProtocol("https", githttp.NewClient(customClient))
  */

	

	// var version string
	// mfa.CatchFatal(DB.QueryRow("select ").Scan(&version))
	// log.Print("version", version)
// }


/*

pull schema from database

# assuming state
C 	CHECK_CONSTRAINT
F 	FOREIGN_KEY_CONSTRAINT
U 	USER_TABLE
IT	INTERNAL_TABLE
S 	SYSTEM_TABLE
D 	DEFAULT_CONSTRAINT
PK	PRIMARY_KEY_CONSTRAINT
TF	SQL_TABLE_VALUED_FUNCTION
TR	SQL_TRIGGER
UQ	UNIQUE_CONSTRAINT

# assuming state-less
#   .net running on sql server o_O
FS	CLR_SCALAR_FUNCTION
PC	CLR_STORED_PROCEDURE

#   dlls referenced from sql server O_o
X 	EXTENDED_STORED_PROCEDURE

#   things that I used before
P 	SQL_STORED_PROCEDURE
FN	SQL_SCALAR_FUNCTION

#   everything else
SN	SYNONYM
IF	SQL_INLINE_TABLE_VALUED_FUNCTION
SQ	SERVICE_QUEUE
AF	AGGREGATE_FUNCTION
V 	VIEW

# no idea
TT	TYPE_TABLE


USE [your_database_name_here];
GO
SELECT * FROM sys.all_objects
WHERE ([type] = 'P' OR [type] = 'X' OR [type] = 'PC')
ORDER BY [name];
go


for each object function
	EXEC sp_helptext N'...'

	or see here: https://docs.microsoft.com/en-us/sql/relational-databases/stored-procedures/view-the-definition-of-a-stored-procedure?view=sql-server-ver15


for formatting of the pulled objects, see here
	https://github.com/mjibson/sqlfmt

*/