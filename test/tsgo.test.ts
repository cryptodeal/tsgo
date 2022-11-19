import { describe, expect, it } from 'bun:test'
import { disposePtr, _IntTest, _Int32ArrayTest, _Int64ArrayTest, _Float32ArrayTest, _Float64ArrayTest, _Uint32ArrayTest, _Uint64ArrayTest, _TestStruct, _StringTest, ArraySize, type StructBar } from '@tsgo/abstract'
import { toArrayBuffer } from 'bun:ffi'

describe('tsgo - gen CGo Code + Bindings Proof of Concept', () => {
  it('basic; should work - returns `int`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _IntTest(foo)
    expect(typeof bar).toBe('number')
  })

   it('should work - returns `*float32[]`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Float32ArrayTest(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float32Array(toArrayBuffer(bar, 0, ArraySize(bar) * 4, disposePtr()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

   it('should work - returns `*float64[]`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Float64ArrayTest(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float64Array(toArrayBuffer(bar, 0, ArraySize(bar) * 8, disposePtr()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

  it('should work - returns `*int32[]`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Int32ArrayTest(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Int32Array(toArrayBuffer(bar, 0, ArraySize(bar) * 4, disposePtr()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

  it('should work - returns `*int64[]`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Int32ArrayTest(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigInt64Array(toArrayBuffer(bar, 0, ArraySize(bar) * 8, disposePtr()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('bigint')
    }
  })

  it('should work - returns `*uint32[]`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Uint32ArrayTest(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Uint32Array(toArrayBuffer(bar, 0, ArraySize(bar) * 4, disposePtr()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

  it('should work - returns `*uint64[]`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Uint64ArrayTest(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigUint64Array(toArrayBuffer(bar, 0, ArraySize(bar) * 8, disposePtr()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('bigint')
    }
  })

  it('should work - returns Go struct as JSON (`json.Marshal` struct)', () => {
    const struct = <StructBar>JSON.parse(_TestStruct().toString())
    expect(typeof struct).toBe('object')
    const keys = Object.keys(struct)
    for (let i = 0; i < keys.length; i++) {
      switch (keys[i]) {
        case 'field':
          expect(typeof struct.field).toBe('string')
          expect(struct.field).toBe('foo')
          break
        case 'weird':
          expect(typeof struct.weird).toBe('number')
          expect(struct.weird).toBe(123)
          break
        case 'field_that_should_be_optional': 
          if (struct.field_that_should_be_optional) {
            expect(typeof struct.field_that_should_be_optional).toBe('string')
          }
          break
        case 'field_that_should_not_be_optional':
          expect(typeof struct.field_that_should_not_be_optional).toBe('string')
          expect(struct.field_that_should_not_be_optional).toBe('bar')
          break
        case 'field_that_should_be_readonly':
          expect(typeof struct.field_that_should_be_readonly).toBe('string')
          expect(struct.field_that_should_be_readonly).toBe('readonly')
          break
        default:
          console.error(`Error: field ${keys[i]} not found in struct')}`)
      }
    }
  })

  it('should work - returns string (as cstring)', () => {
    const str = _StringTest().toString()
    console.log(str)
    expect(typeof str).toBe('string')
    expect(str).toBe('Hello, World!')
  })

})