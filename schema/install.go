package schema

import (
	"fmt"
	"log"
	"golang.org/x/mod/semver"

	_ "github.com/denisenkom/go-mssqldb"

	"bufio"
  "os"

  "github.com/muxmuse/schema/mfa"
  "sort"
  "regexp"
  "gopkg.in/yaml.v3"
  "path/filepath"
  "gopkg.in/src-d/go-git.v4/plumbing"
)

func AskForConfirmation() bool {
	fmt.Print("Confirm (y/n): ")
  rune, _, err := bufio.NewReader(os.Stdin).ReadRune()
  mfa.CatchFatal(err)

  return rune == 'y' || rune == 'Y'
}


func dropSchemaInfo(schema *TSchema) {
	functionName := "[" + schema.Name + "].[SCHEMA_INFO]"
	
	// SQL Server > 2016 only
	// _, err := DB.Exec("DROP FUNCTION IF EXISTS "+ functionName)

	_, err := DB.Exec("if object_id('" + functionName + "') is not null DROP FUNCTION " + functionName)
	mfa.CatchFatal(err)
}

func createSchemaInfo(schema *TSchema) {
	marshalledSchema, err := yaml.Marshal(schema)
	mfa.CatchFatal(err)
	_, err = DB.Exec(fmt.Sprintf(`CREATE FUNCTION [%s].[SCHEMA_INFO]() RETURNS varchar(max) AS BEGIN RETURN '%s' END`, schema.Name, marshalledSchema))
	mfa.CatchFatal(err)
}

// 0.1
// 1.0
// 1.2
// 2.1
// 2.3
// 3.2
// 5.3
// 3.5
// 5.7
// 7.5
//
// 1 -> 3 : 1.2, 2.3
// 2 -> 0 : 2.1, 1.0
// 2 -> 4 : 2.3
// 6 -> 1 : 5.3, 3.2, 2.1

// s_1 := migration script start version
// s_2 := migration script end version
// c := installed version
// t := target version
//
// up: s_1 < s_2 s_2 <= t && s_1 >= c
// down: s_1 > s_2 && s_2 >= t && s_1 <= c
//
func collectMigrationScripts(fromVersion string, toVersion string, schema *TSchema) (error, []string, string) {
	if !semver.IsValid(fromVersion) || !semver.IsValid(toVersion) {
		return fmt.Errorf("Invalid version. Must be semver 2.0.0 prefixed with 'v'."), nil, ""
	}

	// -1 : migrate up
	// +1 : migrate down
	// 0  : don't migrate
	mode := semver.Compare(fromVersion, toVersion)

	var result []string
	re := regexp.MustCompile(`(?:.+/)?(v[^_]+)_(v[^_]+)\.migrate\.sql`)
	for _, path := range schema.MigrateScripts() {
		v := re.FindStringSubmatch(path)
		
		if v == nil {
			return fmt.Errorf("Unable to find version in", path), nil, ""
		}

		if mode < 0 && semver.Compare(v[1], v[2]) < 0 && semver.Compare(v[2], toVersion) <= 0 && semver.Compare(v[1], fromVersion) >= 0 {
			result = append(result, path)
		} else if mode > 0 && semver.Compare(v[1], v[2]) > 0 && semver.Compare(v[2], toVersion) >= 0 && semver.Compare(v[1], fromVersion) <= 0 {
			result = append(result, path)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		v_i := re.FindStringSubmatch(result[i])
		v_j := re.FindStringSubmatch(result[j])

		if mode < 0 {
			return semver.Compare(v_i[1], v_j[1]) < 0
		} else {
			return semver.Compare(v_i[1], v_j[1]) > 0
		}
	})

	return nil, result, fromVersion
}

func getInstalledVersion(schemaName string) *TSchema {
	managedInstalledSchemas, otherInstalledSchemas := listSchemas(DB)

	for _, s := range managedInstalledSchemas {
		if(s.Name == schemaName) {
			return &s
		}
	}

	// TODO [mfa] This check can not be guessed from function name
	for _, s := range otherInstalledSchemas {
		if(s.Name == schemaName) {
			log.Fatal("Schema ", schemaName, " already exists in the database but is not managed by schemapm.", "See schema create --from-db-schema on how to manage a unmanaged schema")
		}
	}

	return nil
}

