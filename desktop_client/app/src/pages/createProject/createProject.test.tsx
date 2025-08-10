import { fireEvent, waitFor } from "@testing-library/react";
import { CreateFileInput } from "@/components/input/createFileInput";
import { DbExtension } from "@/lib/databaseTypes";
import { getProjectClient } from "@/lib/grpcClient/client";
import {
  type CreateProjectReq,
  type ProjectHandle,
} from "@/lib/grpcClient/pb/project";
import { newResult } from "@/lib/testUtils/grpcTestUtils";
import { wrappedRender } from "@/lib/testUtils/wrappedRender";
import { CreateProjectPage as Page } from "./createProject";

vi.mock("../../lib/grpcClient/client.ts");
const mockedClient = vi.mocked(getProjectClient());
const id = "f23c1618-39c3-4368-a66b-c50ed7187ea6";

vi.mock("../../components/input/createFileInput.tsx");
const mockedCreateFileInput = vi.mocked(CreateFileInput);

describe("Create project", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  beforeEach(() => {
    mockedCreateFileInput.mockReturnValue(<p>Mocked input</p>);
  });

  it("Should render", () => {
    expect(wrappedRender(<Page />)).toBeDefined();
  });

  it("Should create a project with the provided settings", async () => {
    mockedClient.createProject.mockResolvedValueOnce(
      newResult({} as CreateProjectReq, { id } as ProjectHandle),
    );

    const dom = wrappedRender(<Page />);
    const testName = "Testing The Horrors";
    const projectNameInput = await dom.findByLabelText(/project name/i);

    fireEvent.change(projectNameInput, {
      target: { value: testName },
    });

    await waitFor(() => {
      expect((projectNameInput as HTMLInputElement).value).toBe(testName);
    });

    const createButton = await dom.findByText(/create project/i);
    fireEvent.click(createButton);

    await waitFor(() => {
      //eslint-disable-next-line @typescript-eslint/unbound-method
      expect(getProjectClient().createProject).toHaveBeenCalledWith({
        fileName: `${testName}${DbExtension}`,
        projectName: testName,
      } as CreateProjectReq);
    });
  });
});
