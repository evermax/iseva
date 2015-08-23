#Iseva

[![Build Status](https://img.shields.io/travis/evermax/iseva.svg?style=flat-square)](https://travis-ci.org/evermax/iseva)
[![Build Status](https://drone.io/github.com/evermax/iseva/status.png)](https://drone.io/github.com/evermax/iseva/latest)

## A JSON server for frontend developer

This project is meant to help intensive frontend development to rely less on the backend development.

That way, frontend developer can start working on a feature without having the backend ready for it, as everything can be server by this server.

So this project has been meant with this needs to decouple backend and frontend in mind.

## Basic feature
Given a json file with this format:

```
{
  "urls": {
    "/test/simple": {
      "json": {
        "field1": "",
        "field2": "",
      },
    },
    "/test/objects": {
      "json": {
        Any kind of complex JSON object is allowed here.
      },
    },
    "/test/array": {
      "json": [An array works as well here!],
    }
  }
}
```
The server will serve the `"json"` object on the url provided as a key. It will return 404 if the url ask isn't one of the urls provided.

The server answer any OPTIONS call with status 204 and the following headers:

```
Access-Control-Allow-Header: Content-Type, X-Requested-With
Access-Control-Allow-Origin: {The orgin header of the request}
Content-Type: application/json; charset=utf-8
```

## Important remark
The chain of characters `---` is reserved to seperate the list of urls part and the templating (that is following). So it shouldn't appear anywhere in the JSON file.
Later on, it might become a parameter passed to the command if needed.

## Simple templating.

If you want to randomise a bit your datas, or reuse some values that you don't want to copy-paste or might change often in a lot of places, you might want to use the templating feature.

The template JSON setup is another object that you is to be added to the JSON file.

```
{
  URL part as shown before
}
---
{
  Templating part that will be explained after
}
```


### Variables
You can specify the variables you want in that fashion inside the templating area:

```
  "variables": {
    "name1": "value1",
    "name2": "value2",
    ...
  }
```

Then you can use then inside the actual JSON area in the following way:

```
  "urls": {
    "/test/variables": {
      "json": {
        "field1": "{{.name1}}",
        "field2": "",
      },
    },
    ...
  }
```

### Functions
Currently, there are only two types of functions:

- Random functions
- Array functions

They are specified in the JSON file as the variables, using the following notations:

```
"functions": {
  "rand": {
  here the list of random functions
  },
  "array": {
  here the list of array functions
  }
}
```

#### Random functions
Random funtions are using the following notation:

```
  "name": {
    type: "",
    size: ,
    max: ,
    min:
  }
```

Three types (`type` parameter) are currently available:

- int
- float
- string


`max`and `min` parameters are for the `int`, `float` and `string` types. Those parameters are integers, and it means the values between which the random value will be assigned. It can be negative values and has the following constrain `min < max`. If it is not the case, this function will just be ignored and the program will return an error when trying to use it saying that the function is not found.
The default value for those parameters is `0`.
For int and float they represent the interval in which the random value will be  chosen, whereas for the string, it represents the interval in which the size of random string will be picked up.

`size` is only for the `string` one. It represents the size of the random string. If the `min`and `max` are provided and `min < max` the `size`parameter will be ignored.


#### Array of random value with size
```
{
  "name": {
    "type": "",
    "arraysize": ,
    "size": ,
    "max": ,
    "min":
  }
}
```
Same `type`s, usage of `size`, `max`, `min` as in the random function, same rules. They apply to the elements of the array.
The `arraysize` is a fixed value that will be the size of the array.

#### Functions usage
You call a function in the JSON part using `{{functionName}}`.
Example:

```
{
  "urls": {
    "/randomint": {
      "json": {
        "key": {{randomint}}
      }
    }
  }
}
---
{
  "functions": {
    "rand": {
      "randomint": {
        "type": "int",
        "min": 58,
        "max": 5943
      }
    }
  }
}
```
And this will return, when calling `/randomint`: `{ "key": 495 }`, for example.

## Example
Here you have a complete example of how you could work with this:
It is worth mentionning again that the variables are called like this: `.variableName`, with a dot(`.`), whereas functions are called that way: `functionName` without a dot.

```
{
    "urls": {
        "/test/variables": {
            "json": {
                "var1": "{{.variable1}}",
                "var1": "{{.variable2}}",
            }
        }
        "/test/random": {
            "json": {
                "randomInt": {{randomInt}},
                "randomFloat": {{randomFloat}},
                "randomStringSize": "{{randomStringSize}}",
                "randomStringMinMax": "{{randomStringMinMax}}"
            }
        },
        "/test/array": {
            "json": {
                "arrayInt": {{arrayInt}},
                "arrayFloat": {{arrayFloat}},
                "arrayStringSize": {{arrayStringSize}},
                "arrayStringMinMax": {{arrayStringMinMax}}
            }
        }
    }
}
---
{
    "variables": {
        "variable1": "value1",
        "variable2": "value2",
    },
    "functions": {
        "rand": {
            "randomInt": {
                "type": "int",
                "max": 10,
                "min": 1
            },
            "randomFloat": {
                "type": "float",
                "max": 5550,
                "min": 234
            },
            "randomStringSize": {
                "type": "string",
                "size": 40
            },
            "randomStringMinMax": {
                "type": "string",
                "max": 35,
                "min": 10
            }
        },
        "array": {
            "arrayInt": {
                "type": "int",
                "arraysize": 40,
                "max": 200,
                "min": 11
            },
            "arrayFloat": {
                "type": "float",
                "arraysize": 50,
                "max": 243,
                "min": 15
            },
            "arrayStringSize": {
                "type": "string",
                "arraysize": 32,
                "size": 20
            },
            "arrayStringMinMax": {
                "type": "string",
                "arraysize": 27,
                "max": 79,
                "min": 0
            }
        }
    }
}
```

## Next steps
Add the object templating to the template section.
Add a few more element to the configuration:

- port
- leading path like `/api`

## Contributions
Contributions are more than welcome, you can talk to me on Twitter via [@MaximeLasserre](https://twitter.com/MaximeLasserre) or send me an email to [maxlasserre@free.fr](mailto:maxlasserre@free.fr).
I am also on the Golang slack, @maxime, so feel free to drop by and chat!
