import { fireEvent, waitFor } from "@testing-library/react";
import { wrappedRender } from "@/lib/testUtils/wrappedRender";
import { CreateFileInput } from "./createFileInput";
import { electronDialog } from "@/lib/electron";

vi.mock("@/lib/electron", () => ({
  electronDialog: {
    showSaveDialog: vi.fn(),
  },
}));

const mockedElectronDialog = vi.mocked(electronDialog);

describe("CreateFileInput", () => {
  const defaultProps = {
    onChange: vi.fn(),
    value: "",
    filters: [],
  };

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("Should render with the correct initial value", async () => {
    const { findByDisplayValue } = wrappedRender(
      <CreateFileInput {...defaultProps} value="initial/path/value" />,
    );
    const input = await findByDisplayValue("initial/path/value");
    expect(input).toBeInTheDocument();
  });

  it("Should call onChange with the selected file path when a file is chosen", async () => {
    const filePath = "/path/to/save/file.db";
    mockedElectronDialog.showSaveDialog.mockResolvedValue({
      canceled: false,
      filePath: filePath,
    });

    const { findByRole } = wrappedRender(<CreateFileInput {...defaultProps} />);
    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(defaultProps.onChange).toHaveBeenCalledWith(filePath);
    });
  });

  it("Should not call onChange when the save dialog is canceled", async () => {
    mockedElectronDialog.showSaveDialog.mockResolvedValue({
      canceled: true,
      filePath: undefined,
    });

    const { findByRole } = wrappedRender(<CreateFileInput {...defaultProps} />);
    const button = await findByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(defaultProps.onChange).not.toHaveBeenCalled();
    });
  });

  it("Should call onChange when the user types in the input", async () => {
    const { findByPlaceholderText } = wrappedRender(
      <CreateFileInput {...defaultProps} placeholder="test-placeholder" />,
    );
    const input = await findByPlaceholderText("test-placeholder");
    fireEvent.change(input, { target: { value: "new-value" } });

    await waitFor(() => {
      expect(defaultProps.onChange).toHaveBeenCalledWith("new-value");
    });
  });
});
