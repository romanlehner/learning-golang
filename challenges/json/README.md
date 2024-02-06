# JSON processing

## Problem statement

For a given json data create a function that returns the item with the highest price. Example input:

```json
items := `[
    {"title":"ItemA","price":42},
    {"title":"ItemB","price":99}
  ]`
```

Example output:

```json
    {"title":"ItemB","price":99}
```

Run tests with:

```bash
cd challenges
go test ../challenges/json/
```
