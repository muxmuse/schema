package schema

import (
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/muxmuse/schema/mfa"
)

func Diff() {
	stmt, err := DB.Prepare(`select
	    [name] = '[' + l.schema_name + '].[' + o.name + ']',
			o.modify_date
		from sys.objects o
		join
		(
		    select
		        [schema_name] = s.name,
		        s.schema_id,
		        -- o.name,
		        o.object_id,
		        o.modify_date
		    from sys.objects o
		    join sys.schemas s on o.schema_id = s.schema_id
		    where o.[name] = 'SCHEMA_INFO'
		) l
		on l.schema_id = o.schema_id
		where o.modify_date > l.modify_date`);

	mfa.CatchFatal(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	mfa.CatchFatal(err)
	defer rows.Close()

	for rows.Next() {
			var name string
			var modifiedAt string

			rows.Scan(&name, &modifiedAt)
			fmt.Println(name, modifiedAt)
	}
	mfa.CatchFatal(rows.Err())
}