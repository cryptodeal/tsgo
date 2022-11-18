import { describe, expect, it } from 'bun:test'
import { _TestFunc, _TestFunc2, ArraySize } from '@tsgo/abstract'
import { toArrayBuffer } from 'bun:ffi'

describe('tsgo - gen CGo Code + Bindings Proof of Concept', () => {
  it('should work', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _TestFunc(foo)
    console.log('(logged from Bun) bar:', bar)
    expect(typeof bar).toBe('number')
  })

   it('should work', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _TestFunc2(foo)
    expect(typeof bar).toBe('number')
    const out = new Float32Array(toArrayBuffer(bar, 0, ArraySize(bar) * 4))
    for (let i = 0; i < out.length; i++) {
      expect(typeof out[i]).toBe('number')
    }
  })
})