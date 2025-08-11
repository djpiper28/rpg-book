interface Window {
  electron: {
    shell: (typeof import("electron"))["shell"];
    dialog: (typeof import("electron"))["dialog"];
    getSystemVersion: () => string;
    getBuildType: () => string;
  };
  getEnv: () => {
    [key: string]: string;
  };
}
