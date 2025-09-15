import { fireEvent, waitFor } from "@testing-library/react";
import { electronDialog } from "@/lib/electron";
import { wrappedRender } from "@/lib/testUtils/wrappedRender";
import { IconSelector } from "./iconSelector";

vi.mock("@/lib/electron", () => ({
  electronDialog: {
    showOpenDialog: vi.fn(),
  },
}));

const mockedElectronDialog = vi.mocked(electronDialog);

describe("IconSelector", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("Should render", () => {
    expect(
      wrappedRender(
        <IconSelector
          description="test"
          filepath={undefined}
          setFilepath={vi.fn()}
        />,
      ),
    ).toBeDefined();
  });

  it("Should update the image source when a file is selected", async () => {
    const filePath = "/path/to/test/image.jpeg";
    const setFilepath = vi.fn();

    mockedElectronDialog.showOpenDialog.mockResolvedValue({
      canceled: false,
      filePaths: [filePath],
    });

    const { findByRole } = wrappedRender(
      <IconSelector
        description="test"
        filepath={undefined}
        setFilepath={setFilepath}
      />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(setFilepath).toHaveBeenCalledWith(`file://${filePath}`);
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
        filepath={undefined}
        setFilepath={setFilepath}
      />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(setFilepath).not.toHaveBeenCalled();
    });
  });
});
