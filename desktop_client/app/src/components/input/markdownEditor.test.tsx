import { fireEvent, waitFor } from "@testing-library/react";
import { wrappedRender } from "@/lib/testUtils/wrappedRender";
import { MarkdownEditor } from "./markdownEditor";

describe("MarkdownEditor", () => {
  it("Should render with the correct initial value", async () => {
    const label = "Test Label";
    const value = "Initial Value";
    const setValue = vi.fn();

    const { findByLabelText } = wrappedRender(
      <MarkdownEditor label={label} setValue={setValue} value={value} />,
    );

    const textarea = await findByLabelText(label);
    expect(textarea).toBeInTheDocument();
    expect((textarea as HTMLInputElement).value).toBe(value);
  });

  it("Should toggle the preview when the preview button is clicked", async () => {
    const label = "Test Label";
    const value = "## Test Header";
    const setValue = vi.fn();

    const { findByText, queryByText } = wrappedRender(
      <MarkdownEditor label={label} setValue={setValue} value={value} />,
    );

    const showButton = await findByText(/show preview/i);
    fireEvent.click(showButton);

    await waitFor(async () => {
      const previewHeader = await findByText("Test Header");
      expect(previewHeader).toBeInTheDocument();
    });

    const hideButton = await findByText(/hide preview/i);
    fireEvent.click(hideButton);

    await waitFor(() => {
      const previewHeader = queryByText("Test Header");
      expect(previewHeader).not.toBeInTheDocument();
    });
  });

  it("Should call setValue when the textarea value changes", async () => {
    const label = "Test Label";
    const value = "Initial Value";
    const setValue = vi.fn();

    const { findByLabelText } = wrappedRender(
      <MarkdownEditor label={label} setValue={setValue} value={value} />,
    );

    const textarea = await findByLabelText(label);
    fireEvent.change(textarea, { target: { value: "New Value" } });

    await waitFor(() => {
      expect(setValue).toHaveBeenCalledWith("New Value");
    });
  });
});
