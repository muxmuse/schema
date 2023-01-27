package schema

/*
TODO [mfa] consider standard text representation for geography and geometry
DECLARE @g geography;  
SET @g = geography::STGeomFromText('LINESTRING(-122.360 47.656, -122.343 47.656)', 4326);  
SELECT @g.ToString();
*/

/*
T-SQL Datatypes
- Exact numerics
  - bigint
  - numeric
  - bit
  - smallint
  - decimal
  - smallmoney
  - int
  - tinyint
  - money

- Approximate numerics
	- float
	- real

- Date and time
  - date
  - datetimeoffset
  - datetime2
  - smalldatetime
  - datetime
  - time

- Binary strings
	- binary
	- varbinary
	- image

- Character strings
	- char
	- varchar
	- text

- Unicode character strings
	- nchar
	- nvarchar
	- ntext

- Other data types
	- cursor
	- rowversion
	- hierarchyid
	- uniqueidentifier
	- sql_variant
	- xml
	- Spatial Geometry Types
	- Spatial Geography Types
	- table

*/


import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/muxmuse/schema/mfa"
	"strings"
	"errors"
	"log"
	"fmt"
)

type TColumn struct {
	Name string
	Type string 
}

type TTable struct {
	Schema string
	Name string
	Dependencies string
	Level string
	Columns []TColumn
}

func (c *TColumn) FqName() (string) {
	return "[" + c.Name + "]"
}

func (c *TColumn) DumpStatement() (string, error) {
	if c.Type == "geometry" || 
	   c.Type == "geography" {
		return c.FqName() + " = cast(" + c.FqName() + " as varbinary(max))", nil
	}

	if c.Type == "bigint" ||
	   c.Type == "numeric" ||
	   c.Type == "bit" ||
	   c.Type == "smallint" ||
	   c.Type == "decimal" ||
	   c.Type == "smallmoney" ||
	   c.Type == "int" ||
	   c.Type == "tinyint" ||
	   c.Type == "money" ||
	   c.Type == "float" ||
	   c.Type == "real" ||
	   c.Type == "date" ||
	   c.Type == "datetimeoffset" ||
	   c.Type == "datetime2" ||
	   c.Type == "smalldatetime" ||
	   c.Type == "datetime" ||
	   c.Type == "time" ||
	   c.Type == "binary" ||
	   c.Type == "varbinary" ||
	   c.Type == "char" ||
	   c.Type == "varchar" ||
	   c.Type == "nchar" ||
	   c.Type == "nvarchar" ||
	   c.Type == "uniqueidentifier" ||
	   c.Type == "xml" {
		return c.FqName(), nil
	}

	if c.Type == "image" || 
	   c.Type == "text" ||
	   c.Type == "ntext" ||
	   c.Type == "hierarchyid" ||
	   c.Type == "sql_variant" {
	   return "", errors.New("Datatype not supported for json dumps: " + c.Type)
	}

	if c.Type == "rowversion" ||
	   c.Type == "timestamp" {
	   log.Println("Column of datatype " + c.Type + " will be ignored during dump")
	   return "", nil
	}

	return "", errors.New("Unknown datatype: " + c.Type)
}

func (c *TColumn) WithType() (string) {
	if c.Type == "binary" ||
	   c.Type == "varbinary" {
		return "varbinary(max)"
	}

	if c.Type == "char" ||
	   c.Type == "varchar" {
		return "varchar(max)"
	}

	if c.Type == "nchar" ||
	   c.Type == "nvarchar" {
		return "nvarchar(max)"
	}

	return c.Type
}

func (c *TColumn) LoadStatement() (string, error) {
	if c.Type == "geometry" || 
	   c.Type == "geography" {
		return c.FqName() + " = cast(" + c.FqName() + " as " + c.Type  + ")", nil
	}

	if c.Type == "bigint" ||
	   c.Type == "numeric" ||
	   c.Type == "bit" ||
	   c.Type == "smallint" ||
	   c.Type == "decimal" ||
	   c.Type == "smallmoney" ||
	   c.Type == "int" ||
	   c.Type == "tinyint" ||
	   c.Type == "money" ||
	   c.Type == "float" ||
	   c.Type == "real" ||
	   c.Type == "date" ||
	   c.Type == "datetimeoffset" ||
	   c.Type == "datetime2" ||
	   c.Type == "smalldatetime" ||
	   c.Type == "datetime" ||
	   c.Type == "time" ||
	   c.Type == "binary" ||
	   c.Type == "varbinary" ||
	   c.Type == "char" ||
	   c.Type == "varchar" ||
	   c.Type == "nchar" ||
	   c.Type == "nvarchar" ||
	   c.Type == "uniqueidentifier" ||
	   c.Type == "xml" {
		return c.FqName(), nil
	}

	if c.Type == "image" || 
	   c.Type == "text" ||
	   c.Type == "ntext" ||
	   c.Type == "hierarchyid" ||
	   c.Type == "sql_variant" {
	   return "", errors.New("Datatype not supported for json dumps: " + c.Type)
	}

	if c.Type == "rowversion" ||
	   c.Type == "timestamp" {
	   log.Println("Column of datatype " + c.Type + " will be ignored during dump")
	   return "", nil
	}

	return "", errors.New("Unknown datatype: " + c.Type)
}

func (table *TTable) FqName() string {
	return "[" + table.Schema + "].[" + table.Name + "]"
}

