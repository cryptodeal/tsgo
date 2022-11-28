// Code generated by tsgo. DO NOT EDIT.
/* eslint-disable */
import { dlopen, FFIType, ptr, toArrayBuffer } from 'bun:ffi';
export type Something = string | number;

//////////
// source: iota.go

export type MyIotaType = number /* int */;

export const Zero: MyIotaType = 0;
export const One: MyIotaType = 1;
export const Two: MyIotaType = 2;
export const Four: MyIotaType = 4;
export const FourString: string = "four";
export const AlsoFourString: string = "four";
export const Five = 5;
export const FiveAgain = 5;
export const Sixteen = 16;
export const Seventeen = 17;

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

/**
 * Comment for the const group declaration
 */
export const ConstNumberValue = 123; // Line comment behind field with value 123
/**
 * Individual comment for field ConstStringValue
 */
export const ConstStringValue = "abc";
/**
 * Comment for the const group declaration
 */
export const ConstFooValue: Foo = "foo_const_value";
export const Alice = "Alice";
/**
 * Multiline comment for StructBar
 * Some more text
 */
export interface StructBar {
  /**
   * Comment for field Field of type Foo
   */
  Field: Foo;
  FieldWithWeirdJSONTag: number /* int64 */;
  FieldThatShouldBeOptional?: string;
  FieldThatShouldNotBeOptional?: string;
  FieldThatShouldBeReadonly: string;
  ArrayField: number /* float32 */[] | Float32Array;
  StructField?: DemoStruct;
}
/**
 * Another example multiline comment
 * for DemoStruct
 */
export interface DemoStruct {
  ArrayField?: number /* float32 */[] | Float32Array;
  FieldToAnotherStruct?: DemoStruct2;
}
export interface DemoStruct2 {
  AnotherArray?: number /* float64 */[] | Float64Array;
  BacktoAnotherStruct?: DemoStruct3;
}
export interface DemoStruct3 {
  AnotherArray?: number /* float32 */[] | Float32Array;
}

//////////
// Generated config for Bun FFI

