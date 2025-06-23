SimpleSchema is a programming language to define exchange formats and
generate code for different targets to encode/decode data via SDL3
stream functions.

```
module name {
  option prop1 = "value"
  option prop2 = 123
  option prop3 = prop2 * 2 // todo
  option fmb_encoder_py {
    endianness = "be"
  }
}

type example1 struct {
  field : Type = Value
  field_with_options : u32 {
    option encoding = "utf-8"
    option lua_bindings {
      name_in_table = "field_1"
    }
  }
}

type example2 enum {
  BLACK // zero by default
  WHITE = 0xFFFFFF
  RED = 0xFF0000 {
    option lua_bindings {
      prefix = "RGB_"
    }
  }

  option lua_bindings {
    prefix = "COLOR_"
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
  // This will enforce the encoder/decoder to write/read the literal "FMF::EXAMPLE"
  magic : const String = "FMF::EXAMPLE" {
    // We dont want to actually encode/decode this data to/from lua table
    option lua_bindings.exclude = true
  }
}

type vec2 float[4] {
  /**
  * add. Adds the components from src1 to src2 and stores the result in dest.
  */
  proc add(src1 : const vec2, src2 : const vec2, dest : vec2) -> void
}

type overwrite_example struct {
  option example {
    name = "hello world!"
  }

  x : int {
    option example extend {
      // name = "hello world!"
      salute = "amigo!"
    }

    option example overwrite {
      // no name property
      salute = "hi!"
    }
  }
}
```

