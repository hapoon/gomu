package gomu

import (
	"database/sql"
	"fmt"
	"testing"

	gorp "gopkg.in/gorp.v1"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabase(t *testing.T) {
	dbmap := initDB()
	defer dbmap.Db.Close()

	err := dbmap.TruncateTables()
	checkError(err)
	/*
		t1 := TestTable{
			Name:      StringFrom("test"),
			Age:       IntFrom(20),
			IsBanned:  BoolFrom(true),
			CreatedAt: TimeFrom(time.Now()),
		}
		err = dbmap.Insert(&t1)
		checkError(err)
	*/

	t2 := TestTable2{
		ID:   3,
		Name: "cccc",
	}
	err = dbmap.Insert(&t2)
	checkError(err)

	var t3 []TestTable2
	_, err = dbmap.Select(&t3, "SELECT * FROM test_table2")
	checkError(err)
	for i, v := range t3 {
		fmt.Printf("primitive string[%d] %v\n", i, v)
	}

	t4 := TestTable3{
		Name: sql.NullString{
			String: "dddd",
			Valid:  true,
		},
	}
	err = dbmap.Insert(&t4)
	checkError(err)

	var t3s []TestTable3
	_, err = dbmap.Select(&t3s, "SELECT * FROM test_table3")
	checkError(err)
	for i, v := range t3s {
		fmt.Printf("database/sql.NullString[%d] %+v\n", i, v)
	}

	t5 := TestTable4{
		Name: StringFrom("eeee"),
	}
	t6 := TestTable4{
		Name: NewString("ffff", true, true),
	}
	t7 := TestTable4{
		Name: NewString("gggg", false, false),
	}
	err = dbmap.Insert(&t5, &t6, &t7)
	checkError(err)

	var t4s []TestTable4
	_, err = dbmap.Select(&t4s, "SELECT * FROM test_table4")
	checkError(err)
	for i, v := range t4s {
		fmt.Printf("gomu.String[%d] %+v\n", i, v)
	}
}

func initDB() *gorp.DbMap {
	db, err := sql.Open("mysql", "valencia:valencia@/valencia_test")
	checkError(err)

	dbmap := &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "UTF8",
		},
	}

	dbmap.AddTableWithName(TestTable{}, "test_table").SetKeys(true, "ID")
	dbmap.AddTableWithName(TestTable2{}, "test_table2").SetKeys(true, "ID")
	dbmap.AddTableWithName(TestTable3{}, "test_table3").SetKeys(true, "ID")
	dbmap.AddTableWithName(TestTable4{}, "test_table4").SetKeys(true, "ID")

	err = dbmap.CreateTablesIfNotExists()
	checkError(err)

	return dbmap
}

type TestTable struct {
	ID        int    `db:"id"`
	Name      String `db:"name"`
	Age       Int    `db:"age"`
	IsBanned  Bool   `db:"is_banned"`
	CreatedAt Time   `db:"created_at"`
}

// primitiveなstringでの挙動確認用
type TestTable2 struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// database/sql.NullStringでの挙動確認用
type TestTable3 struct {
	ID   int            `db:"id"`
	Name sql.NullString `db:"name"`
}

// gomu.Stringの挙動確認用
type TestTable4 struct {
	ID   int    `db:"id"`
	Name String `db:"name"`
}
