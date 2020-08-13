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

	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
  "path/filepath"

  "gopkg.in/src-d/go-git.v4"

  "strings"
  "errors"
)

type TSchema struct {
	Name string
	Description string
	Version string

	InstallScripts []string
	UninstallScripts []string
	Packages []TSchema

	Url string
	Getter string
	Dependencies []TSchema

	Owner string
	Dir string
}

var DB *sql.DB

func getConnectedDatabase(con TConnectionConfig) (*sql.DB) {
	db, err := sql.Open(
		"sqlserver",
		mfa.Format("url", "sqlserver://{{.User}}:{{.Password}}@{{.Url}}", con))

	mfa.CatchFatal(err)

	mfa.CatchFatal(db.Ping())
	mfa.CatchFatal(err)

	_, err = db.Exec("use " + con.Name)
	mfa.CatchFatal(err)
	
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
		        'INFORMATION_SCHEMA')`)

	mfa.CatchFatal(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()

	var schemas []TSchema
	for rows.Next() {
			var schema TSchema
			rows.Scan(&schema.Name, &schema.Owner)
			schemas = append(schemas, schema)
	}
	mfa.CatchFatal(rows.Err())

	return schemas
}


func List() {
	schemas := listSchemas(DB)
	
	var dbname string
	mfa.CatchFatal(DB.QueryRow("select DB_NAME()").Scan(&dbname))
	if SelectedConnectionConfig.Name != dbname {
		log.Fatal("Actual database is " + dbname)
	}

	fmt.Println("Installed packages on " + SelectedConnectionConfig.Name)

	var err error
	var result []TSchema
	for _, schema := range schemas {
		err = DB.QueryRow("select [" + schema.Name + "].[SCHEMA_INFO]()").Scan(&schema.Version)
		if err != nil {
		} else {
			result = append(result, schema)
			fmt.Println("-", schema.Name, "@", schema.Version)
		}
	}
}

func schemaLocallyAvailable(schema *TSchema) bool {
	var schemaRoot string
	
	switch {
	case len(schema.Dir) > 0:
		schemaRoot = schema.Dir
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

// Find schema locally or download it
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
		schema.Dir = filepath.Join(WorkingDirectory, "schemas", schema.Name)
		if schemaLocallyAvailable(schema) {
			defer os.RemoveAll(tmpDir)
		} else {
			mfa.CatchFatal(os.Rename(tmpDir, schema.Dir))
		}

	case "sqlserver":
		// TODO [mfa] read schema contents from dabatase

	default:
		// Treat as local directory
		schema.Dir = schema.Url
		yamlFile, err := ioutil.ReadFile(filepath.Join(schema.Dir, "schema.yaml"))
		mfa.CatchFatal(err)

		err = yaml.Unmarshal(yamlFile, &schema)
		mfa.CatchFatal(err)
	}

	fmt.Println("[pulled] " + schema.Name + " at " + schema.Dir)
}

func Scan(schema *TSchema) {
	// Scan package contents
	yamlFile, err := ioutil.ReadFile(filepath.Join(schema.Dir, "schema.yaml"))
	mfa.CatchFatal(err)

	err = yaml.Unmarshal(yamlFile, &schema)
	mfa.CatchFatal(err)

	fileInfos, err := ioutil.ReadDir(schema.Dir)
	mfa.CatchFatal(err)

	for _, fileInfo := range fileInfos {
		switch {
		
		case strings.HasSuffix(fileInfo.Name(), "uninstall.sql"):
			schema.UninstallScripts = append(schema.UninstallScripts, fileInfo.Name())
		
		case strings.HasSuffix(fileInfo.Name(), "install.sql"):
			schema.InstallScripts = append(schema.InstallScripts, fileInfo.Name())
		
		case fileInfo.IsDir():
			if _, err := os.Stat(filepath.Join(schema.Dir, fileInfo.Name(), "schema.yaml")); err == nil {
			  var subPackage TSchema
				subPackage.Dir = filepath.Join(schema.Dir, fileInfo.Name())
				schema.Packages = append(schema.Packages, subPackage)

				fmt.Println("[found sub-package]", subPackage.Dir)
				Scan(&subPackage)
			} else {
				fmt.Println("ignoring ", fileInfo.Name())
			}
		}
	}
}

func Uninstall(schema *TSchema) {
	Get(schema)
	Scan(schema)

	// Remove subpackages
	for i := 0; i < len(schema.Packages); i++ {
		Uninstall(&schema.Packages[i])
		fmt.Println(schema.Packages[i].Name, schema.Packages[i].Dir)
	}

	// Remove package
	for _, script := range schema.UninstallScripts {
		fmt.Print("[running] ", script, " ...")
		execBatchesFromFile(filepath.Join(schema.Dir, script))
		fmt.Println("done")
	}
}

func Install(schema *TSchema) {
	Get(schema)
	Scan(schema)

	// Install dependencies
	fmt.Println("\nInstalling dependencies for", schema.Name)
	for i := 0; i < len(schema.Dependencies); i++ {
		// Install(&schema.Dependencies[i])
		fmt.Println(schema.Dependencies[i].Name, schema.Dependencies[i].Dir)
	}
	
	// Install package
	for _, script := range schema.InstallScripts {
		fmt.Print("[running] ", script, " ...")
		execBatchesFromFile(filepath.Join(schema.Dir, script))
		fmt.Println("done")
	}

	// Install subpackages
	for i := 0; i < len(schema.Packages); i++ {
		Install(&schema.Packages[i])
		fmt.Println(schema.Packages[i].Name, schema.Packages[i].Dir)
	}
}

func execBatchesFromFile(path string) {
	content, err := ioutil.ReadFile(path)
	mfa.CatchFatal(err)

	for _, batch := range strings.Split(string(content), "GO\n") {
		_, err = DB.Exec(batch)
		mfa.CatchFatal(err)
	}
}

func Connect() {
	DB = getConnectedDatabase(SelectedConnectionConfig)
}


func init() {
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
	// mfa.CatchFatal(DB.QueryRow("select [CALDERA].[SCHEMA___$$$]()").Scan(&version))
	// log.Print("version", version)
}
