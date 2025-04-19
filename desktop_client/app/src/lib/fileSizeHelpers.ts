export function bytesToFriendly(bytes: number): string {
  let exponent = 0;
  while (bytes > 1000) {
    exponent++;
    bytes /= 1000;
  }

  const size = bytes.toPrecision(3);
  switch (exponent) {
    case 0:
      return `${size} bytes`;
    case 1:
      return `${size} KB`;
    case 2:
      return `${size} MB`;
    case 3:
      return `${size} GB`;
    case 4:
      return `${size} TB`;
    default:
      throw new Error("The FILE IS HUGE!!");
  }
}
