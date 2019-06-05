let
  _pkgs = import <nixpkgs> {};
in
{ pkgs ? import (_pkgs.fetchFromGitHub { owner = "NixOS";
                                         repo = "nixpkgs-channels";
                                         # nixos-unstable @2019-05-06
                                         rev = "2ec5e9595becf05b052ce4c61a05d87ce95d19af";
                                         sha256 = "1z9ajsff9iv0n70aa4ss5nqi21m8fvs27g88lyjbh43wgjbrc2sy";
                                       }) {}
}:

with pkgs;

mkShell {
  buildInputs = [
    dep
    go
  ];
}
