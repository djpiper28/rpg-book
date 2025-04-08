{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    bash
    curl
    git
    nodejs_22
    pnpm
    go
    golangci-lint
    electron_35
    p7zip
    protobuf_27
    protoc-gen-js
    protoc-gen-grpc-web
    grpc-tools
  ];
}
