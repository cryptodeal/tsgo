// Code generated by tsgo. DO NOT EDIT.
/* eslint-disable */
import { dlopen, FFIType, toArrayBuffer } from 'bun:ffi';
export type Something = string | number;

//////////
// source: iota.go

export type MyIotaType = number /* int */;
export enum MyEnum {
  Zero = 0,
  One = 1,
  Two = 2,
  Four = 4,
  FourString = "four",
  AlsoFourString = "four",
  Five = 5,
  FiveAgain = 5,
  Sixteen = 16,
  Seventeen = 17
}

//////////
// source: misc.go
/*
Package level
Second line of package level comment.
*/

/**
 * Comment belonging to Foo
 */
export type Foo = string;
export type FooInt64 = number /* int64 */;
export enum FooEnum {
  /**
   * Comment for the const group declaration
   */
  ConstNumberValue = 123, // Line comment behind field with value 123
  /**
   * Individual comment for field ConstStringValue
   */
  ConstStringValue = "abc",
  /**
   * Comment for the const group declaration
   */
  ConstFooValue = "foo_const_value"
}
export const Alice = "Alice"
/**
 * Multiline comment for StructBar
 * Some more text
 */
export interface StructBar {
  /**
   * Comment for field Field of type Foo
   */
  field: Foo; // Line Comment for field Field of type Foo
  weird: number /* int64 */;
  field_that_should_be_optional?: string;
  field_that_should_not_be_optional: string;
  readonly field_that_should_be_readonly: string;
  ArrayField: number /* float32 */[];
  StructField?: DemoStruct;
}
/**
 * Another example multiline comment
 * for DemoStruct
 */
export interface DemoStruct {
  ArrayField?: number /* float32 */[];
}

//////////
// Generated config for Bun FFI

export const {
  symbols: {
    _Float32ArrayTest,
    arraySize,
    _Float64ArgTest,
    _Int64ArgTest,
    _Uint32ArgTest,
    _IntTest,
    _Float64ArrayTest,
    _Uint32ArrayTest,
    _Uint64ArrayTest,
    _Int32ArgTest,
    _TestStruct,
    _DISPOSE_Struct
    _GET_StructBar_Field,
    _GET_StructBar_FieldWithWeirdJSONTag,
    _GET_StructBar_FieldThatShouldBeOptional,
    _GET_StructBar_FieldThatShouldNotBeOptional,
    _GET_StructBar_FieldThatShouldBeReadonly,
    _GET_StructBar_ArrayField,
    _GET_StructBar_StructField,
    _GET_DemoStruct_ArrayField,
    _TestStruct2,
    _TestMap,
    genDisposePtr,
    _Int32ArrayTest,
    _Int64ArrayTest,
    _Float32ArgTest,
    _Uint64ArgTest,
    _StringTest
  }
} = dlopen(import.meta.dir + '/abstract/gen_bindings.dylib', {
  genDisposePtr: {
    returns: FFIType.ptr
  },
  _Int32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Int64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Float32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Uint64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _StringTest: {
    returns: FFIType.cstring
  },
  _Float32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  arraySize: {
    args: [FFIType.ptr],
    returns: FFIType.u64_fast
  },
  _Float64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Int64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Uint32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _TestStruct: {
    returns: FFIType.ptr
  },
  _DISPOSE_Struct: {
    args: [FFIType.ptr]
  },
  _GET_StructBar_Field: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _GET_StructBar_FieldWithWeirdJSONTag: {
    args: [FFIType.ptr],
    returns: FFIType.i64_fast
  },
  _GET_StructBar_FieldThatShouldBeOptional: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _GET_StructBar_FieldThatShouldNotBeOptional: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _GET_StructBar_FieldThatShouldBeReadonly: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _GET_StructBar_ArrayField: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_StructBar_StructField: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_DemoStruct_ArrayField: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _TestStruct2: {
    returns: FFIType.ptr
  },
  _TestMap: {
    returns: FFIType.cstring
  },
  _IntTest: {
    args: [FFIType.cstring],
    returns: FFIType.int
  },
  _Float64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Uint32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Uint64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Int32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  }
})

const registry = new FinalizationRegistry((disp: { cb: (ptr: number) => void; ptr: number}) => {
  const { cb, ptr } = disp;
  return cb(ptr);
});

export class _StructBar {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr });
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct(ptr);
  }

  get Field(): string {
    return _GET_StructBar_Field(this._ptr).toString();
  }

  get FieldWithWeirdJSONTag(): number {
    return _GET_StructBar_FieldWithWeirdJSONTag(this._ptr);
  }

  get FieldThatShouldBeOptional(): string | undefined {
    return _GET_StructBar_FieldThatShouldBeOptional(this._ptr).toString();
  }

  get FieldThatShouldNotBeOptional(): string {
    return _GET_StructBar_FieldThatShouldNotBeOptional(this._ptr).toString();
  }

  get FieldThatShouldBeReadonly(): string {
    return _GET_StructBar_FieldThatShouldBeReadonly(this._ptr).toString();
  }

  get ArrayField(): Float32Array | undefined {
    const ptr = _GET_StructBar_ArrayField(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float32Array(toArrayBuffer(ptr, 0, arraySize(ptr) * 4, genDisposePtr.native()));
  }

  get StructField(): _DemoStruct | undefined {
    const ptr = _GET_StructBar_StructField(this._ptr);
    if (!ptr) return undefined;
    return new _DemoStruct(ptr);
  }
}

export class _DemoStruct {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr });
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct(ptr);
  }

  get ArrayField(): Float32Array | undefined {
    const ptr = _GET_DemoStruct_ArrayField(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float32Array(toArrayBuffer(ptr, 0, arraySize(ptr) * 4, genDisposePtr.native()));
  }
}