func (table *TTable) DumpStatement() (string, error) {
	columnDumpStatements := make([]string, 0)
	
	for _, c := range table.Columns {
		columnDumpStatement, err := c.DumpStatement()
		
		if err != nil {
			return "", err
		} else if len(columnDumpStatement) > 0 {
			columnDumpStatements = append(columnDumpStatements, columnDumpStatement)
		}
	}

	return "select (select " + strings.Join(columnDumpStatements, ",") + " for json path, without_array_wrapper) from " + table.FqName(), nil
}

func (table *TTable) InsertStatement() (string, string, error) {
	columnLoadStatements := make([]string, 0)
	usedColumnNames := make([]string, 0)
	for _, c := range table.Columns {
		columnLoadStatement, err := c.LoadStatement()
		
		if err != nil {
			return "", "", err
		} else if len(columnLoadStatement) > 0 {
			usedColumnNames = append(usedColumnNames, c.FqName()) 
			columnLoadStatements = append(columnLoadStatements, columnLoadStatement)
		}
	}

	columnWithStatements := make([]string, len(table.Columns))
	for i, c := range table.Columns {
		columnWithStatements[i] = c.FqName() + " " + c.WithType()
	}

	prefix := "insert " + table.FqName() + " (" + strings.Join(usedColumnNames, ",") + ") " + "\nselect " + strings.Join(columnLoadStatements, ",") + "\nfrom OPENJSON("
	postifx := ")\nwith (" + strings.Join(columnWithStatements, ",") + ")\nGO\n\n"

	return prefix, postifx, nil
}

func (table *TTable) Dump() (error) {
	query, err := table.DumpStatement()
	if err != nil {
		return err
	}

	stmt, err := DB.Prepare(query)
	mfa.CatchFatal(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()


	prefix, postfix, err := table.InsertStatement()
	if err != nil {
		return err
	}

	i := 1
	batchSize := 10
	for ; rows.Next(); i++ {
			var rowValue string

			rows.Scan(&rowValue)
			rowValue = strings.ReplaceAll(rowValue, "'", "''")
			
			if i == 1 {
				fmt.Print(prefix + "N'[" + rowValue)
			} else if i == batchSize {
				fmt.Print("," + rowValue + "]' " + postfix)
				i = 0
			} else {
				fmt.Print("," + rowValue)
			}
	}

	if i < (batchSize +1) && i > 1 {
		fmt.Print("]' " + postfix)
	}
	
 	return nil
}

func (table *TTable) LoadColumnsFromDb() {
	query := `select 
		COLUMN_NAME, 
		DATA_TYPE 
		from INFORMATION_SCHEMA.COLUMNS 
		where TABLE_SCHEMA = '` + table.Schema + `'
		  and TABLE_NAME = '` + table.Name + `' 
		order by ORDINAL_POSITION`;

	// TODO [mfa] use parametrized statement
	stmt, err := DB.Prepare(query)

	mfa.CatchFatal(err)
	defer stmt.Close()
	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()

	for rows.Next() {
		var column TColumn
		rows.Scan(&column.Name, &column.Type)
		table.Columns = append(table.Columns, column)
	}
}


func DumpDataJson() (error) {
	// CREDITS: https://www.mssqltips.com/sqlservertip/6179/sql-server-foreign-key-hierarchy-order-and-dependency-list-script/
	stmt, err := DB.Prepare(`
		WITH dependencies -- Get object with FK dependencies
		AS (
		    SELECT FK.TABLE_NAME AS Obj
		        , PK.TABLE_NAME AS Depends
		    FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS C
		    INNER JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS FK
		        ON C.CONSTRAINT_NAME = FK.CONSTRAINT_NAME
		    INNER JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS PK
		        ON C.UNIQUE_CONSTRAINT_NAME = PK.CONSTRAINT_NAME
		    ), 
		no_dependencies -- The first [level] are objects with no dependencies 
		AS (
		    SELECT 
		        name AS Obj
		    FROM sys.objects
		    WHERE name NOT IN (SELECT obj FROM dependencies) --we remove objects with dependencies from first CTE
		    AND type = 'U' -- Just tables
		    ), 
		recursiv -- recursive CTE to get dependencies
		AS (
		    SELECT Obj AS [name]
		        , CAST('' AS VARCHAR(max)) AS [dependencies]
		        , 0 AS [level] -- [Level] 0 indicate tables with no dependencies
		    FROM no_dependencies
		 
		    UNION ALL
		 
		    SELECT d.Obj AS [name]
		        , CAST(IIF([level] > 0, r.[dependencies] + ' > ', '') + d.Depends AS VARCHAR(max)) -- visually reflects hierarchy
		        , R.[level] + 1 AS [level]
		    FROM dependencies d
		    INNER JOIN recursiv r
		        ON d.Depends = r.[name]
		    )
		-- The final result, with some extra fields for more information
		SELECT DISTINCT SCHEMA_NAME(O.schema_id) AS [
					schema]
		    , R.[name]
		   	, R.[dependencies]
		    , R.[level]
		FROM recursiv R
		INNER JOIN sys.objects O
		    ON R.[name] = O.name
		ORDER BY [level], [name]`)

	mfa.CatchFatal(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()

	var tables []TTable
	for rows.Next() {
			var table TTable
			rows.Scan(&table.Schema, &table.Name, &table.Dependencies, &table.Level)
			tables = append(tables, table)
	}
	mfa.CatchFatal(rows.Err())
	
	for _, table := range tables {
		table.LoadColumnsFromDb()
		err := table.Dump()
		if err != nil {
			log.Println("[failure] " + table.FqName())
			return err
		}
		log.Println("[success] " + table.FqName())
	}

	return nil
}
