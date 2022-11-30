// Code generated by tsgo. DO NOT EDIT.
/* eslint-disable */
import { dlopen, FFIType, ptr, toArrayBuffer } from 'bun:ffi';
export type Something = string | number;

//////////
// source: iota.go

export type MyIotaType = number /* int */;

export const Zero: IMyIotaType = Iiota;
export const One: IMyIotaType = Iiota;
export const Two: IMyIotaType = Iiota;
export const Four: IMyIotaType = Iiota;
export const FourString: Istring = "four";
export const AlsoFourString: Istring = "four";
export const Five = 5;
export const FiveAgain = 5;
export const Sixteen = Iiota + 6;
export const Seventeen = Iiota + 6;

//////////
// source: misc.go
/*
Package level
Second line of package level comment.
*/

/**
 * Comment belonging to Foo
 */
export type Foo = Istring;

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
export const ConstFooValue: IFoo = "foo_const_value";
export const Alice = "Alice";
/**
 * Multiline comment for StructBar
 * Some more text
 */
export interface IStructBar {
  /**
   * Comment for field Field of type Foo
   */
  Field: IIFoo;
  FieldWithWeirdJSONTag: Inumber /* int64 */;
  FieldThatShouldBeOptional?: IIstring;
  FieldThatShouldNotBeOptional: IIstring;
  FieldThatShouldBeReadonly: IIstring;
  ArrayField: Inumber /* float32 */[] | Float32Array;
  StructField?: IIDemoStruct | DemoStruct;
}
/**
 * Another example multiline comment
 * for DemoStruct
 */
export interface IDemoStruct {
  ArrayField?: Inumber /* float32 */[] | Float32Array;
  FieldToAnotherStruct?: IIDemoStruct2 | DemoStruct2;
}
export interface IDemoStruct2 {
  AnotherArray?: Inumber /* float64 */[] | Float64Array;
  BacktoAnotherStruct?: IIDemoStruct3 | DemoStruct3;
}
export interface IDemoStruct3 {
  AnotherArray?: Inumber /* float32 */[] | Float32Array;
}

//////////
// Generated config for Bun FFI

