import { fireEvent, waitFor } from "@testing-library/react";
import { act } from "react";
import { DbExtension } from "@/lib/databaseTypes";
import { projectClient } from "@/lib/grpcClient/client";
import {
  type CreateProjectReq,
  type ProjectHandle,
} from "@/lib/grpcClient/pb/project";
import { newResult } from "@/lib/testUtils/grpcTestUtils";
import { wrappedRender } from "@/lib/testUtils/wrappedRender";
import { CreateProjectPage } from "./createProject";

vi.mock("../../lib/grpcClient/client.ts");
const mockedClient = vi.mocked(projectClient);
const id = "f23c1618-39c3-4368-a66b-c50ed7187ea6";

describe("Create project", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("Should render", () => {
    expect(wrappedRender(<CreateProjectPage />)).toBeDefined();
  });

  it("Should create a project with the provided settings", async () => {
    mockedClient.createProject.mockResolvedValueOnce(
      newResult({} as CreateProjectReq, { id } as ProjectHandle),
    );

    const dom = wrappedRender(<CreateProjectPage />);
    const testName = "Testing The Horrors";

    await act(async () => {
      const projectName = await dom.findByLabelText(/project name/i);

      fireEvent.change(projectName, {
        target: {
          value: testName,
        },
      });
    });

    await waitFor(async () => {
      const projectNameInput = (await dom.findByLabelText(
        /project name/i,
      )) as HTMLInputElement;

      expect(projectNameInput.value).toBe(testName);
    });

    await act(async () => {
      const createButton = await dom.findByText(/create project/i);
      fireEvent.click(createButton);
    });

    //eslint-disable-next-line @typescript-eslint/unbound-method
    expect(mockedClient.createProject).toHaveBeenCalledWith({
      fileName: `${testName}${DbExtension}`,
      projectName: testName,
    } as CreateProjectReq);
  });
});
