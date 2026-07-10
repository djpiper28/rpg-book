{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    bash
    curl
    git
    nodejs_22
    pnpm
    go
    gopls
    gofumpt
    golangci-lint
    goimports-reviser
    electron
    p7zip
    protolint
    protobuf_27
    protoc-gen-go
    protoc-gen-go-grpc
    grpc-tools
  ];

  shellHook = ''
    # Force npm/pnpm to use the patched NixOS electron binary instead of downloading its own
    export ELECTRON_OVERRIDE_DIST_PATH="${pkgs.electron}/bin"
  '';
}
