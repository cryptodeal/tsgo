import { describe, expect, it } from 'bun:test'
import { _TestFunc } from '@tsgo/abstract'

describe('tsgo - gen CGo Code + Bindings Proof of Concept', () => {
  it('should work', () => {
    const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8')
    const bar = _TestFunc(foo)
    console.log('(logged from Bun) bar:', bar)
    expect(typeof bar).toBe('number')
  })
})