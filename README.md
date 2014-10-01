# mogi: SQL mocking library for Go
Sometimes you want to mock an SQL connection. Maybe you want to write tests that don't require a database installed. 

[mogi](http://github.com/guregu/mogi) is a library that makes it easy to write DRY stubs corresponding to specific or general SQL queries. This is thanks to the [vitess](https://github.com/youtube/vitess) SQL parser.

Let's take a look at how to use mogi.

## Some code to test
First, some code that we will write tests for. We define a Beer type and a function called GetBeer that will return a beer or an error given an ID. Assume some setup code connects to the database.

```
var db *sql.DB

type Beer struct {
	ID   int64
	Name string
	Pct  float32
}

func GetBeer(id int64) (beer Beer, err error) {
	query := `SELECT id, name, pct FROM beer WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&beer.ID, &beer.Name, &beer.Pct)
	return
}

```

## mogi setup
We need to make sure our tests use mogi as the database driver. One way to do that is to set our global db variable in an init function in a test file, like `main_test.go`

```
import (
	"database/sql"

	"github.com/guregu/mogi"
)

func init() {
	db, _ = sql.Open("mogi", "")
	// Log unstubbed queries
	mogi.Verbose(true)
}
```

Opening a new database connection in every test is another option.


## Writing the tests
### The fixture
Create a file called `beer_test.go`. First let's define a beer fixture containing some arbitrary data. 

```
package main

var beerFixture = Beer{
	ID:   42,
	Name: "Yona Yona Ale",
	Pct:  5.5,
}
```

### A failing test
Let's write a failing test. We will sprinkle some mogi magic afterwards to make it pass.

```
func TestGetBeer(t *testing.T) {
	beer, err := GetBeer(42)
	if err != nil {
		t.Fatal("err should be nil, but is:", err)
	}
	// Here's a lazy way to compare our results and our expectations.
	if !reflect.DeepEqual(beer, beerFixture) {
		t.Errorf("%#v ≠ %#v", beer, beerFixture)
	}
}
```

Run `go test`, you'll get something like this.

```
2014/10/02 01:04:05 Unstubbed query: SELECT id, name, pct FROM beer WHERE id = ? [42]
--- FAIL: TestGetBeer (0.00 seconds)
	beer_test.go:26: err should be nil, but is: mogi: query not stubbed
```
Unstubbed queries will return the above error, `mogi.ErrUnstubbed`.

### Stub it
Now we want to stub that query. Here's one way:

```
	mogi.Select("id", "name", "pct").StubCSV(`42,Yona Yona Ale,5.5`)
```
`mogi.Select(...)` takes column names and matches it against incoming queries. `StubCSV()` takes a CSV string and translates it into `driver.Value`s to answer the query.

This is cool, but it's not a very good test. It's not specific enough. For example, `GetBeer(99)` would still return Beer #42 and the test would pass. 

mogi lets us add filters our stubs to make them more specific. Let's try and match our query as specifically as possible.

#### Filter by table
`mogi.Select()` returns a `*mogi.Stub`, with it you can use method chaining to add as many filters as you like.

Since we know our query needs the `beer` table, let's express that in our test. Just add `From("beer")` to your stub.

```
	mogi.Select("id", "name", "pct").
		From("beer").
		StubCSV(`42,Yona Yona Ale,5.5`)
```

#### Filter by the WHERE clause
We still need to solve our ID problem. let's use `Where(column, value)` to add another filter to our stub. 

```
	mogi.Select("id", "name", "pct").
		From("beer").
		Where("id", 42).
		StubCSV(`42,Yona Yona Ale,5.5`)
``` 

Reads kind of like an SQL statement, doesn't it? This way, we can write stubs for specific queries.

### Multiple stubs and mogi.Reset()
You can stub as many queries as you'd like at the same time. When mogi receives a query, it checks all registered stubs in order of specificity and uses the first match. This generally does what you want, but you can manually tweak the priority of your stub with the `Priority(int)` method.

`mogi.Reset()` will clear all registered stubs. It's a good idea to call this at the end of every test. 

### Putting it all together
Now we have a thorough, passing test.

```
func TestGetBeer(t *testing.T) {
	defer mogi.Reset()
	mogi.Select("id", "name", "pct").
		From("beer").
		Where("id", 1).
		StubCSV(`42,Yona Yona Ale,5.5`)

	beer, err := GetBeer(42)
	if err != nil {
		t.Fatal("err should be nil, but is:", err)
	}
	if !reflect.DeepEqual(beer, beerFixture) {
		t.Errorf("%#v ≠ %#v", beer, beerFixture)
	}
}
```

### More ways to stub
Here's some other useful ways to register a stub.

#### Stub an error
Let's say we want to make sure a missing record returns an error. Although you can leave a query unstubbed to have it error out, it's better to stub `sql.ErrNoRows`.

```
	mogi.Select().
		From("beer").
		Where("id", 99).
		StubError(sql.ErrNoRows)
```

Notice how we didn't specify any columns in `mogi.Select()`. Specifying columns is optional: using `Select()` like this lets us match any SELECT statement.