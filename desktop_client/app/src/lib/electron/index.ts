/// <reference path="./electron.d.ts" />

export const electronShell = window.electron.shell;
export const electronDialog = window.electron.dialog;
export const getEnv = () => window.getEnv();
export const getSystemVersion = () => window.electron.getSystemVersion();
export const getBuildType = () => window.electron.getBuildType();
export const filesToOpen = () => window.electron.filesToOpen();
