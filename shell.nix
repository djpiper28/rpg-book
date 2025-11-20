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
    protoc-gen-go
    protoc-gen-go-grpc
    grpc-tools
  ];
}