export const {
  symbols: {
    _Float64ArgTest,
    _Int64ArgTest,
    _TestStruct,
    _DISPOSE_Struct,
    _INIT_StructBar,
    _SET_StructBar_Field,
    _GET_StructBar_Field,
    _SET_StructBar_FieldWithWeirdJSONTag,
    _GET_StructBar_FieldWithWeirdJSONTag,
    _SET_StructBar_FieldThatShouldBeOptional,
    _GET_StructBar_FieldThatShouldBeOptional,
    _SET_StructBar_FieldThatShouldNotBeOptional,
    _GET_StructBar_FieldThatShouldNotBeOptional,
    _SET_StructBar_FieldThatShouldBeReadonly,
    _GET_StructBar_FieldThatShouldBeReadonly,
    _SET_StructBar_ArrayField,
    _GET_StructBar_ArrayField,
    _SET_StructBar_StructField,
    _GET_StructBar_StructField,
    _INIT_DemoStruct,
    _SET_DemoStruct_ArrayField,
    _GET_DemoStruct_ArrayField,
    _SET_DemoStruct_FieldToAnotherStruct,
    _GET_DemoStruct_FieldToAnotherStruct,
    _INIT_DemoStruct2,
    _SET_DemoStruct2_AnotherArray,
    _GET_DemoStruct2_AnotherArray,
    _SET_DemoStruct2_BacktoAnotherStruct,
    _GET_DemoStruct2_BacktoAnotherStruct,
    _INIT_DemoStruct3,
    _SET_DemoStruct3_AnotherArray,
    _GET_DemoStruct3_AnotherArray,
    _Int32ArgTest,
    _Uint32ArgTest,
    _TestMap,
    _Float32ArrayTest,
    _Float64ArrayTest,
    _StringTest,
    _Float32ArgTest,
    _Uint32ArrayTest,
    _Uint64ArrayTest,
    _IntTest,
    genDisposePtr,
    arraySize,
    _Int32ArrayTest,
    _Int64ArrayTest,
    _Uint64ArgTest,
    _TestStruct2,
  }
} = dlopen(import.meta.dir + '/abstract/gen_bindings.dylib', {
  _Float64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Int64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _TestStruct: {
    returns: FFIType.ptr
  },
  _DISPOSE_Struct: {
    args: [FFIType.ptr]
  },
  _INIT_StructBar: {
    args: [FFIType.ptr, FFIType.u64_fast, FFIType.ptr, FFIType.u64_fast, FFIType.ptr, FFIType.u64_fast, FFIType.ptr, FFIType.u64_fast, FFIType.ptr, FFIType.u64_fast, FFIType.ptr, FFIType.u64_fast, FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _SET_StructBar_Field: {
    args: [FFIType.ptr, FFIType.cstring],
    returns: FFIType.ptr
  },
  _GET_StructBar_Field: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _SET_StructBar_FieldWithWeirdJSONTag: {
    args: [FFIType.ptr, FFIType.i64_fast],
    returns: FFIType.ptr
  },
  _GET_StructBar_FieldWithWeirdJSONTag: {
    args: [FFIType.ptr],
    returns: FFIType.i64_fast
  },
  _SET_StructBar_FieldThatShouldBeOptional: {
    args: [FFIType.ptr, FFIType.cstring],
    returns: FFIType.ptr
  },
  _GET_StructBar_FieldThatShouldBeOptional: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _SET_StructBar_FieldThatShouldNotBeOptional: {
    args: [FFIType.ptr, FFIType.cstring],
    returns: FFIType.ptr
  },
  _GET_StructBar_FieldThatShouldNotBeOptional: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _SET_StructBar_FieldThatShouldBeReadonly: {
    args: [FFIType.ptr, FFIType.cstring],
    returns: FFIType.ptr
  },
  _GET_StructBar_FieldThatShouldBeReadonly: {
    args: [FFIType.ptr],
    returns: FFIType.cstring
  },
  _SET_StructBar_ArrayField: {
    args: [FFIType.ptr, FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _GET_StructBar_ArrayField: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _SET_StructBar_StructField: {
    args: [FFIType.ptr, FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_StructBar_StructField: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _INIT_DemoStruct: {
    args: [FFIType.ptr, FFIType.u64_fast, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _SET_DemoStruct_ArrayField: {
    args: [FFIType.ptr, FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _GET_DemoStruct_ArrayField: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _SET_DemoStruct_FieldToAnotherStruct: {
    args: [FFIType.ptr, FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_DemoStruct_FieldToAnotherStruct: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _INIT_DemoStruct2: {
    args: [FFIType.ptr, FFIType.u64_fast, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _SET_DemoStruct2_AnotherArray: {
    args: [FFIType.ptr, FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _GET_DemoStruct2_AnotherArray: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _SET_DemoStruct2_BacktoAnotherStruct: {
    args: [FFIType.ptr, FFIType.ptr],
    returns: FFIType.ptr
  },
  _GET_DemoStruct2_BacktoAnotherStruct: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _INIT_DemoStruct3: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _SET_DemoStruct3_AnotherArray: {
    args: [FFIType.ptr, FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _GET_DemoStruct3_AnotherArray: {
    args: [FFIType.ptr],
    returns: FFIType.ptr
  },
  _Uint32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _TestMap: {
    returns: FFIType.cstring
  },
  _Float32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Float64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _StringTest: {
    returns: FFIType.cstring
  },
  _Float32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Int32ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _Uint64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _IntTest: {
    args: [FFIType.cstring],
    returns: FFIType.int
  },
  genDisposePtr: {
    returns: FFIType.ptr
  },
  arraySize: {
    args: [FFIType.ptr],
    returns: FFIType.u64_fast
  },
  _Int32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Uint32ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Int64ArrayTest: {
    args: [FFIType.cstring],
    returns: FFIType.ptr
  },
  _Uint64ArgTest: {
    args: [FFIType.ptr, FFIType.u64_fast],
    returns: FFIType.ptr
  },
  _TestStruct2: {
    returns: FFIType.ptr
  },
})

const registry = new FinalizationRegistry((disp: { cb: (ptr: number) => void; ptr: number}) => {
  const { cb, ptr } = disp;
  return cb(ptr);
});

export class StructBar implements IStructBar {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  get ptr(): number {
    return this._ptr;
  }

  get Field(): string {
    return <Foo>_GET_StructBar_Field(this._ptr).toString();
  }

  set Field(val: StructBar["Field"]) {
    const parsed_val = Buffer.from(val + "\x00", "utf8");
    this._updatePtr(_SET_StructBar_Field(this._ptr, ptr(parsed_val)));
  }

  get FieldWithWeirdJSONTag(): number {
    return _GET_StructBar_FieldWithWeirdJSONTag.native(this._ptr);
  }

  set FieldWithWeirdJSONTag(val: StructBar["FieldWithWeirdJSONTag"]) {
    this._updatePtr(_SET_StructBar_FieldWithWeirdJSONTag(this._ptr, val));
  }

  get FieldThatShouldBeOptional(): string | undefined {
    return _GET_StructBar_FieldThatShouldBeOptional(this._ptr).toString();
  }

  set FieldThatShouldBeOptional(val: StructBar["FieldThatShouldBeOptional"]) {
    const parsed_val = Buffer.from(val + "\x00", "utf8");
    this._updatePtr(_SET_StructBar_FieldThatShouldBeOptional(this._ptr, ptr(parsed_val)));
  }

  get FieldThatShouldNotBeOptional(): string {
    return _GET_StructBar_FieldThatShouldNotBeOptional(this._ptr).toString();
  }

  set FieldThatShouldNotBeOptional(val: StructBar["FieldThatShouldNotBeOptional"]) {
    const parsed_val = Buffer.from(val + "\x00", "utf8");
    this._updatePtr(_SET_StructBar_FieldThatShouldNotBeOptional(this._ptr, ptr(parsed_val)));
  }

  get FieldThatShouldBeReadonly(): string {
    return _GET_StructBar_FieldThatShouldBeReadonly(this._ptr).toString();
  }

  set FieldThatShouldBeReadonly(val: StructBar["FieldThatShouldBeReadonly"]) {
    const parsed_val = Buffer.from(val + "\x00", "utf8");
    this._updatePtr(_SET_StructBar_FieldThatShouldBeReadonly(this._ptr, ptr(parsed_val)));
  }

  get ArrayField(): Float32Array | undefined {
    const ptr = _GET_StructBar_ArrayField.native(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float32Array(toArrayBuffer(ptr, 0, arraySize.native(ptr) * 4, genDisposePtr.native()));
  }

  set ArrayField(val: StructBar["ArrayField"]) {
    if (!(val instanceof Float32Array)) val = new Float32Array(val);
    this._updatePtr(_SET_StructBar_ArrayField(this._ptr, ptr(val), val.length));
  }

  get StructField(): _DemoStruct | undefined {
    const ptr = _GET_StructBar_StructField.native(this._ptr);
    if (!ptr) return undefined;
    return new _DemoStruct(ptr);
  }

  set StructField(val: StructBar["StructField"]) {
    const parsed_value = val instanceof _DemoStruct ? val : _DemoStruct.init(val);
    this._updatePtr(_SET_StructBar_StructField(this._ptr, parsed_value.ptr));
  }

  static init(struct: StructBar | IStructBar): StructBar {
    const { Field, FieldWithWeirdJSONTag, FieldThatShouldBeOptional, FieldThatShouldNotBeOptional, FieldThatShouldBeReadonly, StructField } = struct;
    let { ArrayField } = struct;
    const _Field = Buffer.from(Field + "\x00", "utf8");
    const _FieldThatShouldBeOptional = Buffer.from(FieldThatShouldBeOptional + "\x00", "utf8");
    const _FieldThatShouldNotBeOptional = Buffer.from(FieldThatShouldNotBeOptional + "\x00", "utf8");
    const _FieldThatShouldBeReadonly = Buffer.from(FieldThatShouldBeReadonly + "\x00", "utf8");
    const _StructField = StructField instanceof DemoStruct ? StructField : DemoStruct.init(StructField);
    if (!(ArrayField instanceof Float32Array)) ArrayField = new Float32Array(ArrayField);
    return new StructBar(_INIT_StructBar.native(ptr(_Field), FieldWithWeirdJSONTag, ptr(_FieldThatShouldBeOptional), ptr(_FieldThatShouldNotBeOptional), ptr(_FieldThatShouldBeReadonly), ptr(ArrayField), ArrayField.length, _StructField.ptr));
  }

  private _updatePtr(ptr: number): void {
    registry.unregister(this);
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct.native(ptr);
  }
}

export class DemoStruct implements IDemoStruct {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  get ptr(): number {
    return this._ptr;
  }

  get ArrayField(): Float32Array | undefined {
    const ptr = _GET_DemoStruct_ArrayField.native(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float32Array(toArrayBuffer(ptr, 0, arraySize.native(ptr) * 4, genDisposePtr.native()));
  }

  set ArrayField(val: DemoStruct["ArrayField"]) {
    if (!(val instanceof Float32Array)) val = new Float32Array(val);
    this._updatePtr(_SET_DemoStruct_ArrayField(this._ptr, ptr(val), val.length));
  }

  get FieldToAnotherStruct(): _DemoStruct2 | undefined {
    const ptr = _GET_DemoStruct_FieldToAnotherStruct.native(this._ptr);
    if (!ptr) return undefined;
    return new _DemoStruct2(ptr);
  }

  set FieldToAnotherStruct(val: DemoStruct["FieldToAnotherStruct"]) {
    const parsed_value = val instanceof _DemoStruct2 ? val : _DemoStruct2.init(val);
    this._updatePtr(_SET_DemoStruct_FieldToAnotherStruct(this._ptr, parsed_value.ptr));
  }

  static init(struct: DemoStruct | IDemoStruct): DemoStruct {
    const { FieldToAnotherStruct } = struct;
    let { ArrayField } = struct;
    const _FieldToAnotherStruct = FieldToAnotherStruct instanceof DemoStruct2 ? FieldToAnotherStruct : DemoStruct2.init(FieldToAnotherStruct);
    if (!(ArrayField instanceof Float32Array)) ArrayField = new Float32Array(ArrayField);
    return new DemoStruct(_INIT_DemoStruct.native(ptr(ArrayField), ArrayField.length, _FieldToAnotherStruct.ptr));
  }

  private _updatePtr(ptr: number): void {
    registry.unregister(this);
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct.native(ptr);
  }
}

export class DemoStruct2 implements IDemoStruct2 {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  get ptr(): number {
    return this._ptr;
  }

  get AnotherArray(): Float64Array | undefined {
    const ptr = _GET_DemoStruct2_AnotherArray.native(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float64Array(toArrayBuffer(ptr, 0, arraySize.native(ptr) * 8, genDisposePtr.native()));
  }

  set AnotherArray(val: DemoStruct2["AnotherArray"]) {
    if (!(val instanceof Float64Array)) val = new Float64Array(val);
    this._updatePtr(_SET_DemoStruct2_AnotherArray(this._ptr, ptr(val), val.length));
  }

  get BacktoAnotherStruct(): _DemoStruct3 | undefined {
    const ptr = _GET_DemoStruct2_BacktoAnotherStruct.native(this._ptr);
    if (!ptr) return undefined;
    return new _DemoStruct3(ptr);
  }

  set BacktoAnotherStruct(val: DemoStruct2["BacktoAnotherStruct"]) {
    const parsed_value = val instanceof _DemoStruct3 ? val : _DemoStruct3.init(val);
    this._updatePtr(_SET_DemoStruct2_BacktoAnotherStruct(this._ptr, parsed_value.ptr));
  }

  static init(struct: DemoStruct2 | IDemoStruct2): DemoStruct2 {
    const { BacktoAnotherStruct } = struct;
    let { AnotherArray } = struct;
    const _BacktoAnotherStruct = BacktoAnotherStruct instanceof DemoStruct3 ? BacktoAnotherStruct : DemoStruct3.init(BacktoAnotherStruct);
    if (!(AnotherArray instanceof Float64Array)) AnotherArray = new Float64Array(AnotherArray);
    return new DemoStruct2(_INIT_DemoStruct2.native(ptr(AnotherArray), AnotherArray.length, _BacktoAnotherStruct.ptr));
  }

  private _updatePtr(ptr: number): void {
    registry.unregister(this);
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct.native(ptr);
  }
}

export class DemoStruct3 implements IDemoStruct3 {
  private _ptr: number;

  constructor(ptr: number) {
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  get ptr(): number {
    return this._ptr;
  }

  get AnotherArray(): Float32Array | undefined {
    const ptr = _GET_DemoStruct3_AnotherArray.native(this._ptr);
    if (!ptr) return undefined;
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    return new Float32Array(toArrayBuffer(ptr, 0, arraySize.native(ptr) * 4, genDisposePtr.native()));
  }

  set AnotherArray(val: DemoStruct3["AnotherArray"]) {
    if (!(val instanceof Float32Array)) val = new Float32Array(val);
    this._updatePtr(_SET_DemoStruct3_AnotherArray(this._ptr, ptr(val), val.length));
  }

  static init(struct: DemoStruct3 | IDemoStruct3): DemoStruct3 {
    let { AnotherArray } = struct;
    if (!(AnotherArray instanceof Float32Array)) AnotherArray = new Float32Array(AnotherArray);
    return new DemoStruct3(_INIT_DemoStruct3.native(ptr(AnotherArray), AnotherArray.length));
  }

  private _updatePtr(ptr: number): void {
    registry.unregister(this);
    this._ptr = ptr;
    registry.register(this, { cb: this._gc_dispose, ptr }, this);
  }

  public _gc_dispose(ptr: number): void {
    return _DISPOSE_Struct.native(ptr);
  }
}

