import { EnvVarCertificate, EnvVarPort } from "@/lib/launcherTypes/index.ts";

Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// import matchers from "@testing-library/jest-dom/matchers";
// expect.extend(matchers);

process.env = {
  [EnvVarPort]: "9000",
  [EnvVarCertificate]: "Testing-cert",
};
