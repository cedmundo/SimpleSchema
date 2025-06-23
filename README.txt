SimpleSchema is a programming language to define exchange formats and
generate code for different targets to encode/decode data via SDL3
stream functions.

```
module name options {
  name = "hello"
  value = "world"
}

type example1 struct {
  field : Type = Value
  field_with_options : u32 options {
    endianness = "be"
  }

  options {
    debug_info.name = "component:example1"
  }
}

type example2 enum {
  BLACK // zero by default
  WHITE = 0xFFFFFF
  RED = 0xFF0000 options {
    py_fmd.name = "red_color"
  }
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

type vec2 float[4] {
  /**
  * add. Adds the components from src1 to src2 and stores the result in dest.
  */
  proc add(src1 : const vec2, src2 : const vec2, dest : vec2) -> void
}

type add_fn proc(int, int) -> int // typedef int add_fn(int, int)
```

