import { useGlobalErrorStore, type RpgBookError } from "./globalErrorStore";

describe("globalErrorStore", () => {
  it("should set an error correctly", () => {
    const error: RpgBookError = {
      title: "Test Error",
      body: "This is a test error.",
    };

    // Set the error
    useGlobalErrorStore.getState().setError(error);

    // Check if the error was set
    expect(useGlobalErrorStore.getState().currentError).toEqual(error);
  });

  it("should clear the error by setting a new one with an empty body", () => {
    // Set an initial error
    useGlobalErrorStore.getState().setError({
      title: "Initial Error",
      body: "Some error.",
    });

    // Clear the error
    useGlobalErrorStore.getState().setError({ body: "" });

    // Check if the error is cleared (or set to the new 'empty' state)
    expect(useGlobalErrorStore.getState().currentError).toEqual({ body: "" });
  });
});
