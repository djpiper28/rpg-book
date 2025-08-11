import { fireEvent, waitFor } from "@testing-library/react";
import { CreateFileInput } from "@/components/input/createFileInput";
import { DbExtension } from "@/lib/databaseTypes";
import { getProjectClient, initializeClients } from "@/lib/grpcClient/client";
import {
  type CreateProjectReq,
  type ProjectHandle,
} from "@/lib/grpcClient/pb/project";
import { type ProjectSvcClient } from "@/lib/grpcClient/pb/project.client";
import { newResult } from "@/lib/testUtils/grpcTestUtils";
import { wrappedRender } from "@/lib/testUtils/wrappedRender";
import { Component as Page } from "./createProject";

const mockedGetProjectClient = vi.mocked(getProjectClient);

const mockedClient = {
  createProject: vi.fn(),
} as Partial<ProjectSvcClient>;

const id = "f23c1618-39c3-4368-a66b-c50ed7187ea6";

vi.mock("../../components/input/createFileInput.tsx");
const mockedCreateFileInput = vi.mocked(CreateFileInput);

describe("Create project", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  beforeEach(() => {
    initializeClients();
    mockedGetProjectClient.mockReturnValue(mockedClient);
    mockedCreateFileInput.mockReturnValue(<p>Mocked input</p>);
  });

  it("Should render", () => {
    expect(wrappedRender(<Page />)).toBeDefined();
  });

  it("Should create a project with the provided settings", async () => {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call
    mockedClient.createProject!.mockResolvedValueOnce(
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
      expect(mockedClient.createProject).toHaveBeenCalledWith({
        fileName: `${testName}${DbExtension}`,
        projectName: testName,
      } as CreateProjectReq);
    });
  });
});
