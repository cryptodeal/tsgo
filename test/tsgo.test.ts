import {
	_DemoStruct,
	_DemoStruct2,
	_DemoStruct3,
	_Float32ArgTest,
	_Float32ArrayTest,
	_Float64ArgTest,
	_Float64ArrayTest,
	_Int32ArgTest,
	_Int32ArrayTest,
	_Int64ArgTest,
	_Int64ArrayTest,
	_IntTest,
	_StringTest,
	_StructBar,
	_TestMap,
	_TestStruct,
	_TestStruct2,
	_Uint32ArgTest,
	_Uint32ArrayTest,
	_Uint64ArgTest,
	_Uint64ArrayTest,
	arraySize,
	genDisposePtr
} from '@tsgo/abstract';
import { ptr, toArrayBuffer } from 'bun:ffi';
import { describe, expect, it } from 'bun:test';

describe('tsgo', () => {
	it('returns `int`', () => {
		const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8');
		const bar = _IntTest.native(foo);
		expect(typeof bar).toBe('number');
	});

	it('returns `*[]float32`', () => {
		const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8');
		const bar = _Float32ArrayTest.native(foo);
		expect(typeof bar).toBe('number');
		const out = new Float32Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(bar, 0, arraySize.native(bar) * 4, genDisposePtr.native())
		);
		for (let i = 0; i < out.length; i++) {
			expect(typeof out[i]).toBe('number');
		}
	});

	it('returns `*[]float64`', () => {
		const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8');
		const bar = _Float64ArrayTest.native(foo);
		expect(typeof bar).toBe('number');
		const out = new Float64Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(bar, 0, arraySize.native(bar) * 8, genDisposePtr.native())
		);
		for (let i = 0; i < out.length; i++) {
			expect(typeof out[i]).toBe('number');
		}
	});

	it('returns `*[]int32`', () => {
		const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8');
		const bar = _Int32ArrayTest.native(foo);
		expect(typeof bar).toBe('number');
		const out = new Int32Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(bar, 0, arraySize.native(bar) * 4, genDisposePtr.native())
		);
		for (let i = 0; i < out.length; i++) {
			expect(typeof out[i]).toBe('number');
		}
	});

	it('returns `*[]int64`', () => {
		const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8');
		const bar = _Int32ArrayTest.native(foo);
		expect(typeof bar).toBe('number');
		const out = new BigInt64Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(bar, 0, arraySize.native(bar) * 8, genDisposePtr.native())
		);
		for (let i = 0; i < out.length; i++) {
			expect(typeof out[i]).toBe('bigint');
		}
	});

	it('returns `*[]uint32`', () => {
		const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8');
		const bar = _Uint32ArrayTest.native(foo);
		expect(typeof bar).toBe('number');
		const out = new Uint32Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(bar, 0, arraySize.native(bar) * 4, genDisposePtr.native())
		);
		for (let i = 0; i < out.length; i++) {
			expect(typeof out[i]).toBe('number');
		}
	});

	it('returns `*[]uint64`', () => {
		const foo = Buffer.from(`Message that originated from Bun.js runtime as a string!\0`, 'utf8');
		const bar = _Uint64ArrayTest.native(foo);
		expect(typeof bar).toBe('number');
		const out = new BigUint64Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(bar, 0, arraySize.native(bar) * 8, genDisposePtr.native())
		);
		for (let i = 0; i < out.length; i++) {
			expect(typeof out[i]).toBe('bigint');
		}
	});

	it('returns string (as cstring)', () => {
		const str = _StringTest().toString();
		expect(typeof str).toBe('string');
		expect(str).toBe('Hello, World!');
	});

	it('returns map (Record<number, string>)', () => {
		const str = <Record<number, string>>JSON.parse(_TestMap().toString());
		const keys = Object.keys(str);
		for (let i = 0; i < keys.length; i++) {
			expect(typeof str[keys[i]]).toBe('string');
		}
	});

	it('round trip `Float32Array`; mutate underlying data', () => {
		const test = new Float32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]);
		const temp_ptr = ptr(test);
		const res = _Float32ArgTest.native(temp_ptr, test.length);
		const out = new Float32Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(res, 0, arraySize.native(res) * 4, genDisposePtr.native())
		);
		for (let i = 0; i < test.length; i++) {
			expect(out[i]).toBe(test[i]);
		}
	});

	it('round trip `Float64Array`; mutate underlying data', () => {
		const test = new Float64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]);
		const temp_ptr = ptr(test);
		const res = _Float64ArgTest.native(temp_ptr, test.length);
		const out = new Float64Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(res, 0, arraySize.native(res) * 8, genDisposePtr.native())
		);
		for (let i = 0; i < test.length; i++) {
			expect(out[i]).toBe(test[i]);
		}
	});

	it('round trip `Int32Array`; mutate underlying data', () => {
		const test = new Int32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]);
		const temp_ptr = ptr(test);
		const res = _Int32ArgTest(temp_ptr, test.length);

		const out = new Int32Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(res, 0, arraySize.native(res) * 4, genDisposePtr.native())
		);
		for (let i = 0; i < test.length; i++) {
			expect(out[i]).toBe(test[i]);
		}
	});

	it('round trip `BigInt64Array`; mutate underlying data', () => {
		const test = new BigInt64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((v) => BigInt(v)));
		const temp_ptr = ptr(test);
		const res = _Int64ArgTest.native(temp_ptr, test.length);
		const out = new BigInt64Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(res, 0, arraySize.native(res) * 8, genDisposePtr.native())
		);
		for (let i = 0; i < test.length; i++) {
			expect(out[i]).toBe(test[i]);
		}
	});

	it('round trip `Uint32Array`; mutate underlying data', () => {
		const test = new Uint32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]);
		const temp_ptr = ptr(test);
		const res = _Uint32ArgTest.native(temp_ptr, test.length);
		const out = new Uint32Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(res, 0, arraySize.native(res) * 4, genDisposePtr.native())
		);
		for (let i = 0; i < test.length; i++) {
			expect(out[i]).toBe(test[i]);
		}
	});

	it('round trip `Uint64Array`; mutate underlying data', () => {
		const test = new BigUint64Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((v) => BigInt(v)));
		const temp_ptr = ptr(test);
		const res = _Uint64ArgTest.native(temp_ptr, test.length);

		const out = new BigUint64Array(
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore - overload toArrayBuffer params
			toArrayBuffer(res, 0, arraySize.native(res) * 8, genDisposePtr.native())
		);
		for (let i = 0; i < test.length; i++) {
			expect(out[i]).toBe(test[i]);
		}
	});

	it('returns Go *struct (wrapped class)', () => {
		const StructBar = new _StructBar(_TestStruct2());
		expect(typeof StructBar).toBe('object');
		expect(typeof StructBar.Field).toBe('string');
		expect(typeof StructBar.FieldWithWeirdJSONTag).toBe('number');
		expect(typeof StructBar.FieldThatShouldBeOptional).toBe('string');
		expect(typeof StructBar.FieldThatShouldNotBeOptional).toBe('string');
		expect(typeof StructBar.FieldThatShouldBeReadonly).toBe('string');
		expect(typeof StructBar.ArrayField).toBe('undefined');
	});

	it('returns Go struct (wrapped class); contains field w ref to dif struct', () => {
		const StructBar = new _StructBar(_TestStruct());
		expect(typeof StructBar).toBe('object');
		expect(typeof StructBar.Field).toBe('string');
		expect(typeof StructBar.FieldWithWeirdJSONTag).toBe('number');
		expect(typeof StructBar.FieldThatShouldBeOptional).toBe('string');
		expect(typeof StructBar.FieldThatShouldNotBeOptional).toBe('string');
		expect(typeof StructBar.FieldThatShouldBeReadonly).toBe('string');
		expect(StructBar.StructField instanceof _DemoStruct).toBe(true);
		expect(StructBar.StructField.ArrayField instanceof Float32Array).toBe(true);
		expect(StructBar.StructField.FieldToAnotherStruct instanceof _DemoStruct2).toBe(true);
		expect(StructBar.StructField.FieldToAnotherStruct.AnotherArray instanceof Float64Array).toBe(
			true
		);
		expect(
			StructBar.StructField.FieldToAnotherStruct.BacktoAnotherStruct instanceof _DemoStruct3
		).toBe(true);
	});

	it('works to initialize Go structs', () => {
		const AnotherArray = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
		const DemoStruct3 = _DemoStruct3.init({ AnotherArray });
		expect(DemoStruct3 instanceof _DemoStruct3).toBe(true);
		for (let i = 0; i < AnotherArray.length; i++) {
			expect(DemoStruct3.AnotherArray[i]).toBe(AnotherArray[i]);
		}
		const ArrayField = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
		const DemoStruct = _DemoStruct.init({
			ArrayField,
			FieldToAnotherStruct: { AnotherArray, BacktoAnotherStruct: { AnotherArray } }
		});
		expect(DemoStruct instanceof _DemoStruct).toBe(true);
		for (let i = 0; i < ArrayField.length; i++) {
			expect(DemoStruct.ArrayField[i]).toBe(ArrayField[i]);
		}
		expect(DemoStruct.FieldToAnotherStruct instanceof _DemoStruct2).toBe(true);
		expect(DemoStruct.FieldToAnotherStruct.AnotherArray instanceof Float64Array).toBe(true);
		for (let i = 0; i < AnotherArray.length; i++) {
			expect(DemoStruct.FieldToAnotherStruct.AnotherArray[i]).toBe(AnotherArray[i]);
		}
		expect(
			DemoStruct.FieldToAnotherStruct.BacktoAnotherStruct.AnotherArray instanceof Float32Array
		).toBe(true);
		for (let i = 0; i < AnotherArray.length; i++) {
			expect(DemoStruct.FieldToAnotherStruct.BacktoAnotherStruct.AnotherArray[i]).toBe(
				AnotherArray[i]
			);
		}
		const TestStruct = _StructBar.init({
			Field: 'hello',
			FieldWithWeirdJSONTag: 123,
			ArrayField,
			FieldThatShouldBeOptional: 'optional',
			FieldThatShouldNotBeOptional: 'not optional',
			FieldThatShouldBeReadonly: 'readonly',
			StructField: DemoStruct
		});
		expect(TestStruct instanceof _StructBar).toBe(true);
	});

	it('struct field setters work', () => {
		const StructBar = new _StructBar(_TestStruct());
		StructBar.Field = 'hello';
		expect(StructBar.Field).toBe('hello');
		StructBar.FieldWithWeirdJSONTag = 1234;
		expect(StructBar.FieldWithWeirdJSONTag).toBe(123);
	});
});
