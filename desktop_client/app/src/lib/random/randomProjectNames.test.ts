import { randomPlace as randomProjectName } from "./randomProjectName";

describe("Random project names", () => {
  it("Should provide a non-empty name", () => {
    expect(randomProjectName()).not.toBe("");
  });

  it("Test a bunch", () => {
    for (let i = 0; i < 100; i++) {
      expect(randomProjectName()).not.toBe("");
    }
  });
});
