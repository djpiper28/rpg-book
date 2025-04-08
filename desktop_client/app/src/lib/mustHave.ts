export function mustHave(x: string | undefined): string {
  if (x) {
    return x;
  }

  throw new Error("Cannot find env vars");
}
