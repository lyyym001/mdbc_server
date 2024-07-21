package api

import (
	"database/sql"
	"fmt"
	"mdbc_server/lframework/utils"
)

func RegisterWorkRecord() int {
	var maxUniqueId sql.NullInt32

	db := utils.GlobalObject.SqliteInst.GetDB()
	rows, _ := db.Query("select max(uniqueid) from tb_work; ")
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&maxUniqueId); err == nil {
			//allStudent.Students = append(allStudent.Students, studentInfo)
		} else {
			fmt.Println("RegisterWorkRecord,", err)
		}
	}
	fmt.Println("maxUniqueId = ", maxUniqueId)
	return 1
}
