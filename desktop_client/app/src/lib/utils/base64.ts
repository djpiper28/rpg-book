export function uint8ArrayToBase64(uint8Array: Uint8Array | undefined): string {
  if (!uint8Array) {
    return "";
  }

  let binaryString = "";
  for (let i = 0; i < uint8Array.length; i++) {
    binaryString += String.fromCodePoint(uint8Array[i]);
  }
  return btoa(binaryString);
}

export function base64ToUint8Array(base64String: string): Uint8Array {
  const binaryString = atob(base64String);
  const len = binaryString.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.codePointAt(i)!;
  }
  return bytes;
}
