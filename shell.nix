let _pkgs = import <nixpkgs> { };
in { pkgs ? import (_pkgs.fetchFromGitHub {
  owner = "NixOS";
  repo = "nixpkgs";
  #branch@date: nixpkgs-unstable@2021-10-11
  rev = "2cdd608fab0af07647da29634627a42852a8c97f";
  sha256 = "1szv364xr25yqlljrlclv8z2lm2n1qva56ad9vd02zcmn2pimdih";
}) { } }:

with pkgs;

mkShell {
  buildInputs = [
    go
    goimports
    golangci-lint
  ];
}
