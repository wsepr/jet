package mysql

import (
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-jet/jet/v2/generator/mysql"
	"github.com/go-jet/jet/v2/internal/testutils"
	"github.com/go-jet/jet/v2/tests/dbconfig"
)

const genTestDirRoot = "./.gentestdata3"
const genTestDir3 = "./.gentestdata3/mysql"

func TestGenerator(t *testing.T) {

	for i := 0; i < 3; i++ {
		err := mysql.Generate(genTestDir3, dbConnection("dvds"))

		require.NoError(t, err)

		assertGeneratedFiles(t)
	}

	for i := 0; i < 3; i++ {
		dsn := dbconfig.MySQLConnectionString(sourceIsMariaDB(), "dvds")

		err := mysql.GenerateDSN(dsn, genTestDir3)

		require.NoError(t, err)
		assertGeneratedFiles(t)
	}

	err := os.RemoveAll(genTestDirRoot)
	require.NoError(t, err)
}

func TestCmdGenerator(t *testing.T) {
	err := os.RemoveAll(genTestDir3)
	require.NoError(t, err)

	var cmd *exec.Cmd

	if sourceIsMariaDB() {
		cmd = exec.Command("jet",
			"-source=MariaDB",
			"-dbname=dvds",
			"-host="+dbconfig.MariaDBHost,
			"-port="+strconv.Itoa(dbconfig.MariaDBPort),
			"-user="+dbconfig.MariaDBUser,
			"-password="+dbconfig.MariaDBPassword,
			"-path="+genTestDir3)
	} else {
		cmd = exec.Command("jet",
			"-source=MySQL",
			"-dbname=dvds",
			"-host="+dbconfig.MySqLHost,
			"-port="+strconv.Itoa(dbconfig.MySQLPort),
			"-user="+dbconfig.MySQLUser,
			"-password="+dbconfig.MySQLPassword,
			"-path="+genTestDir3)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	require.NoError(t, err)

	assertGeneratedFiles(t)

	err = os.RemoveAll(genTestDirRoot)
	require.NoError(t, err)

	// check that generation via DSN works
	dsn := "mysql://" + dbconfig.MySQLConnectionString(sourceIsMariaDB(), "dvds")
	cmd = exec.Command("jet", "-dsn="+dsn, "-path="+genTestDir3)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	require.NoError(t, err)
}

func TestIgnoreTablesViewsEnums(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "with dsn",
			args: []string{
				"-dsn=mysql://" + dbconfig.MySQLConnectionString(sourceIsMariaDB(), "dvds"),
				"-ignore-tables=actor,ADDRESS,Category, city ,country,staff,store,rental",
				"-ignore-views=actor_info,CUSTomER_LIST, film_list",
				"-ignore-enums=film_list_rating,film_rating",
				"-path=" + genTestDir3,
			},
		},
		{
			name: "without dsn",
			args: []string{
				"-source=MySQL",
				"-dbname=dvds",
				"-host=" + dbconfig.MySqLHost,
				"-port=" + strconv.Itoa(dbconfig.MySQLPort),
				"-user=" + dbconfig.MySQLUser,
				"-password=" + dbconfig.MySQLPassword,
				"-ignore-tables=actor,ADDRESS,Category, city ,country,staff,store,rental",
				"-ignore-views=actor_info,CUSTomER_LIST, film_list",
				"-ignore-enums=film_list_rating,film_rating",
				"-path=" + genTestDir3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("jet", tt.args...)

			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			err := cmd.Run()
			require.NoError(t, err)

			testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/table", "customer.go", "film.go", "film_actor.go",
				"film_category.go", "film_text.go", "inventory.go", "language.go", "payment.go", "table_use_schema.go")

			testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/view", "nicer_but_slower_film_list.go",
				"sales_by_film_category.go", "sales_by_store.go", "staff_list.go", "view_use_schema.go")

			testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/enum", "nicer_but_slower_film_list_rating.go")

			testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/model",
				"customer.go", "film.go", "film_actor.go", "film_category.go", "film_text.go", "inventory.go", "language.go",
				"payment.go", "nicer_but_slower_film_list_rating.go", "nicer_but_slower_film_list.go", "sales_by_film_category.go",
				"sales_by_store.go", "staff_list.go")
		})
	}
}

