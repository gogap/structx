# structx

structx is iterator for golang struct

### Example


```go

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gogap/structx"
)

type Person struct {
	Name   string
	Number int
}

type Access struct {
	Persons       []*Person
	HasPermission bool      `structs:"-"`
	Password      string    `encrypt:"data"`
	LastAccessed  time.Time `structs:"omitnested"`
	X             map[string]*Person
}

func main() {

	s := &Access{
		Persons: []*Person{
			{Name: "fatih", Number: 1234567},
			{Name: "zeal", Number: 123123345},
		},
		Password:      "1234567",
		LastAccessed:  time.Now(),
		HasPermission: true,
		X:             map[string]*Person{"aaa": {Name: "fatih", Number: 1234567}},
	}

	printJSON("Before", s)

	mapv, err := structx.IterateObject(s, f1)
	if err != nil {
		fmt.Println(err)
		return
	}

	printJSON("After", s)
	printJSON("Values", mapv)

}

func printJSON(tips string, v interface{}) {
	data, _ := json.MarshalIndent(v, "", "    ")
	fmt.Println(tips + "\n" + string(data))
}

func f1(path string, field structx.Field) (err error) {

	tagVal, exist := field.Tag.Lookup("encrypt")
	if !exist {
		return
	}

	if field.Val.Kind() == reflect.String {
		if field.Val.CanSet() {
			field.Val.Set(reflect.ValueOf("This field's content is changed by f1, and the field tag is:" + tagVal))
		} else {
			err = fmt.Errorf("the path of %s could not set", path)
			return
		}
	}

	return
}

```

**result**

```json
Before
{
    "Persons": [
        {
            "Name": "fatih",
            "Number": 1234567
        },
        {
            "Name": "zeal",
            "Number": 123123345
        }
    ],
    "HasPermission": true,
    "Password": "1234567",
    "LastAccessed": "2018-12-07T10:26:20.246868+08:00",
    "X": {
        "aaa": {
            "Name": "fatih",
            "Number": 1234567
        }
    }
}
After
{
    "Persons": [
        {
            "Name": "fatih",
            "Number": 1234567
        },
        {
            "Name": "zeal",
            "Number": 123123345
        }
    ],
    "HasPermission": true,
    "Password": "This field's content is changed by f1, and the field tag is:data",
    "LastAccessed": "2018-12-07T10:26:20.246868+08:00",
    "X": {
        "aaa": {
            "Name": "fatih",
            "Number": 1234567
        }
    }
}
Values
{
    "Access.LastAccessed": {
        "Val": {},
        "Tag": "structs:\"omitnested\""
    },
    "Access.Password": {
        "Val": {},
        "Tag": "encrypt:\"data\""
    },
    "Access.Persons[0].Name": {
        "Val": {},
        "Tag": ""
    },
    "Access.Persons[0].Number": {
        "Val": {},
        "Tag": ""
    },
    "Access.Persons[1].Name": {
        "Val": {},
        "Tag": ""
    },
    "Access.Persons[1].Number": {
        "Val": {},
        "Tag": ""
    },
    "Access.X[aaa].Name": {
        "Val": {},
        "Tag": ""
    },
    "Access.X[aaa].Number": {
        "Val": {},
        "Tag": ""
    }
}
```


> some code is copied from https://github.com/fatih/structs