export const {
  symbols: {
    _Float32ArgTest,
    _Int64ArgTest,
    _TestStruct2,
    _DISPOSE_Struct,
    _GET_StructBar_Field,
    _GET_StructBar_FieldWithWeirdJSONTag,
    _GET_StructBar_FieldThatShouldBeOptional,
    _GET_StructBar_FieldThatShouldNotBeOptional,
    _GET_StructBar_FieldThatShouldBeReadonly,
    _GET_StructBar_ArrayField,
    _GET_StructBar_StructField,
    _GET_DemoStruct_ArrayField,
    _GET_DemoStruct_FieldToAnotherStruct,
    _GET_DemoStruct2_AnotherArray,
    _GET_DemoStruct2_BacktoAnotherStruct,
    _GET_DemoStruct3_AnotherArray,
    _IntTest,
    arraySize,
    _Float64ArrayTest,
    _StringTest,
    _TestMap,
    genDisposePtr,
    _Uint32ArrayTest,
    _Uint64ArrayTest,
    _Int32ArgTest,
    _Uint64ArgTest,
    _TestStruct,
    _Float32ArrayTest,
    _Int32ArrayTest,
    _Int64ArrayTest,
    _Float64ArgTest,
    _Uint32ArgTest
  }
} = dlopen(import.meta.dir + '/abstract/gen_bindings.dylib', {
  genDisposePtr: {
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
  },
  _Uint64ArgTest: {
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
  _GET_DemoStruct_FieldToAnotherStruct: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_DemoStruct2_AnotherArray: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_DemoStruct2_BacktoAnotherStruct: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_DemoStruct3_AnotherArray: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _Float32ArrayTest: {
    args: [FFIType.cstring],
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
  _Float64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Uint32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
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
  _TestStruct2: {
    returns: FFIType.ptr
  },
  _IntTest: {
    args: [FFIType.cstring],
    returns: FFIType.int
  },
  arraySize: {
    args: [FFIType.ptr],
    returns: FFIType.u64_fast
  },
  _Float64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _StringTest: {
    returns: FFIType.cstring
  },
  _TestMap: {
    returns: FFIType.cstring
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

  get ptr(): number {
    return this._ptr;
  }

  get Field(): string {
    return <Foo>_GET_StructBar_Field(this._ptr).toString();
  }

  get FieldWithWeirdJSONTag(): number {
    return _GET_StructBar_FieldWithWeirdJSONTag(this._ptr);
  }

  get FieldThatShouldBeOptional(): string | undefined {
    return _GET_StructBar_FieldThatShouldBeOptional(this._ptr).toString();
  }

  get FieldThatShouldNotBeOptional(): string | undefined {
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
  static init(struct: StructBar): _StructBar {
    const { Field, FieldWithWeirdJSONTag, FieldThatShouldBeOptional, FieldThatShouldNotBeOptional, FieldThatShouldBeReadonly } = struct;
    let { ArrayField, StructField } = struct;
const _Field = Buffer.from(Field + '/0', "utf8");
const _FieldThatShouldBeOptional = Buffer.from(FieldThatShouldBeOptional + '/0', "utf8");
const _FieldThatShouldNotBeOptional = Buffer.from(FieldThatShouldNotBeOptional + '/0', "utf8");
const _FieldThatShouldBeReadonly = Buffer.from(FieldThatShouldBeReadonly + '/0', "utf8");
    if (!(ArrayField instanceof Float32Array)) ArrayField = new Float32Array(ArrayField);
    if (!(StructField instanceof _DemoStruct)) StructField = _DemoStruct.init(StructField);
    return new _StructBar(DUMMY_INIT_FN_NAME(ptr(_Field), FieldWithWeirdJSONTag, ptr(_FieldThatShouldBeOptional), ptr(_FieldThatShouldNotBeOptional), ptr(_FieldThatShouldBeReadonly), ptr(ArrayField), ptr(StructField)));
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct(ptr);
  }

}

export class _DemoStruct {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr });
  }

  get ptr(): number {
    return this._ptr;
  }

  get ArrayField(): Float32Array | undefined {
    const ptr = _GET_DemoStruct_ArrayField(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float32Array(toArrayBuffer(ptr, 0, arraySize(ptr) * 4, genDisposePtr.native()));
  }

  get FieldToAnotherStruct(): _DemoStruct2 | undefined {
    const ptr = _GET_DemoStruct_FieldToAnotherStruct(this._ptr);
    if (!ptr) return undefined;
    return new _DemoStruct2(ptr);
  }
  static init(struct: DemoStruct): _DemoStruct {
    let { ArrayField, FieldToAnotherStruct } = struct;
    if (!(ArrayField instanceof Float32Array)) ArrayField = new Float32Array(ArrayField);
    if (!(FieldToAnotherStruct instanceof _DemoStruct2)) FieldToAnotherStruct = _DemoStruct2.init(FieldToAnotherStruct);
    return new _DemoStruct(DUMMY_INIT_FN_NAME(ptr(ArrayField), ptr(FieldToAnotherStruct)));
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct(ptr);
  }

}

export class _DemoStruct2 {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr });
  }

  get ptr(): number {
    return this._ptr;
  }

  get AnotherArray(): Float64Array | undefined {
    const ptr = _GET_DemoStruct2_AnotherArray(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float64Array(toArrayBuffer(ptr, 0, arraySize(ptr) * 8, genDisposePtr.native()));
  }

  get BacktoAnotherStruct(): _DemoStruct3 | undefined {
    const ptr = _GET_DemoStruct2_BacktoAnotherStruct(this._ptr);
    if (!ptr) return undefined;
    return new _DemoStruct3(ptr);
  }
  static init(struct: DemoStruct2): _DemoStruct2 {
    let { AnotherArray, BacktoAnotherStruct } = struct;
    if (!(AnotherArray instanceof Float64Array)) AnotherArray = new Float64Array(AnotherArray);
    if (!(BacktoAnotherStruct instanceof _DemoStruct3)) BacktoAnotherStruct = _DemoStruct3.init(BacktoAnotherStruct);
    return new _DemoStruct2(DUMMY_INIT_FN_NAME(ptr(AnotherArray), ptr(BacktoAnotherStruct)));
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct(ptr);
  }

}

export class _DemoStruct3 {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr });
  }

  get ptr(): number {
    return this._ptr;
  }

  get AnotherArray(): Float32Array | undefined {
    const ptr = _GET_DemoStruct3_AnotherArray(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float32Array(toArrayBuffer(ptr, 0, arraySize(ptr) * 4, genDisposePtr.native()));
  }
  static init(struct: DemoStruct3): _DemoStruct3 {
    let { AnotherArray } = struct;
    if (!(AnotherArray instanceof Float32Array)) AnotherArray = new Float32Array(AnotherArray);
    return new _DemoStruct3(DUMMY_INIT_FN_NAME(ptr(AnotherArray)));
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct(ptr);
  }

}

