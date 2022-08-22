# jsonfq

> wip

## Getting Started

With query, you can:

- select object value by key
- select array element by index
- apply filter to input data (only for objects and arrays)

## Some examples

### Input

For all examples we use input data

```json
{
  "data": {
    "users": [
      {
        "id": 1,
        "name": "name1"
      },
      {
        "id": 2,
        "name": "name22"
      },
      {
        "id": 3,
        "name": "name333"
      }
    ]
  }
}
```

and run query

```go
res, _ := ParseAndSelect(input, query)
fmt.Printf("%s\n", res)
```

### Queries and results

Select object value by key

Query: `data`

```json
{
  "users": [
    {
      "id": 1,
      "name": "name1"
    },
    {
      "id": 2,
      "name": "name22"
    },
    {
      "id": 3,
      "name": "name333"
    }
  ]
}
```

Select object value by keys chain

Query: `data.users`

```json
 [
  {
    "id": 1,
    "name": "name1"
  },
  {
    "id": 2,
    "name": "name22"
  },
  {
    "id": 3,
    "name": "name333"
  }
]
```

Apply filter

Query: `data.users{int($value.id) = 2}`

```json
 [
  {
    "id": 2,
    "name": "name22"
  }
]
```

Get array element by index

Query: `data.users{int($value.id) = 2}[0]`

```json
{
  "id": 2,
  "name": "name22"
}
```

Get object field value

Query: `data.users{int($value.id) = 2}[0].name`

```json
"name22"
```

## Query Syntax

### Get object value

For get object value you can use path to element, dot-separated.
If object key contains only `0-9a-zA-Z_`, you can not use quotes.

Example input data:

```json
{
  "foo": {
    "bar": {
      "baz23": {
        "quz": 100
      }
    }
  },
  "qwe": {
    "rty space": {
      "uiop": 200
    }
  }
}
```

and query examples

```
foo.bar.baz23.quz
qwe."rty space".uiop
```

### Get array element

For get array element by index you can use habitual form `[int]`

Example input data:

```json
{
  "users": [
    {
      "id": 1
    },
    {
      "id": 2
    }
  ],
  "ab cd": [
    100,
    "s",
    true,
    null
  ]
}
```

and query examples

```
users[1].id
"ab cd"[3]

```

### Filter object or array

You can apply filter to object or array. Expression will be applied to each element of object or array.

Filter expression should return `bool` value. If filter expression returns `true`, element will be appended to result.

Expression should be enclosed with `{}` and follow query path, which returns object or array

Example input

```json
{
  "data1": [...],
  "data2": {
    "sub1": {...},
    "sub2": "some string",
    "sub3": true
  }
}
```

correct examples, because data1 is array and data2.sub1 is object

```
data1{<expression>}
data3.sub1{<expression>}
```

invalid examples, because data2.sub2 and data2.sub3 is not object or array

```
data3.sub2{<expression>}
data3.sub3{<expression>}
```

Inside filter expression you can get key/value for object elements, and index/value for array elements with use keywords `$key`, `$value` and `$index`


Example input:

```json
{
  "foo": {"id":1,"name": "foo"},
  "bar": {"id":2,"name": "bar"},
  "baz": {"id":3,"name": "baz"}
}
```

Valid queries 

```
{$value.name = "foo"}   - returns only first element
{int($value.id) = 2}    - returns second element
{int($value.id) > 1}   - returns last 2 elements
```

#### Filter functions

Inside filter expression you can use some functions. 

Functions list:

##### `int(<bytes>) int`

convert path result to int

example `data.users{int($key) > 10}`
