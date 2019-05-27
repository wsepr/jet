package tests

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/sub0zero/go-sqlbuilder/sqlbuilder"
	"github.com/sub0zero/go-sqlbuilder/tests/.test_files/dvd_rental/dvds/model"
	. "github.com/sub0zero/go-sqlbuilder/tests/.test_files/dvd_rental/dvds/table"
	"gotest.tools/assert"
	"testing"
)

var query = Inventory.
	SELECT(Inventory.AllColumns).
	LIMIT(1).
	ORDER_BY(Inventory.InventoryID)

func TestScanToInvalidDestination(t *testing.T) {

	t.Run("nil dest", func(t *testing.T) {
		err := query.Query(db, nil)

		assert.Error(t, err, "Destination is nil. ")
	})

	t.Run("struct dest", func(t *testing.T) {
		err := query.Query(db, struct{}{})

		assert.Error(t, err, "Destination has to be a pointer to slice or pointer to struct. ")
	})

	t.Run("slice dest", func(t *testing.T) {
		err := query.Query(db, []struct{}{})

		assert.Error(t, err, "Destination has to be a pointer to slice or pointer to struct. ")
	})

	t.Run("slice of pointers to pointer dest", func(t *testing.T) {
		err := query.Query(db, []**struct{}{})

		assert.Error(t, err, "Destination has to be a pointer to slice or pointer to struct. ")
	})

	t.Run("map dest", func(t *testing.T) {
		err := query.Query(db, []map[string]string{})

		assert.Error(t, err, "Destination has to be a pointer to slice or pointer to struct. ")
	})
}

func TestScanToValidDestination(t *testing.T) {
	t.Run("pointer to struct", func(t *testing.T) {
		err := query.Query(db, &struct{}{})

		assert.NilError(t, err)
	})

	t.Run("pointer to slice", func(t *testing.T) {
		err := query.Query(db, &[]struct{}{})

		assert.NilError(t, err)
	})

	t.Run("pointer to slice of pointer to structs", func(t *testing.T) {
		err := query.Query(db, &[]*struct{}{})

		assert.NilError(t, err)
	})

	t.Run("pointer to slice of strings", func(t *testing.T) {
		err := query.Query(db, &[]int32{})

		assert.NilError(t, err)
	})

	t.Run("pointer to slice of strings", func(t *testing.T) {
		err := query.Query(db, &[]*int32{})

		assert.NilError(t, err)
	})
}

func TestScanToStruct(t *testing.T) {
	query := Inventory.
		SELECT(Inventory.AllColumns).
		ORDER_BY(Inventory.InventoryID)

	t.Run("one struct", func(t *testing.T) {
		dest := model.Inventory{}
		err := query.LIMIT(1).Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, inventory1, dest)
	})

	t.Run("multiple structs, just first one used", func(t *testing.T) {
		dest := model.Inventory{}
		err := query.LIMIT(10).Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, inventory1, dest)
	})

	t.Run("one struct", func(t *testing.T) {
		dest := struct {
			model.Inventory
		}{}
		err := query.LIMIT(1).Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, inventory1, dest.Inventory)
	})

	t.Run("one struct", func(t *testing.T) {
		dest := struct {
			*model.Inventory
		}{}
		err := query.LIMIT(1).Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, inventory1, *dest.Inventory)
	})

	t.Run("invalid dest", func(t *testing.T) {
		dest := struct {
			Inventory **model.Inventory
		}{}

		err := query.Query(db, &dest)

		assert.Error(t, err, "Unsupported dest type: Inventory **model.Inventory")
	})

	t.Run("invalid dest 2", func(t *testing.T) {
		dest := struct {
			Inventory ***model.Inventory
		}{}

		err := query.Query(db, &dest)

		assert.Error(t, err, "Unsupported dest type: Inventory ***model.Inventory")
	})

	t.Run("custom struct", func(t *testing.T) {
		type Inventory struct {
			InventoryID *int32 `sql:"unique"`
			FilmID      int16
			StoreID     *int16
		}

		dest := Inventory{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)

		assert.Equal(t, *dest.InventoryID, int32(1))
		assert.Equal(t, dest.FilmID, int16(1))
		assert.Equal(t, *dest.StoreID, int16(1))
	})

	t.Run("type mismatch base type", func(t *testing.T) {
		type Inventory struct {
			InventoryID int
			FilmID      string
		}

		dest := Inventory{}

		err := query.Query(db, &dest)

		assert.Error(t, err, `Scan: can't set int32 to int, at struct field: InventoryID int of type tests.Inventory. `)

		fmt.Println(err)
	})

	t.Run("type mismatch scanner type", func(t *testing.T) {
		type Inventory struct {
			InventoryID uuid.UUID
			FilmID      string
		}

		dest := Inventory{}

		err := query.Query(db, &dest)

		assert.Error(t, err, `Scan: unable to scan type int32 into UUID, at struct field: InventoryID uuid.UUID of type tests.Inventory. `)
	})

}

