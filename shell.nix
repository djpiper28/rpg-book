{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    bash
    curl
    git
    nodejs_22
    go
    golangci-lint
    electron_35
    p7zip
    protoc-gen-js
  ];
}
