export function bytesToFriendly(bytes: number): string {
  let exponent = 0;
  while (bytes > 1024) {
    exponent++;
    bytes /= 1024;
  }

  switch (bytes) {
    case 0:
      return `${bytes} bytes`;
    case 1:
      return `${bytes} KiB)`;
    case 2:
      return `${bytes} MiB)`;
    case 3:
      return `${bytes} GiB)`;
    case 4:
      return `${bytes} TiB)`;
    default:
      throw new Error("The FILE IS HUGE!!");
  }
}