func assertGeneratedFiles(t *testing.T) {
	// Table SQL Builder files
	testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/table", "actor.go", "address.go", "category.go", "city.go", "country.go",
		"customer.go", "film.go", "film_actor.go", "film_category.go", "film_text.go", "inventory.go", "language.go",
		"payment.go", "rental.go", "staff.go", "store.go", "table_use_schema.go")

	testutils.AssertFileContent(t, genTestDir3+"/dvds/table/actor.go", actorSQLBuilderFile)
	testutils.AssertFileContent(t, genTestDir3+"/dvds/table/table_use_schema.go", tableUseSchemaFile)

	// View SQL Builder files
	testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/view", "actor_info.go", "film_list.go", "nicer_but_slower_film_list.go",
		"sales_by_film_category.go", "customer_list.go", "sales_by_store.go", "staff_list.go", "view_use_schema.go")

	testutils.AssertFileContent(t, genTestDir3+"/dvds/view/actor_info.go", actorInfoSQLBuilderFile)
	testutils.AssertFileContent(t, genTestDir3+"/dvds/view/view_use_schema.go", viewUseSchemaFile)

	// Enums SQL Builder files
	testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/enum", "film_rating.go", "film_list_rating.go", "nicer_but_slower_film_list_rating.go")
	testutils.AssertFileContent(t, genTestDir3+"/dvds/enum/film_rating.go", mpaaRatingEnumFile)

	// Model files
	testutils.AssertFileNamesEqual(t, genTestDir3+"/dvds/model", "actor.go", "address.go", "category.go", "city.go", "country.go",
		"customer.go", "film.go", "film_actor.go", "film_category.go", "film_text.go", "inventory.go", "language.go",
		"payment.go", "rental.go", "staff.go", "store.go",
		"film_rating.go", "film_list_rating.go", "nicer_but_slower_film_list_rating.go",
		"actor_info.go", "film_list.go", "nicer_but_slower_film_list.go", "sales_by_film_category.go",
		"customer_list.go", "sales_by_store.go", "staff_list.go")

	testutils.AssertFileContent(t, genTestDir3+"/dvds/model/actor.go", actorModelFile)
}

func TestModelColumnComment(t *testing.T) {
	testutils.AssertFileContent(t, "./../.gentestdata/mysql/test_sample/model/link.go", `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type Link struct {
	ID          int32   `+"`sql:\"primary_key\"`"+` // this is link id
	URL         string  // link url
	Name        string  // link name
	Description *string // this is link description
}
`)
}

var mpaaRatingEnumFile = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/mysql"

var FilmRating = &struct {
	G    mysql.StringExpression
	Pg   mysql.StringExpression
	Pg13 mysql.StringExpression
	R    mysql.StringExpression
	Nc17 mysql.StringExpression
}{
	G:    mysql.NewEnumValue("G"),
	Pg:   mysql.NewEnumValue("PG"),
	Pg13: mysql.NewEnumValue("PG-13"),
	R:    mysql.NewEnumValue("R"),
	Nc17: mysql.NewEnumValue("NC-17"),
}
`

var actorSQLBuilderFile = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/mysql"
)

var Actor = newActorTable("dvds", "actor", "")

type actorTable struct {
	mysql.Table

	// Columns
	ActorID    mysql.ColumnInteger
	FirstName  mysql.ColumnString
	LastName   mysql.ColumnString
	LastUpdate mysql.ColumnTimestamp

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type ActorTable struct {
	actorTable

	NEW actorTable
}

// AS creates new ActorTable with assigned alias
func (a ActorTable) AS(alias string) *ActorTable {
	return newActorTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ActorTable with assigned schema name
func (a ActorTable) FromSchema(schemaName string) *ActorTable {
	return newActorTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ActorTable with assigned table prefix
func (a ActorTable) WithPrefix(prefix string) *ActorTable {
	return newActorTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ActorTable with assigned table suffix
func (a ActorTable) WithSuffix(suffix string) *ActorTable {
	return newActorTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newActorTable(schemaName, tableName, alias string) *ActorTable {
	return &ActorTable{
		actorTable: newActorTableImpl(schemaName, tableName, alias),
		NEW:        newActorTableImpl("", "new", ""),
	}
}

func newActorTableImpl(schemaName, tableName, alias string) actorTable {
	var (
		ActorIDColumn    = mysql.IntegerColumn("actor_id")
		FirstNameColumn  = mysql.StringColumn("first_name")
		LastNameColumn   = mysql.StringColumn("last_name")
		LastUpdateColumn = mysql.TimestampColumn("last_update")
		allColumns       = mysql.ColumnList{ActorIDColumn, FirstNameColumn, LastNameColumn, LastUpdateColumn}
		mutableColumns   = mysql.ColumnList{FirstNameColumn, LastNameColumn, LastUpdateColumn}
	)

	return actorTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ActorID:    ActorIDColumn,
		FirstName:  FirstNameColumn,
		LastName:   LastNameColumn,
		LastUpdate: LastUpdateColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
`

