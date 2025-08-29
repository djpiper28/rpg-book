{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    bash
    curl
    git
    nodejs_22
    pnpm
    go
    gofumpt
    golangci-lint
    electron
    p7zip
    protobuf_27
    grpc-tools
  ];
}
