import { base64ToUint8Array, uint8ArrayToBase64 } from "./base64";

describe("Test base64 to and from", () => {
  it("Should work", () => {
    const data = new Uint8Array([1, 2, 3]);
    expect(base64ToUint8Array(uint8ArrayToBase64(data))).toEqual(data);
  });
});
