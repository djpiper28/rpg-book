import { type RpgBookError, useGlobalErrorStore } from "./globalErrorStore";

describe("globalErrorStore", () => {
  it("should set an error correctly", () => {
    const error: RpgBookError = {
      body: "This is a test error.",
      title: "Test Error",
    };

    // Set the error
    useGlobalErrorStore.getState().setError(error);

    // Check if the error was set
    expect(useGlobalErrorStore.getState().currentError).toEqual(error);
  });

  it("should clear the error by setting a new one with an empty body", () => {
    // Set an initial error
    useGlobalErrorStore.getState().setError({
      body: "Some error.",
      title: "Initial Error",
    });

    // Clear the error
    useGlobalErrorStore.getState().setError({ body: "" });

    // Check if the error is cleared (or set to the new 'empty' state)
    expect(useGlobalErrorStore.getState().currentError).toEqual({ body: "" });
  });
});
