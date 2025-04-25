import { MantineProvider } from "@mantine/core";
import { render } from "@testing-library/react";
import { CreateProjectPage } from "./createProject";

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
