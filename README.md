# 🎑 TSGo

TSGo is a tool for generating Typescript typings from Golang source files that just works. At it's core, TSGo started as an opinionated fork of [gzuidhof/tygo](https://github.com/gzuidhof/tygo), which provided support for Typescript `enum` values with minimal configuration. The original implementation still exists in [cryptodeal/tygo](https://github.com/cryptodeal/tygo), but this implementation aims to extend the original implementation yet again to support optionally enabling generating a Bun FFI wrapper for the library using CGo.

It preserves comments, understands constants and also supports non-struct `type` expressions. It's perfect for generating equivalent types for a Golang REST API to be used in your front-end codebase.

**🚀 Now supports Golang 1.18 generic types**

## Installation

```shell
# Go >= 1.17
go install github.com/cryptodeal/tsgo@latest
# Go < 1.17:
go install github.com/cryptodeal/tsgo
```

## Example

_Golang input file_

```go
// Comments are kept :)
type ComplexType map[string]map[uint16]*uint32

type UserRole = string
const (
	UserRoleDefault UserRole = "viewer"
	UserRoleEditor  UserRole = "editor" // Line comments are also kept
)

type UserEntry struct {
	// Instead of specifying `tstype` we could also declare the typing
	// for uuid.NullUUID in the config file.
	ID uuid.NullUUID `json:"id" tstype:"string | null"`

	Preferences map[string]struct {
		Foo uint32 `json:"foo"`
		// An unknown type without a `tstype` tag or mapping in the config file
		// becomes `any`
		Bar uuid.UUID `json:"bar"`
	} `json:"prefs"`

	MaybeFieldWithStar *string  `json:"address"`
	Nickname           string   `json:"nickname,omitempty"`
	Role               UserRole `json:"role"`

	Complex    ComplexType `json:"complex"`
	unexported bool        // Unexported fields are omitted
	Ignored    bool        `tstype:"-"` // Fields with - are omitted too
}

type ListUsersResponse struct {
	Users []UserEntry `json:"users"`
}
```

_Typescript output_

```typescript
/**
 * Comments are kept :)
 */
export type ComplexType = {
  [key: string]: {
    [key: number /* uint16 */]: number /* uint32 */ | undefined;
  };
};
export type UserRole = string;
export const UserRoleDefault: UserRole = "viewer";
export const UserRoleEditor: UserRole = "editor"; // Line comments are also kept
export interface UserEntry {
  /**
   * Instead of specifying `tstype` we could also declare the typing
   * for uuid.NullUUID in the config file.
   */
  id: string | null;
  prefs: {
    [key: string]: {
      foo: number /* uint32 */;
      /**
       * An unknown type without a `tstype` tag or mapping in the config file
       * becomes `any`
       */
      bar: any /* uuid.UUID */;
    };
  };
  address?: string;
  nickname?: string;
  role: UserRole;
  complex: ComplexType;
}
export interface ListUsersResponse {
  users: UserEntry[];
}
```

