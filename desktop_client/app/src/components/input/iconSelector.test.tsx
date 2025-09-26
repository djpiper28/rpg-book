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
          imagePath={undefined}
          setImagePath={vi.fn()}
        />,
      ),
    ).toBeDefined();
  });

  it("Should update the image source when a file is selected", async () => {
    const filePath = "/path/to/test/image.jpeg";
    const setImagePath = vi.fn();

    mockedElectronDialog.showOpenDialog.mockResolvedValue({
      canceled: false,
      filePaths: [filePath],
    });

    const { findByRole } = wrappedRender(
      <IconSelector
        description="test"
        imagePath={undefined}
        setImagePath={setImagePath}
      />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(setImagePath).toHaveBeenCalledWith(filePath);
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
        imagePath={undefined}
        setImagePath={setFilepath}
      />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(setFilepath).not.toHaveBeenCalled();
    });
  });
});
