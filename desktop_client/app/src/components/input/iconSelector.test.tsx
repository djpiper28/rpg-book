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
    expect(wrappedRender(<IconSelector description="test" />)).toBeDefined();
  });

  it("Should update the image source when a file is selected", async () => {
    const filePath = "/path/to/test/image.jpeg";

    mockedElectronDialog.showOpenDialog.mockResolvedValue({
      canceled: false,
      filePaths: [filePath],
    });

    const { findByAltText, findByRole } = wrappedRender(
      <IconSelector description="test" />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(async () => {
      const img = await findByAltText("User selected");
      expect(img).toHaveAttribute("src", `file://${filePath}`);
    });
  });

  it("Should not update the image source when the dialog is canceled", async () => {
    mockedElectronDialog.showOpenDialog.mockResolvedValue({
      canceled: true,
      filePaths: [],
    });

    const { findByRole, queryByAltText } = wrappedRender(
      <IconSelector description="test" />,
    );

    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      const img = queryByAltText("User selected");
      expect(img).not.toBeInTheDocument();
    });
  });
});
