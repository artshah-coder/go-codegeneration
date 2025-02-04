It is necessary to write a code generator that finds struct methods marked with a special tag and generates the following code for them:

* HTTP wrappers for these methods
* Authorization checks
* Method checks (GET/POST)
* Parameter validation
* Filling of the method parameter struct
* Handling of unknown errors
 
The code for which the generator needs to be written is located in api.go. The code generator must work with unknown code similar to api.go.
 
The only allowed dependency is the `type ApiError struct` for error checking.
 
The code generator will handle the following struct field types:
* `int`
* `string`
 
Available validator/filler tags (`apivalidator`):
* `required` - field must not be empty (no default value)
* `paramname` - if specified, use this parameter name; otherwise, use the `lowercase` field name
* `enum` - "one of" validation
* `default` - if specified and the incoming value is empty (default), set the tag’s value
* `min` - >= X for `int`, `len(str)` >= for strings
* `max` - <= X for `int`
 
Authorization is checked by verifying the `X-Auth` header contains the value `100500`.
 
The generated code will follow this structure:
 
`ServeHTTP` - handles all methods via the multiplexer. If a method is found, it calls `handler$methodName`; otherwise, returns `404`.
`handler$methodName` - wrapper around the struct method `$methodName`. Performs all checks, outputs errors or results in `JSON` format.
`$methodName` - the actual struct method (prefixed with `apigen:api` followed by `json` metadata for method name, type, and auth requirement). Do not generate this — it already exists.
 
Code generator structure:

* Find all struct methods.
* For each method, generate parameter validation and checks in `handler$methodName`.
* Generate the `ServeHTTP` wrapper for all methods.
 
Error handling in the generator can be simplified: assume input parameters are guaranteed to be valid.

Directory structure:

* handlers_gen/codegen.go -  code generator
* handlers_gen/templates - template directory
* api.go -  input file for the generator
* main.go
* main_test.go - tests after generation

Test commands:
``` shell
# From the handlers_gen directory, generate code:  
go build codegen.go
./codegen ../api.go ../api_handlers.go
# Run tests from the parent directory: 
go test -v
```