For a real baptism by fire example, [here is a Gist with output for the Go built-in `net/http` and `time` package](https://gist.github.com/cryptodeal/7e192a2f33d8a4f5bde5b77fb2c5048c).

## Usage

### Option A: CLI (recommended)

Create a file `tsgo.yaml` in which you specify which packages are to be converted and any special type mappings you want to add.

```yaml
packages:
  - path: "github.com/cryptodeal/tsgo/examples/bookstore"
    type_mappings:
      time.Time: "string /* RFC3339 */"
      null.String: "null | string"
      null.Bool: "null | boolean"
      uuid.UUID: "string /* uuid */"
      uuid.NullUUID: "null | string /* uuid */"
```

Then run

```shell
tsgo generate
```

The output Typescript file will be next to the Go source files.

### Option B: Library-mode

```go
config := &tsgo.Config{
  Packages: []*tsgo.PackageConfig{
      &tsgo.PackageConfig{
          Path: "github.com/cryptodeal/tsgo/examples/bookstore",
      },
  },
}
gen := tsgo.New(config)
err := gen.Generate()
```

## Config

```yaml
# You can specify more than one package
packages:
  # The package path just like you would import it in Go
  - path: "github.com/my/package"

    # Where this output should be written to.
    # If you specify a folder it will be written to a file `index.ts` within that folder. By default it is written into the Golang package folder.
    output_path: "webapp/api/types.ts"

    # Customize the indentation (use \t if you want tabs)
    indent: "    "

    # Specify your own custom type translations, useful for custom types, `time.Time` and `null.String`.
    # Be default unrecognized types will be `any`.
    type_mappings:
      time.Time: "string"
      my.Type: "SomeType"

    # This content will be put at the top of the output Typescript file, useful for importing custom types.
    frontmatter: |
      import { SomeType } from '../lib/sometype.ts'

    # *Magically* converts structs that match the names listed here to typescript enums :)
    enum_structs:
      - "StructToConvert"
      - "AnotherStructToConvert"
      - "AndAnotherOne"
      - "LGTMStruct"

    # Filenames of Go source files that should not be included
    # in the output.
    exclude_files:
      - "private_stuff.go"
```

See also the source file [tsgo/config.go](./tsgo/config.go).

## Type hints through tagging

You can tag struct fields with `tstype` to specify their output Typescript type.

### Custom type mapping

```golang
// Golang input

type Book struct {
	Title    string    `json:"title"`
	Genre    string    `json:"genre" tstype:"'novel' | 'crime' | 'fantasy'"`
}
```

```typescript
// Typescript output

export interface Book {
  title: string;
  genre: "novel" | "crime" | "fantasy";
}
```

**Alternative**

You could use the `frontmatter` field in the config to inject `export type Genre = "novel" | "crime" | "fantasy"` at the top of the file, and use `tstype:"Genre"`. I personally prefer that as we may use the `Genre` type more than once.

### Required fields

Pointer type fields usually become optional in the Typescript output, but sometimes you may want to require it regardless.

You can add `,required` to the `tstype` tag to mark a pointer type as required.

```golang
// Golang input
type Nicknames struct {
	Alice   *string `json:"alice"`
	Bob     *string `json:"bob" tstype:"BobCustomType,required"`
	Charlie *string `json:"charlie" tstype:",required"`
}
```

```typescript
// Typescript output
export interface Nicknames {
  alice?: string;
  bob: BobCustomType;
  charlie: string;
}
```

### Readonly fields

Sometimes a field should be immutable, you can `,readonly` to the `tstype` tag to mark a field as `readonly`.

```golang
// Golang input
type Cat struct {
	Name    string `json:"name,readonly"`
	Owner   string `json:"owner"`
}
```

```typescript
// Typescript output
export interface Cat {
  readonly name: string;
  owner: string;
}
```

## Generics

TSGo supports generic types (Go version >= 1.18) out of the box.

```go
// Golang input
type UnionType interface {
	uint64 | string
}

type ABCD[A, B string, C UnionType, D int64 | bool] struct {
	A A `json:"a"`
	B B `json:"b"`
	C C `json:"c"`
	D D `json:"d"`
}
```

```typescript
// Typescript output
export type UnionType = number /* uint64 */ | string;

export interface ABCD<
  A extends string,
  B extends string,
  C extends UnionType,
  D extends number /* int64 */ | boolean
> {
  a: A;
  b: B;
  c: C;
  d: D;
}
```

## Related projects

- [**typescriptify-golang-structs**](https://github.com/tkrajina/typescriptify-golang-structs): Probably the most popular choice. The downside of this package is that it relies on reflection rather than parsing, which means that certain things can't be kept such as comments without adding a bunch of tags to your structs. The CLI generates a Go file which is then executed and reflected on, and its library requires you to manually specify all types that should be converted.
- [**go2ts**](https://github.com/StirlingMarketingGroup/go2ts): A transpiler with a web interface, this project was based off this project. It's perfect for quick one-off transpilations. There is no CLI, no support for `const` and there are no ways to customize the output.

If this repository was useful for your project, consider leaving a star.

## License

[MIT](./LICENSE)
