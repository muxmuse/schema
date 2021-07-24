package schema

import (
  "fmt"
  //"log"

  // "crypto/tls"
  // "net/http"
  // "gopkg.in/src-d/go-git.v4/plumbing/transport/client"
  // githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"

  "github.com/muxmuse/schema/mfa"

  // "github.com/gookit/color"

  "gopkg.in/yaml.v3"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"

  "os"
  "io/ioutil"
  "path/filepath"

  "strings"
)

type TModule struct {
  Name string
  Dependencies []string
  installScripts [][2]string
  migrateScripts []string

  modules []TModule

  schema *TSchema
}

func (self *TModule) LocalDir() string {
  return filepath.Join(self.schema.localDir, self.Name)
} 

type TSchema struct {
  Name string
  Description string
  GitTag string `yaml:"gitTag"`
  GitRepoUrl string `yaml:"gitRepoUrl"`

  // Checked out schemas
  devMode bool
  localDir string
  rootModule TModule

  // Installed schemas
  dbOwner string
}

func (self *TSchema) LocalDir() string {
  return self.localDir
}

func (self *TModule) InstallScripts() [][2]string {
   var result [][2]string

  for _, s := range self.installScripts {
    do := filepath.Join(self.LocalDir(), s[0])
    undo := filepath.Join(self.LocalDir(), s[1])
    result = append(result, [2]string{do, undo})
  }

  return result
}

func (self *TModule) UninstallScripts() [][2]string {
   var result [][2]string

  for _, s := range self.installScripts {
    do := filepath.Join(self.LocalDir(), s[1])
    undo := filepath.Join(self.LocalDir(), s[0])
    result = append(result, [2]string{do, undo})
  }

  return result
}

func (self *TSchema) InstallScripts() [][2]string {
  result := self.rootModule.InstallScripts()

  for _, m := range self.rootModule.modules {
    result = append(result, m.InstallScripts()...)
  }

  return result
}

func (self *TSchema) UninstallScripts() [][2]string {
  result := self.rootModule.UninstallScripts()
  
  for _, m := range self.rootModule.modules {
    result = append(result, m.UninstallScripts()...)
  }

  return result
}

func (self *TModule) MigrateScripts() []string {
   var result []string

  for _, s := range self.migrateScripts {
    result = append(result, filepath.Join(self.LocalDir(), s))
  }

  return result
}

func (self *TSchema) MigrateScripts() []string {
  result := self.rootModule.MigrateScripts()
  
  for _, m := range self.rootModule.modules {
    result = append(result, m.MigrateScripts()...)
  }

  return result
}

func CheckoutDev(localDir string) (error, *TSchema) {
  var schema TSchema
  
  yamlFile, err := ioutil.ReadFile(filepath.Join(localDir, "schema.yaml"))
  if err == nil {
    err = yaml.Unmarshal(yamlFile, &schema)  
    if err == nil {
      schema.localDir = localDir
      schema.rootModule.schema = &schema
      schema.devMode = true
      fmt.Println("[schema]", schema.Name, schema.GitTag, "(in development)")
      scanModule(&schema.rootModule)
      return nil, &schema
    }
  }
  
  if err != nil {
    return err, nil
  } else {
    return nil, &schema
  }
}


func Checkout(gitRepoUrl string, gitReferenceName plumbing.ReferenceName) (error, *TSchema) {
  // Create temporary directory to clone the repository (will be moved)
  tmpDir, err := ioutil.TempDir(SchemasDir, "getting")
  mfa.CatchFatal(err)

  fmt.Println("Checking out", gitRepoUrl, string(gitReferenceName))

  _, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
      URL: gitRepoUrl,
      ReferenceName: gitReferenceName,
      SingleBranch: true })

  var schema TSchema
  if err == nil {
    yamlFile, err := ioutil.ReadFile(filepath.Join(tmpDir, "schema.yaml"))
    if err == nil {
      err = yaml.Unmarshal(yamlFile, &schema)  
      if err == nil {
        schema.localDir = filepath.Join(SchemasDir, strings.ReplaceAll(schema.Name + "-" + string(gitReferenceName), "/", "-"))
        schema.rootModule.schema = &schema
        schema.devMode = false
        fmt.Println("[schema]", schema.Name, schema.GitTag)
        mfa.CatchFatal(os.Rename(tmpDir, schema.localDir))
        scanModule(&schema.rootModule)
        return nil, &schema
      }
    }
  }
  
  if err != nil {
    defer os.RemoveAll(tmpDir)
    return err, nil
  } else {
    return nil, &schema
  }
}

func scanModule(module *TModule) {
  fileInfos, err := ioutil.ReadDir(module.LocalDir())
  mfa.CatchFatal(err)

  for _, fileInfo := range fileInfos {
    switch {

    case strings.HasSuffix(fileInfo.Name(), "uninstall.sql"):
      a := fileInfo.Name()
      b := a[:len(a)-len("uninstall.sql")] + "install.sql"
      module.installScripts = append(module.installScripts, [2]string{ b, a })
    
    case strings.HasSuffix(fileInfo.Name(), "migrate.sql"):
      module.migrateScripts = append(module.migrateScripts, fileInfo.Name())
    
    case fileInfo.IsDir():
      if _, err := os.Stat(filepath.Join(module.LocalDir(), fileInfo.Name(), "schema.yaml")); err == nil {
        var subModule TModule
        subModule.schema = module.schema
        subModule.Name = fileInfo.Name()
        scanModule(&subModule)
        
        module.modules = append(module.modules, subModule)
        fmt.Println("[module]", subModule.Name)
      }
    }
  }
}
