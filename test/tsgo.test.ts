import { describe, expect, it } from 'bun:test'
import { genDisposePtr, _IntTest, _Int32ArrayTest, _Int64ArrayTest, _Float32ArrayTest, _Float64ArrayTest, _Uint32ArrayTest, _Uint64ArrayTest, _TestStruct, _StringTest, _TestMap, ArraySize, _Float32ArgTest, _Float64ArgTest, type StructBar, _Int64ArgTest, _Uint32ArgTest, _Uint64ArgTest, _Int32ArgTest } from '@tsgo/abstract'
import { ptr, toArrayBuffer } from 'bun:ffi'

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
    const out = new Float32Array(toArrayBuffer(bar,0, ArraySize(bar) * 4, genDisposePtr()))
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
    const out = new Float64Array(toArrayBuffer(bar, 0, ArraySize(bar) * 8, genDisposePtr()))
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
    const out = new Int32Array(toArrayBuffer(bar, 0, ArraySize(bar) * 4, genDisposePtr()))
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
    const out = new BigInt64Array(toArrayBuffer(bar, 0, ArraySize(bar) * 8, genDisposePtr()))
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
    const out = new Uint32Array(toArrayBuffer(bar, 0, ArraySize(bar) * 4, genDisposePtr()))
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
    const out = new BigUint64Array(toArrayBuffer(bar, 0, ArraySize(bar) * 8, genDisposePtr()))
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
    expect(typeof str).toBe('string')
    expect(str).toBe('Hello, World!')
  })

  it('should work - returns map (Record<number, string>)', () => {
    const str = <Record<number, string>>JSON.parse(_TestMap().toString())
    const keys = Object.keys(str)
    for (let i = 0; i < keys.length; i++) {
      expect(typeof str[keys[i]]).toBe('string')
    }
  })

  it('should work - round trip `Float32Array`; mutate underlying data', () => {
    const test = new Float32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Float32ArgTest(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float32Array(toArrayBuffer(res, 0, ArraySize(res) * 4, genDisposePtr()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('should work - round trip `Float64Array`; mutate underlying data', () => {
    const test = new Float64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Float64ArgTest(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float64Array(toArrayBuffer(res, 0, ArraySize(res) * 8, genDisposePtr()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('should work - round trip `Int32Array`; mutate underlying data', () => {
    const test = new Int32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Int32ArgTest(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Int32Array(toArrayBuffer(res, 0, ArraySize(res) * 4, genDisposePtr()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('should work - round trip `BigInt64Array`; mutate underlying data', () => {
    const test = new BigInt64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(v => BigInt(v)))
    const temp_ptr = ptr(test)
    const res = _Int64ArgTest(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigInt64Array(toArrayBuffer(res, 0, ArraySize(res) * 8, genDisposePtr()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('should work - round trip `Uint32Array`; mutate underlying data', () => {
    const test = new Uint32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Uint32ArgTest(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Uint32Array(toArrayBuffer(res, 0, ArraySize(res) * 4, genDisposePtr()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('should work - round trip `Uint64Array`; mutate underlying data', () => {
    const test = new BigUint64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(v => BigInt(v)))
    const temp_ptr = ptr(test)
    const res = _Uint64ArgTest(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigUint64Array(toArrayBuffer(res, 0, ArraySize(res) * 8, genDisposePtr()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })


  
})