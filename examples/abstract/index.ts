// Code generated by tsgo. DO NOT EDIT.
import { dlopen, FFIType } from 'bun:ffi';
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
}

//////////
// Generated config for Bun FFI

export const {
  symbols: {
    _Int64ArgTest,
    _Uint64ArgTest,
    _IntTest,
    _Float32ArgTest,
    _Uint32ArrayTest,
    _Uint64ArrayTest,
    _Uint32ArgTest,
    _Float64ArrayTest,
    _Int32ArrayTest,
    arraySize,
    _Int64ArrayTest,
    _StringTest,
    _Float64ArgTest,
    _Int32ArgTest,
    _TestStruct,
    _Float32ArrayTest,
    genDisposePtr,
    _TestMap
  }
} = dlopen(import.meta.dir + '/abstract/gen_bindings.dylib', {
  _Uint32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Uint64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Float64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Int32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Uint32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Float64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Int32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _TestStruct: {
    returns: FFIType.ptr
  },
  _Float32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  genDisposePtr: {
    returns: FFIType.ptr
  },
  arraySize: {
    args: [FFIType.ptr],
    returns: FFIType.u64_fast
  },
  _Int64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _StringTest: {
    returns: FFIType.cstring
  },
  _TestMap: {
    returns: FFIType.cstring
  },
  _IntTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Float32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Int64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Uint64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  }
})