var tableUseSchemaFile = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	Actor = Actor.FromSchema(schema)
	Address = Address.FromSchema(schema)
	Category = Category.FromSchema(schema)
	City = City.FromSchema(schema)
	Country = Country.FromSchema(schema)
	Customer = Customer.FromSchema(schema)
	Film = Film.FromSchema(schema)
	FilmActor = FilmActor.FromSchema(schema)
	FilmCategory = FilmCategory.FromSchema(schema)
	FilmText = FilmText.FromSchema(schema)
	Inventory = Inventory.FromSchema(schema)
	Language = Language.FromSchema(schema)
	Payment = Payment.FromSchema(schema)
	Rental = Rental.FromSchema(schema)
	Staff = Staff.FromSchema(schema)
	Store = Store.FromSchema(schema)
}
`

var actorModelFile = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Actor struct {
	ActorID    uint16 ` + "`sql:\"primary_key\"`" + `
	FirstName  string
	LastName   string
	LastUpdate time.Time
}
`

var actorInfoSQLBuilderFile = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package view

import (
	"github.com/go-jet/jet/v2/mysql"
)

var ActorInfo = newActorInfoTable("dvds", "actor_info", "")

type actorInfoTable struct {
	mysql.Table

	// Columns
	ActorID   mysql.ColumnInteger
	FirstName mysql.ColumnString
	LastName  mysql.ColumnString
	FilmInfo  mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type ActorInfoTable struct {
	actorInfoTable

	NEW actorInfoTable
}

// AS creates new ActorInfoTable with assigned alias
func (a ActorInfoTable) AS(alias string) *ActorInfoTable {
	return newActorInfoTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ActorInfoTable with assigned schema name
func (a ActorInfoTable) FromSchema(schemaName string) *ActorInfoTable {
	return newActorInfoTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ActorInfoTable with assigned table prefix
func (a ActorInfoTable) WithPrefix(prefix string) *ActorInfoTable {
	return newActorInfoTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ActorInfoTable with assigned table suffix
func (a ActorInfoTable) WithSuffix(suffix string) *ActorInfoTable {
	return newActorInfoTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newActorInfoTable(schemaName, tableName, alias string) *ActorInfoTable {
	return &ActorInfoTable{
		actorInfoTable: newActorInfoTableImpl(schemaName, tableName, alias),
		NEW:            newActorInfoTableImpl("", "new", ""),
	}
}

func newActorInfoTableImpl(schemaName, tableName, alias string) actorInfoTable {
	var (
		ActorIDColumn   = mysql.IntegerColumn("actor_id")
		FirstNameColumn = mysql.StringColumn("first_name")
		LastNameColumn  = mysql.StringColumn("last_name")
		FilmInfoColumn  = mysql.StringColumn("film_info")
		allColumns      = mysql.ColumnList{ActorIDColumn, FirstNameColumn, LastNameColumn, FilmInfoColumn}
		mutableColumns  = mysql.ColumnList{ActorIDColumn, FirstNameColumn, LastNameColumn, FilmInfoColumn}
	)

	return actorInfoTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ActorID:   ActorIDColumn,
		FirstName: FirstNameColumn,
		LastName:  LastNameColumn,
		FilmInfo:  FilmInfoColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
`
var viewUseSchemaFile = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package view

// UseSchema sets a new schema name for all generated view SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	ActorInfo = ActorInfo.FromSchema(schema)
	CustomerList = CustomerList.FromSchema(schema)
	FilmList = FilmList.FromSchema(schema)
	NicerButSlowerFilmList = NicerButSlowerFilmList.FromSchema(schema)
	SalesByFilmCategory = SalesByFilmCategory.FromSchema(schema)
	SalesByStore = SalesByStore.FromSchema(schema)
	StaffList = StaffList.FromSchema(schema)
}
`
