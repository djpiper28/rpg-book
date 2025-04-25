import { render } from "@testing-library/react";
import { CreateProjectPage } from "./createProject";
import { MantineProvider } from "@mantine/core";

describe("Create project", () => {
  it("Should render", () => {
    expect(
      render(
        <MantineProvider>
          <CreateProjectPage />
        </MantineProvider>,
      ),
    ).toBeDefined();
  });
});
