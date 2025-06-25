SimpleSchema is a programming language to define exchange formats and
generate code for different targets to encode/decode data via SDL3
stream functions.

```
[[ binding_name = "hello world" ]]
module name

[[ debug_info.name = "hello friend" ]]
type example1 struct {
  field : Type = Value

  [[ endianess = "be" ]]
  field_with_options : u32
}

type example2 enum {
  BLACK // zero by default
  WHITE = 0xFFFFFF

  [[ py_fmd.name = "red_color" ]]
  RED = 0xFF0000
}

type example3 union {
  field : Type
}

type fixed_array struct {
  values : float[4]
}

type dependant_array struct {
  count : u32
  values : byte[count]
}

type magic_prefix_example struct {
  magic : magic_string = "FMD::EXAMPLE"
}

[[ name = "vec2" ]]
type vec2 float[2]

[[ namespace = "vec2", docs = "adds two vec2 into dest" ]]
proc add(src1 : const(vec2), src2 : const(vec2), dest : vec2) -> void

type add_fn proc(int, int) -> int // typedef int add_fn(int, int)
```

