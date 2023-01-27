package schema

import (
	"io"
	"os"
	"fmt"
	"bufio"
)

func execBatchesFromStdin() (error, string, SQLError) {
	reader := bufio.NewReader(os.Stdin)
  
  var batch string
  
  for {
  	line, isPrefix, err := reader.ReadLine()

  	if err != nil {
  		if err == io.EOF {
  			return nil, "", nil	
  		} else {
  			return err, "", nil
  		}
  	}
  	
  	if isPrefix {
  		batch += string(line)
  		continue
  	} else {
  		batch += string(line) + "\n"
  	}

  	if string(line) == "GO" {
  		_, err = DB.Exec(batch)
			if err != nil {
				if sqlError, ok := err.(SQLError); ok {
					return err, batch, sqlError
				} else {
					return err, batch, nil
				}
			}
			batch = ""
  	}
  }

  return nil, "", nil
}

func ExecFromStdin() error {
	scriptErr, batch, sqlError := execBatchesFromStdin()

	if scriptErr != nil {
		fmt.Println("[failure]", "SQL from stdin executed with errors")
		if sqlError != nil {
			fmt.Println("-------------------------------------------------------------------------------")
			printLines(batch, sqlError.SQLErrorLineNo() - 20, sqlError.SQLErrorLineNo() + 20, sqlError.SQLErrorLineNo())
			fmt.Println("-------------------------------------------------------------------------------")
			fmt.Println("[failure]", "Line", sqlError.SQLErrorLineNo(), "in", sqlError.SQLErrorProcName())
			fmt.Println("[failure]", "E" + fmt.Sprint(sqlError.SQLErrorNumber()), sqlError.SQLErrorMessage())
		}
		return scriptErr
	}

	return nil
}