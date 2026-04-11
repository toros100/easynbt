# easynbt

Code generation tool for fast unmarshaling of the minecraft (java edition) [NBT format](https://minecraft.wiki/w/NBT_format#Binary_format) without reflection.

# Modules

* [easynbt](https://github.com/toros100/easynbt/tree/main/easynbt): The code generation tool.
* [nbt](https://github.com/toros100/easynbt/tree/main/nbt): Provides a value-based Option type, as well as definitions and helper functions used by the generated code. No external dependencies to make vendoring as easy as possible.
* [nbt/nbtcmp](https://github.com/toros100/easynbt/tree/main/nbt/nbtcmp): cmp.Option for use with the nbt.Option type and [github.com/google/go-cmp](https://github.com/google/go-cmp) ([only in tests](https://github.com/google/go-cmp/issues/373))

# Overview
Taking inspiration from [easyjson](https://github.com/mailru/easyjson), easynbt lets you define Go types matching the structure of NBT data and automatically generate NBT unmarshaling code. Ideally, this code should be faster than a general, reflection-based implementation, while still being simple and readable enough to modify and optimize by hand.

Using easynbt might look like this (cf. [easynbt/examples/readme_example/](https://github.com/toros100/easynbt/tree/main/easynbt/examples/readme_example)):

```go
//go:generate easynbt -types=Data

type Data struct {
	Hello    string `nbt:"hello"`
	Position struct {
		X int32
		Y int32
	} `nbt:"pos"`
	Numbers nbt.Option[[]int8]
}
```

Based on the underlying type of `Data`, easynbt generates the methods of the interface [nbt.Unmarshaler](https://github.com/toros100/easynbt/blob/main/nbt/nbt.go), allowing you to unmarshal NBT data with the corresponding structure from a byte slice into a value of type `*Data` using the helper function [`nbt.Unmarshal[*Data]`](https://github.com/toros100/easynbt/blob/main/nbt/nbt.go).

The input types (flag `-types`) of easynbt are named types (comma-separated), most commonly with a struct type as their underlying type. A struct type is interpreted as the payload of compound NBT, i.e. several fully formed NBT with unique names. Each field with name `N` and type `T` is interpreted as a NBT with the name `N` and the tag type determined by `T`:

1. `int8`: byte[^1] tag,
2. `int16`: short tag,
3. `int32`: int tag,
4. `int64`: long tag,
5. `float32`: float tag,
6. `float64`: double tag,
7. `string`: string tag,
8. `[]T`: list tag, with the element tag type being the tag type associated with T,
9. `struct {...}`: compound tag, with its fields interpreted the same way recursively,
10. `U` or `*U`, where `U` is a named type and `*U` implements `nbt.Unmarshaler`: as determined by the interface method `TagType`.

For 10., `U` may also be a target type of the same run of easynbt, i.e. a named type that does not yet implement `nbt.Unmarshaler`, but will do so after completion.

The input types of easynbt may also be named types with underlying types as in 1-8 above, with the expected behaviour.

The following keys may be used in struct tags:

* `nbt`: Override the expected name of the NBT implied by the field (default being the field name).
* `nbtignore`: Ignore field. Blank fields (named `_`) are always ignored.
* `nbtoptional`: Mark field as optional. By default, fields are required, i.e. when unmarshaling a compound NBT into a struct, every child tag (as determined by the fields) must be found.

For `nbtignore` and `nbtoptional`, the presence of the key is taken as a boolean flag and the provided value is ignored. Note that this tool treats these as distinct key, i.e. if you want a field to be both renamed and optional, the tag should be `` `nbt:"the_name" nbtoptional:""` `` (cf. `reflect.StructTag`).

When unmarshaling a compound payload into a struct, child tags with names that do not match any expected child tags name are ignored.

To be able to distinguish absent optional values from zero values, the package `nbt` provides the type `nbt.Option[T any]`, which may be used in struct fields and implicitly marks the field as optional. The tag type for `nbt.Option[T]` with a particular `T` is determined by `T`.
Further, the package `nbt` provides types implementing `nbt.Unmarshaler` that may be used to represent byte/int/long array tags. Because array tags are unable to be nested and are the most likely to require custom implementations (e.g. for bit unpacking in chunk data), easynbt uses Go slice types solely for list tags and there is no way to for example unmarshal a long array tag into directly into `[]int64`.


[^1]: The naming is a bit unfortunate, but in NBT data, so-called byte payloads are interpreted as signed 8 bit integers, thus the Go type `byte` would be inappropriate.


# Future work
* Currently, easynbt only supports unmarshaling. While it would be neat to also support marshaling back into bytes, the current architecture only decodes the parts that are actually needed and skips everything else for performance reasons. Presumably, you would want to re-marshal the entire NBT data (and not just the portion implicitly "selected" from a larger compound tag), which would require major changes.
* Obvious optimizations, like support for string interning, optional use of unsafe, pooling etc. Much of that could already be realized with custom `nbt.Unmarshaler` implementations, but some first-class support would be neat.
* Benchmarks and general comparison with other (Go) NBT tools, e.g. [[1]](https://pkg.go.dev/github.com/sandertv/gophertunnel/minecraft/nbt) [[2]](https://pkg.go.dev/github.com/Tnze/go-mc/nbt).
* More helpful user-facing errors in the unmarshaling code, including for example byte offsets where errors occur.
* Currently, when unmarshaling a compound payload, child tags that do not match any expected name are ignored. A flag to error on unexpected fields would be nice, but in my experience not very useful for real-world NBT data.
* Better testing, set up godoc documentation
