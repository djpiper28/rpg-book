import { uint8ArrayToBase64 } from "@/lib/utils/base64";
import { fireEvent, waitFor } from "@testing-library/react";
import { getSystemClient } from "@/lib/grpcClient/client";
import { electronDialog } from "@/lib/electron";
import { wrappedRender } from "@/lib/testUtils/wrappedRender";
import { IconSelector } from "./iconSelector";

vi.mock("@/lib/electron", () => ({
  electronDialog: {
    showOpenDialog: vi.fn(),
  },
}));

vi.mock("@/lib/grpcClient/client", () => ({
  getSystemClient: vi.fn(),
}));

const mockedElectronDialog = vi.mocked(electronDialog);
const mockedGetSystemClient = vi.mocked(getSystemClient);

describe("IconSelector", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("Should render", () => {
    expect(
      wrappedRender(
        <IconSelector
          description="test"
          imageB64={undefined}
          setImageB64={vi.fn()}
        />,
      ),
    ).toBeDefined();
  });

  it("Should update the image source when a file is selected", async () => {
    const filePath = "/path/to/test/image.jpeg";
    const setImageB64 = vi.fn();
    const fakeFileBytes = new Uint8Array([1, 2, 3]);
    const fakeBase64 = uint8ArrayToBase64(fakeFileBytes);

    mockedElectronDialog.showOpenDialog.mockResolvedValue({
      canceled: false,
      filePaths: [filePath],
    });

    const mockSystemClient = {
      readFile: vi.fn().mockResolvedValue({
        response: {
          data: fakeFileBytes,
        },
      }),
    } as unknown as SystemSvcClient;

    mockedGetSystemClient.mockReturnValue(mockSystemClient);

    const { findByRole } = wrappedRender(
      <IconSelector
        description="test"
        imageB64={undefined}
        setImageB64={setImageB64}
      />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(setImageB64).toHaveBeenCalledWith(
        `data:image/jpeg;base64,${fakeBase64}`,
      );
    });
  });

  it("Should not update the image source when the dialog is canceled", async () => {
    const setFilepath = vi.fn();

    mockedElectronDialog.showOpenDialog.mockResolvedValue({
      canceled: true,
      filePaths: [],
    });

    const { findByRole } = wrappedRender(
      <IconSelector
        description="test"
        imageB64={undefined}
        setImageB64={setFilepath}
      />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(setFilepath).not.toHaveBeenCalled();
    });
  });
});
