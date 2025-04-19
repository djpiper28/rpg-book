import { describe, expect, it } from "vitest";
import { bytesToFriendly } from "./fileSizeHelpers";

describe("File size helpers", () => {
  it("bytes", () => {
    expect(bytesToFriendly(123)).toBe("123 bytes");
  });

  it("kilo bytes", () => {
    expect(bytesToFriendly(123_456)).toBe("123 KB");
  });

  it("mega bytes", () => {
    expect(bytesToFriendly(123_456_789)).toBe("123 MB");
  });

  it("giga bytes", () => {
    expect(bytesToFriendly(123_456_789_012)).toBe("123 GB");
  });

  it("terra bytes", () => {
    expect(bytesToFriendly(123_456_789_012_345)).toBe("123 TB");
  });
});
