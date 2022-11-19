import { describe, expect, it } from 'bun:test'
import { disposePtr, _IntTest, _Int32ArrayTest, _Int64ArrayTest, _Float32ArrayTest, _Float64ArrayTest, _Uint32ArrayTest, _Uint64ArrayTest, ArraySize } from '@tsgo/abstract'
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
})