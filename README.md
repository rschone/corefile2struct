# corefile2struct
Automatic parsing of CoreDNS Corefile into a structure based on metadata tags.

This project allows parsing of custom CoreDNS plugin configuration into a structure just by defining tags for parsing. Similarly as it works for parsing JSON files when using the standard golang library.

Supported tags:

* **cf** - defines a string literal that identifies the field in the Corefile
* **default** - states the default value which will be assigned into the field when parsing of the structure begins
* **check** - validations, that shall be performed on the field after the field stucture parsing is finished
* * nonempty - field must not be a default value of its type
* * oneOf(A|...|Z) - field value must be one of given values

Supported types:
* int
* string
* struct
* *struct
* []string
* []int
* time.Duration