func TestScanToNestedStruct(t *testing.T) {
	query := Inventory.
		INNER_JOIN(Film, Inventory.FilmID.Eq(Film.FilmID)).
		INNER_JOIN(Store, Inventory.StoreID.Eq(Store.StoreID)).
		SELECT(Inventory.AllColumns, Film.AllColumns, Store.AllColumns).
		WHERE(Inventory.InventoryID.EqL(1))

	t.Run("embedded structs", func(t *testing.T) {
		dest := struct {
			model.Inventory
			model.Film
			model.Store
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.DeepEqual(t, dest.Film, film1)
		assert.DeepEqual(t, dest.Store, store1)
	})

	t.Run("embedded pointer structs", func(t *testing.T) {
		dest := struct {
			*model.Inventory
			*model.Film
			*model.Store
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, *dest.Inventory, inventory1)
		assert.DeepEqual(t, *dest.Film, film1)
		assert.DeepEqual(t, *dest.Store, store1)
	})

	t.Run("embedded unused structs", func(t *testing.T) {
		dest := struct {
			model.Inventory
			model.Actor //unused
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.DeepEqual(t, dest.Actor, model.Actor{})
	})

	t.Run("embedded unused pointer structs", func(t *testing.T) {
		dest := struct {
			model.Inventory
			*model.Actor //unused
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.DeepEqual(t, dest.Actor, (*model.Actor)(nil))
	})

	t.Run("embedded unused pointer structs", func(t *testing.T) {
		dest := struct {
			model.Inventory
			Actor *model.Actor //unused
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.DeepEqual(t, dest.Actor, (*model.Actor)(nil))
	})

	t.Run("embedded pointer to selected column", func(t *testing.T) {
		query := Inventory.
			INNER_JOIN(Film, Inventory.FilmID.Eq(Film.FilmID)).
			INNER_JOIN(Store, Inventory.StoreID.Eq(Store.StoreID)).
			SELECT(Inventory.AllColumns, Film.AllColumns, Store.AllColumns, Literal("").AS("actor.first_name")).
			WHERE(Inventory.InventoryID.EqL(1))

		dest := struct {
			model.Inventory
			Actor *model.Actor //unused
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.Assert(t, dest.Actor != nil)
	})

	t.Run("struct embedded unused pointer", func(t *testing.T) {
		dest := struct {
			model.Inventory
			Actor *struct {
				model.Actor
			} //unused
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.DeepEqual(t, dest.Actor, (*struct{ model.Actor })(nil))
	})

	t.Run("multiple embedded unused pointer", func(t *testing.T) {
		dest := struct {
			model.Inventory
			Actor *struct {
				model.Actor    //unused
				model.Language //unesed
			}
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.DeepEqual(t, dest.Actor, (*struct {
			model.Actor
			model.Language
		})(nil))
	})

	t.Run("field not nil, embedded selected model", func(t *testing.T) {
		dest := struct {
			model.Inventory
			Actor *struct {
				model.Actor //unselected
				model.Film  //selected
			}
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.Assert(t, dest.Actor != nil)
		assert.DeepEqual(t, dest.Actor.Actor, model.Actor{})
		assert.DeepEqual(t, dest.Actor.Film, film1)
	})

	t.Run("field not nil, deeply nested selected model", func(t *testing.T) {
		dest := struct {
			model.Inventory
			Actor *struct {
				model.Actor //unselected
				Film        *struct {
					*model.Film //selected
				}
			}
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.Assert(t, dest.Actor != nil)
		assert.Assert(t, dest.Actor.Film != nil)
		assert.DeepEqual(t, dest.Actor.Film.Film, &film1)
	})

	t.Run("embedded structs", func(t *testing.T) {
		query := Inventory.
			INNER_JOIN(Film, Inventory.FilmID.Eq(Film.FilmID)).
			INNER_JOIN(Store, Inventory.StoreID.Eq(Store.StoreID)).
			INNER_JOIN(Language, Film.LanguageID.Eq(Language.LanguageID)).
			SELECT(Inventory.AllColumns, Film.AllColumns, Store.AllColumns, Language.AllColumns).
			WHERE(Inventory.InventoryID.EqL(1))

		dest := struct {
			model.Inventory
			Film struct {
				model.Film

				Language  model.Language
				Language2 *model.Language
				Language3 *model.Language `sqlbuilder:"language"`
				Lang      struct {
					model.Language
				}
				Lang2 *struct {
					model.Language
				}
			}
			Store model.Store
		}{}

		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.DeepEqual(t, dest.Inventory, inventory1)
		assert.DeepEqual(t, dest.Film.Film, film1)
		assert.DeepEqual(t, dest.Store, store1)
		assert.DeepEqual(t, dest.Film.Language, language1)
		assert.DeepEqual(t, dest.Film.Lang.Language, language1)
		assert.DeepEqual(t, dest.Film.Lang2.Language, language1)
		assert.DeepEqual(t, dest.Film.Language2, (*model.Language)(nil))
		assert.DeepEqual(t, dest.Film.Language3, &language1)
	})
}

func TestScanToSlice(t *testing.T) {

	t.Run("slice of structs", func(t *testing.T) {
		query := Inventory.
			SELECT(Inventory.AllColumns).
			ORDER_BY(Inventory.InventoryID).
			LIMIT(10)

		t.Run("slice od inventory", func(t *testing.T) {
			dest := []model.Inventory{}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 10)
			assert.DeepEqual(t, dest[0], inventory1)
			assert.DeepEqual(t, dest[1], inventory2)
		})

		t.Run("slice of ints", func(t *testing.T) {
			var dest []int32

			err := query.Query(db, &dest)
			assert.NilError(t, err)
			assert.DeepEqual(t, dest, []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

		})

		t.Run("slice type mismatch ", func(t *testing.T) {
			var dest []int

			err := query.Query(db, &dest)
			assert.Error(t, err, `Scan: can't append int32 to []int slice `)
		})
	})

	t.Run("slice of complex structs", func(t *testing.T) {
		query := Inventory.
			INNER_JOIN(Film, Inventory.FilmID.Eq(Film.FilmID)).
			INNER_JOIN(Store, Inventory.StoreID.Eq(Store.StoreID)).
			SELECT(Inventory.AllColumns, Film.AllColumns, Store.AllColumns).
			ORDER_BY(Inventory.InventoryID).
			LIMIT(10)

		t.Run("struct with slice of ints", func(t *testing.T) {
			var dest struct {
				model.Film
				IDs []int32 `sqlbuilder:"inventory.inventory_id"`
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.DeepEqual(t, dest.Film, film1)
			assert.DeepEqual(t, dest.IDs, []int32{1, 2, 3, 4, 5, 6, 7, 8})
		})

		t.Run("slice of structs with slice of ints", func(t *testing.T) {
			var dest []struct {
				model.Film
				IDs []int32 `sqlbuilder:"inventory.inventory_id"`
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 2)
			assert.DeepEqual(t, dest[0].Film, film1)
			assert.DeepEqual(t, dest[0].IDs, []int32{1, 2, 3, 4, 5, 6, 7, 8})
			assert.DeepEqual(t, dest[1].Film, film2)
			assert.DeepEqual(t, dest[1].IDs, []int32{9, 10})
		})

		t.Run("slice of structs with slice of pointer to ints", func(t *testing.T) {
			var dest []struct {
				model.Film
				IDs []*int32 `sqlbuilder:"inventory.inventory_id"`
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 2)
			assert.DeepEqual(t, dest[0].Film, film1)
			assert.DeepEqual(t, dest[0].IDs, []*int32{int32Ptr(1), int32Ptr(2), int32Ptr(3), int32Ptr(4),
				int32Ptr(5), int32Ptr(6), int32Ptr(7), int32Ptr(8)})
			assert.DeepEqual(t, dest[1].Film, film2)
			assert.DeepEqual(t, dest[1].IDs, []*int32{int32Ptr(9), int32Ptr(10)})
		})

		t.Run("complex struct 1", func(t *testing.T) {
			dest := []struct {
				model.Inventory
				model.Film
				model.Store
			}{}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 10)
			assert.DeepEqual(t, dest[0].Inventory, inventory1)
			assert.DeepEqual(t, dest[0].Film, film1)
			assert.DeepEqual(t, dest[0].Store, store1)

			assert.DeepEqual(t, dest[1].Inventory, inventory2)
		})

		t.Run("complex struct 2", func(t *testing.T) {
			var dest []struct {
				*model.Inventory
				model.Film
				*model.Store
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 10)
			assert.DeepEqual(t, dest[0].Inventory, &inventory1)
			assert.DeepEqual(t, dest[0].Film, film1)
			assert.DeepEqual(t, dest[0].Store, &store1)

			assert.DeepEqual(t, dest[1].Inventory, &inventory2)
		})

		t.Run("complex struct 3", func(t *testing.T) {
			var dest []struct {
				Inventory model.Inventory
				Film      *model.Film
				Store     struct {
					*model.Store
				}
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 10)
			assert.DeepEqual(t, dest[0].Inventory, inventory1)
			assert.DeepEqual(t, dest[0].Film, &film1)
			assert.DeepEqual(t, dest[0].Store.Store, &store1)

			assert.DeepEqual(t, dest[1].Inventory, inventory2)
		})

		t.Run("complex struct 4", func(t *testing.T) {
			var dest []struct {
				model.Film

				Inventories []struct {
					model.Inventory
					model.Store
				}
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 2)
			assert.DeepEqual(t, dest[0].Film, film1)
			assert.DeepEqual(t, len(dest[0].Inventories), 8)
			assert.DeepEqual(t, dest[0].Inventories[0].Inventory, inventory1)
			assert.DeepEqual(t, dest[0].Inventories[0].Store, store1)
		})

		t.Run("complex struct 5", func(t *testing.T) {
			var dest []struct {
				model.Film

				Inventories []struct {
					model.Inventory

					Rentals  *[]model.Rental
					Rentals2 []model.Rental
				}
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)

			assert.Equal(t, len(dest), 2)
			assert.DeepEqual(t, dest[0].Film, film1)
			assert.Equal(t, len(dest[0].Inventories), 8)
			assert.DeepEqual(t, dest[0].Inventories[0].Inventory, inventory1)
			assert.Assert(t, dest[0].Inventories[0].Rentals == nil)
			assert.Assert(t, dest[0].Inventories[0].Rentals2 == nil)
		})
	})

	t.Run("slice of complex structs 2", func(t *testing.T) {
		query := Country.
			INNER_JOIN(City, City.CountryID.Eq(Country.CountryID)).
			INNER_JOIN(Address, Address.CityID.Eq(City.CityID)).
			INNER_JOIN(Customer, Customer.AddressID.Eq(Address.AddressID)).
			SELECT(Country.AllColumns, City.AllColumns, Address.AllColumns, Customer.AllColumns).
			ORDER_BY(Country.CountryID.ASC(), City.CityID.ASC(), Address.AddressID.ASC(), Customer.CustomerID.ASC()).
			LIMIT(1000)

		t.Run("dest1", func(t *testing.T) {
			var dest []struct {
				model.Country

				Cities []struct {
					model.City

					Adresses []struct {
						model.Address

						Customer model.Customer
					}
				}
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 108)
			assert.DeepEqual(t, dest[100].Country, countryUk)
			assert.Equal(t, len(dest[100].Cities), 8)
			assert.DeepEqual(t, dest[100].Cities[2].City, cityLondon)
			assert.Equal(t, len(dest[100].Cities[2].Adresses), 2)
			assert.DeepEqual(t, dest[100].Cities[2].Adresses[0].Address, address256)
			assert.DeepEqual(t, dest[100].Cities[2].Adresses[0].Customer, customer256)
			assert.DeepEqual(t, dest[100].Cities[2].Adresses[1].Address, addres517)
			assert.DeepEqual(t, dest[100].Cities[2].Adresses[1].Customer, customer512)
		})

		t.Run("dest1", func(t *testing.T) {
			var dest []*struct {
				*model.Country

				Cities []*struct {
					*model.City

					Adresses *[]*struct {
						*model.Address

						Customer *model.Customer
					}
				}
			}

			err := query.Query(db, &dest)

			assert.NilError(t, err)
			assert.Equal(t, len(dest), 108)
			assert.DeepEqual(t, dest[100].Country, &countryUk)
			assert.Equal(t, len(dest[100].Cities), 8)
			assert.DeepEqual(t, dest[100].Cities[2].City, &cityLondon)
			assert.Equal(t, len(*dest[100].Cities[2].Adresses), 2)
			assert.DeepEqual(t, (*dest[100].Cities[2].Adresses)[0].Address, &address256)
			assert.DeepEqual(t, (*dest[100].Cities[2].Adresses)[0].Customer, &customer256)
			assert.DeepEqual(t, (*dest[100].Cities[2].Adresses)[1].Address, &addres517)
			assert.DeepEqual(t, (*dest[100].Cities[2].Adresses)[1].Customer, &customer512)
		})

	})

	t.Run("dest1", func(t *testing.T) {
		var dest []*struct {
			*model.Country

			Cities []**struct {
				*model.City
			}
		}

		err := query.Query(db, &dest)

		assert.Error(t, err, "Unsupported dest type: Cities []**struct { *model.City }")
	})
}

var address256 = model.Address{
	AddressID:  256,
	Address:    "1497 Yuzhou Drive",
	Address2:   stringPtr(""),
	District:   "England",
	CityID:     312,
	PostalCode: stringPtr("3433"),
	Phone:      "246810237916",
	LastUpdate: *timestampWithoutTimeZone("2006-02-15 09:45:30", 0),
}

var addres517 = model.Address{
	AddressID:  517,
	Address:    "548 Uruapan Street",
	Address2:   stringPtr(""),
	District:   "Ontario",
	CityID:     312,
	PostalCode: stringPtr("35653"),
	Phone:      "879347453467",
	LastUpdate: *timestampWithoutTimeZone("2006-02-15 09:45:30", 0),
}

var customer256 = model.Customer{
	CustomerID: 252,
	StoreID:    2,
	FirstName:  "Mattie",
	LastName:   "Hoffman",
	Email:      stringPtr("mattie.hoffman@sakilacustomer.org"),
	AddressID:  256,
	Activebool: true,
	CreateDate: *timestampWithoutTimeZone("2006-02-14 00:00:00", 0),
	LastUpdate: timestampWithoutTimeZone("2013-05-26 14:49:45.738", 0),
	Active:     int32Ptr(1),
}

var customer512 = model.Customer{
	CustomerID: 512,
	StoreID:    1,
	FirstName:  "Cecil",
	LastName:   "Vines",
	Email:      stringPtr("cecil.vines@sakilacustomer.org"),
	AddressID:  517,
	Activebool: true,
	CreateDate: *timestampWithoutTimeZone("2006-02-14 00:00:00", 0),
	LastUpdate: timestampWithoutTimeZone("2013-05-26 14:49:45.738", 0),
	Active:     int32Ptr(1),
}

var countryUk = model.Country{
	CountryID:  102,
	Country:    "United Kingdom",
	LastUpdate: *timestampWithoutTimeZone("2006-02-15 09:44:00", 0),
}

var cityLondon = model.City{
	CityID:     312,
	City:       "London",
	CountryID:  102,
	LastUpdate: *timestampWithoutTimeZone("2006-02-15 09:45:25", 0),
}

var inventory1 = model.Inventory{
	InventoryID: 1,
	FilmID:      1,
	StoreID:     1,
	LastUpdate:  *timestampWithoutTimeZone("2006-02-15 10:09:17", 0),
}

var inventory2 = model.Inventory{
	InventoryID: 2,
	FilmID:      1,
	StoreID:     1,
	LastUpdate:  *timestampWithoutTimeZone("2006-02-15 10:09:17", 0),
}

var film1 = model.Film{
	FilmID:          1,
	Title:           "Academy Dinosaur",
	Description:     stringPtr("A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies"),
	ReleaseYear:     int32Ptr(2006),
	LanguageID:      1,
	RentalDuration:  6,
	RentalRate:      0.99,
	Length:          int16Ptr(86),
	ReplacementCost: 20.99,
	Rating:          &pgRating,
	LastUpdate:      *timestampWithoutTimeZone("2013-05-26 14:50:58.951", 3),
	SpecialFeatures: stringPtr("{\"Deleted Scenes\",\"Behind the Scenes\"}"),
	Fulltext:        "'academi':1 'battl':15 'canadian':20 'dinosaur':2 'drama':5 'epic':4 'feminist':8 'mad':11 'must':14 'rocki':21 'scientist':12 'teacher':17",
}

var film2 = model.Film{
	FilmID:          2,
	Title:           "Ace Goldfinger",
	Description:     stringPtr("A Astounding Epistle of a Database Administrator And a Explorer who must Find a Car in Ancient China"),
	ReleaseYear:     int32Ptr(2006),
	LanguageID:      1,
	RentalDuration:  3,
	RentalRate:      4.99,
	Length:          int16Ptr(48),
	ReplacementCost: 12.99,
	Rating:          &gRating,
	LastUpdate:      *timestampWithoutTimeZone("2013-05-26 14:50:58.951", 3),
	SpecialFeatures: stringPtr(`{Trailers,"Deleted Scenes"}`),
	Fulltext:        `'ace':1 'administr':9 'ancient':19 'astound':4 'car':17 'china':20 'databas':8 'epistl':5 'explor':12 'find':15 'goldfing':2 'must':14`,
}

var store1 = model.Store{
	StoreID:        1,
	ManagerStaffID: 1,
	AddressID:      1,
	LastUpdate:     *timestampWithoutTimeZone("2006-02-15 09:57:12", 0),
}

var pgRating = model.MpaaRating_PG
var gRating = model.MpaaRating_G

var language1 = model.Language{
	LanguageID: 1,
	Name:       "English             ",
	LastUpdate: *timestampWithoutTimeZone("2006-02-15 10:02:19", 0),
}
