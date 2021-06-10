package schema

import (
	// "bufio"
  // "fmt"
  "io/ioutil"
  "os"
  "github.com/muxmuse/schema/mfa"
  "path/filepath"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing/object"
)

func CreateNew(name string, path string, version string) {
	readme := []byte(`# ` + name + `

This is a [schema pm](https://github.com/muxmuse/schema) package

Start by editing schema.yaml then select a database and install your schema
` + "``` bash" + `
schema context <your-database>
schema install .
` + "```" + `

Files in top-level sub-direcory containing a file called schema.yaml are also considered in installation.
This shall help to structure the code.
`)

	schemaYaml  := []byte(`kind: Schema
# This is the name of the database schema
name: ` + name + `
description: Describe here what functionality the package provides
# All releases must be tagged with v + semver
gitTag: ` + version + `
# Please set the url to a location reachable for all users of this schema
gitRepoUrl: ` + path + `
`)
	migrate_0_1 := []byte(`-- Migrate mutable parts of ` + name + `
--
-- Schema will ask the user for confirmation before running migrations. 
--
-- For each vA.B.C_vX.Y.Z.migrate.sql there must be the opposite script vX.Y.Z_vA.B.C.migrate.sql
-- Immutable parts of the schema should be defined in files ending on install.sql
`)

	migrate_1_0 := []byte(`-- Migrate mutable parts of ` + name + `
--
-- Schema will ask the user for confirmation before running migrations. 
--
-- For each vA.B.C_vX.Y.Z.migrate.sql there must be the opposite script vX.Y.Z_vA.B.C.migrate.sql
-- Immutable parts of the schema should be defined in files ending on install.sql
`)

	install := []byte(`-- Install immutable parts of ` + name +`
--
-- The file can include multiple batches separated with GO
--
-- For each file ending on install.sql there must be file ending on uninstall.sql
-- To delete data or alter tables, use files ending migrate.sql
--
CREATE FUNCTION [` + name + `].HELLO_WORLD()
	RETURNS varchar(max)
	AS BEGIN
		RETURN 'Enter all the immutable code in install.sql files'
	END
`)
	uninstall := []byte(`-- Uninstall immutable parts of ` + name +`
--
-- The file can include multiple batches separated with GO
--
-- For each file ending on install.sql there must be file ending on uninstall.sql
-- To delete data or alter tables, use files ending migrate.sql
--
DROP FUNCTION [` + name + `].HELLO_WORLD
`)

	os.MkdirAll(path, os.ModePerm)
	mfa.CatchFatal(ioutil.WriteFile(filepath.Join(path, "schema.yaml"), schemaYaml, 0644))
	mfa.CatchFatal(ioutil.WriteFile(filepath.Join(path, "v0.0.0_" + version + ".migrate.sql"), migrate_0_1, 0644))
	mfa.CatchFatal(ioutil.WriteFile(filepath.Join(path, version + "_v0.0.0.migrate.sql"), migrate_1_0, 0644))
	mfa.CatchFatal(ioutil.WriteFile(filepath.Join(path, "install.sql"), install, 0644))
	mfa.CatchFatal(ioutil.WriteFile(filepath.Join(path, "uninstall.sql"), uninstall, 0644))
	mfa.CatchFatal(ioutil.WriteFile(filepath.Join(path, "README.md"), readme, 0644))

	repo, err := git.PlainInit(path, false)
	mfa.CatchFatal(err)

	sig := &object.Signature{
		Email: "",
		Name: "schemapm",
	}

	w, err := repo.Worktree()
	mfa.CatchFatal(err)
	w.Add(filepath.Join("schema.yaml"))
  w.Add(filepath.Join("v0.0.0_" + version + ".migrate.sql"))
  w.Add(filepath.Join(version + "_v0.0.0.migrate.sql"))
  w.Add(filepath.Join("install.sql"))
  w.Add(filepath.Join("uninstall.sql"))
  w.Add(filepath.Join("README.md"))
  hash, err := w.Commit("initial", &git.CommitOptions{
  	Author: sig,
  })
  mfa.CatchFatal(err)
  _, err = repo.CreateTag(version, hash, &git.CreateTagOptions{
  	Tagger: sig,
  	Message: version,
  })
  mfa.CatchFatal(err)
}
