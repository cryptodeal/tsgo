import { describe, expect, it } from 'bun:test'
import { genDisposePtr, _IntTest, _Int32ArrayTest, _Int64ArrayTest, _Float32ArrayTest, _Float64ArrayTest, _Uint32ArrayTest, _Uint64ArrayTest, _StringTest, _TestMap, arraySize, _Float32ArgTest, _Float64ArgTest, type StructBar, _Int64ArgTest, _Uint32ArgTest, _Uint64ArgTest, _Int32ArgTest, _StructBar, _TestStruct } from '@tsgo/abstract'
import { ptr, toArrayBuffer } from 'bun:ffi'

describe('tsgo', () => {
  it('returns `int`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _IntTest.native(foo)
    expect(typeof bar).toBe('number')
  })

   it('returns `*[]float32`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Float32ArrayTest.native(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float32Array(toArrayBuffer(bar,0, arraySize.native(bar) * 4, genDisposePtr.native()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

   it('returns `*[]float64`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Float64ArrayTest.native(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float64Array(toArrayBuffer(bar, 0, arraySize.native(bar) * 8, genDisposePtr.native()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

  it('returns `*[]int32`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Int32ArrayTest.native(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Int32Array(toArrayBuffer(bar, 0, arraySize.native(bar) * 4, genDisposePtr.native()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

  it('returns `*[]int64`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Int32ArrayTest.native(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigInt64Array(toArrayBuffer(bar, 0, arraySize.native(bar) * 8, genDisposePtr.native()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('bigint')
    }
  })

  it('returns `*[]uint32`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Uint32ArrayTest.native(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Uint32Array(toArrayBuffer(bar, 0, arraySize.native(bar) * 4, genDisposePtr.native()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })

  it('returns `*[]uint64`', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _Uint64ArrayTest.native(foo)
    expect(typeof bar).toBe('number')
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigUint64Array(toArrayBuffer(bar, 0, arraySize.native(bar) * 8, genDisposePtr.native()))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('bigint')
    }
  })

  
    it('returns Go struct as JSON (`json.Marshal` struct)', () => {
      let StructBar = new _StructBar(_TestStruct())
      expect(typeof StructBar).toBe('object')
      expect(typeof StructBar.Field).toBe('string')
      expect(typeof StructBar.FieldWithWeirdJSONTag).toBe('number')
      expect(typeof StructBar.FieldThatShouldBeOptional).toBe('string')
      expect(typeof StructBar.FieldThatShouldNotBeOptional).toBe('string')
      expect(typeof StructBar.FieldThatShouldBeReadonly).toBe('string')
      StructBar = null
      Bun.gc(true)
    })
  

  it('returns string (as cstring)', () => {
    const str = _StringTest().toString()
    expect(typeof str).toBe('string')
    expect(str).toBe('Hello, World!')
  })

  it('returns map (Record<number, string>)', () => {
    const str = <Record<number, string>>JSON.parse(_TestMap().toString())
    const keys = Object.keys(str)
    for (let i = 0; i < keys.length; i++) {
      expect(typeof str[keys[i]]).toBe('string')
    }
  })

  it('round trip `Float32Array`; mutate underlying data', () => {
    const test = new Float32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Float32ArgTest.native(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float32Array(toArrayBuffer(res, 0, arraySize.native(res) * 4, genDisposePtr.native()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('round trip `Float64Array`; mutate underlying data', () => {
    const test = new Float64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Float64ArgTest.native(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Float64Array(toArrayBuffer(res, 0, arraySize.native(res) * 8, genDisposePtr.native()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('round trip `Int32Array`; mutate underlying data', () => {
    const test = new Int32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Int32ArgTest(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Int32Array(toArrayBuffer(res, 0, arraySize.native(res) * 4, genDisposePtr.native()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('round trip `BigInt64Array`; mutate underlying data', () => {
    const test = new BigInt64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(v => BigInt(v)))
    const temp_ptr = ptr(test)
    const res = _Int64ArgTest.native(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigInt64Array(toArrayBuffer(res, 0, arraySize.native(res) * 8, genDisposePtr.native()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('round trip `Uint32Array`; mutate underlying data', () => {
    const test = new Uint32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    const temp_ptr = ptr(test)
    const res = _Uint32ArgTest.native(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new Uint32Array(toArrayBuffer(res, 0, arraySize.native(res) * 4, genDisposePtr.native()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })

  it('round trip `Uint64Array`; mutate underlying data', () => {
    const test = new BigUint64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(v => BigInt(v)))
    const temp_ptr = ptr(test)
    const res = _Uint64ArgTest.native(temp_ptr, test.length)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore - overload toArrayBuffer params
    const out = new BigUint64Array(toArrayBuffer(res, 0, arraySize.native(res) * 8, genDisposePtr.native()))
    for (let i = 0; i < test.length; i++) {
      expect(out[i]).toBe(test[i])
    }
  })
})