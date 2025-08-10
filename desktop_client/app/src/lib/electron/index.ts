/// <reference path="./electron.d.ts" />

export const electronShell = window.electron.shell;
export const electronDialog = window.electron.dialog;
export const env = window.env;
export const getSystemVersion = () => window.electron.getSystemVersion();
export const getBuildType = () => window.electron.getBuildType();
