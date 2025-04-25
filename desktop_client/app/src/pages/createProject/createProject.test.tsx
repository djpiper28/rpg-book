import { render } from "@testing-library/react";
import { CreateProjectPage } from "./createProject";

describe("Create project", () => {
  it("Should render", () => {
    expect(render(<CreateProjectPage />)).toBeDefined();
  });
});