func migrate(fromVersion string, toVersion string, schema *TSchema) {
	err, migrationScripts, fromVersion := collectMigrationScripts(fromVersion, toVersion, schema)
	mfa.CatchFatal(err)

	if(len(migrationScripts) >  0) {
		fmt.Println("About to run migrations", fromVersion, "->", toVersion)
		for _, path := range migrationScripts {
			fmt.Println("-", path)
		}

		if(AskForConfirmation()) {
			for _, script := range migrationScripts {
				fmt.Println("[running] ", script)
				mfa.CatchFatal(execBatchesFromFile(filepath.Join(schema.localDir, script)))
				fmt.Println("[success] ", script)
			}

		} else {
			log.Fatal("Aborted by user")
		}
	}
}

func runScriptsOrRollBack(scripts [][2]string) error {	
	var err error
	var index int
	for index = range scripts {
		fmt.Println("[running]", scripts[index][0])
		scriptErr := execBatchesFromFile(scripts[index][0])
		if scriptErr == nil {
			fmt.Println("[success]", scripts[index][0])
		} else {
			fmt.Println("[failure]", scripts[index][0], "executed with errors")
			fmt.Println("[failure]", scriptErr)
			err = scriptErr
			break
		}
	}

	if err != nil {
		fmt.Println("\nAttempting to roll back (ignoring errors)")

		for i := index; i >= 0; i = i -1 {
			fmt.Println("[running]", scripts[i][1])
			scriptErr := execBatchesFromFile(scripts[i][1])
			if scriptErr != nil {
				fmt.Println("[warning]", scripts[i][1], "executed with errors")
				fmt.Println("[warning]", scriptErr)
			}
		}
	}

	return err
}

func Install(schemaToInstall *TSchema) {
	installedSchema := getInstalledVersion(schemaToInstall.Name)
	installedVersion := "v0.0.0"
	
	if installedSchema != nil && installedSchema.GitTag == schemaToInstall.GitTag && !schemaToInstall.devMode {
		fmt.Println(schemaToInstall.Name, schemaToInstall.GitTag, "already installed.")
		return
	}

	if installedSchema != nil {
		installedVersion = installedSchema.GitTag
		if !schemaToInstall.devMode {
			err, _ := Checkout(installedSchema.GitRepoUrl, plumbing.NewTagReferenceName(installedSchema.GitTag))
			mfa.CatchFatal(err)
			mfa.CatchFatal( runScriptsOrRollBack(installedSchema.UninstallScripts()) )
		} else {
			mfa.CatchFatal( runScriptsOrRollBack(schemaToInstall.UninstallScripts()) )	
		}
	}

	migrate(installedVersion, schemaToInstall.GitTag, schemaToInstall)
	mfa.CatchFatal( runScriptsOrRollBack(schemaToInstall.InstallScripts()) )

	dropSchemaInfo(schemaToInstall)
	createSchemaInfo(schemaToInstall)

	fmt.Println()
	if schemaToInstall.devMode {
		fmt.Println("Successfully updated", schemaToInstall.Name, schemaToInstall.GitTag, "(in development)")
	} else {
		fmt.Println("Successfully installed", schemaToInstall.Name, schemaToInstall.GitTag)
	}
	fmt.Println()
}

func Uninstall(schemaName string) {
	installedSchema := getInstalledVersion(schemaName)
	
	if installedSchema == nil {
		fmt.Println(schemaName, "does not seem to be installed.")
		return
	}

	err, installedSchema := Checkout(installedSchema.GitRepoUrl, plumbing.NewTagReferenceName(installedSchema.GitTag))
	mfa.CatchFatal(err)
	mfa.CatchFatal( runScriptsOrRollBack(installedSchema.UninstallScripts()) )

	dropSchemaInfo(installedSchema)
	migrate(installedSchema.GitTag, "v0.0.0", installedSchema)
	
	fmt.Println()
	fmt.Println("Successfully removed", installedSchema.Name, installedSchema.GitTag)
	fmt.Println()	
}