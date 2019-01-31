let
  _pkgs = import <nixpkgs> {};
in
{ pkgs ? import (_pkgs.fetchFromGitHub { owner = "NixOS";
                                         repo = "nixpkgs";
                                         rev = "18.09";
                                         sha256 = "1ib96has10v5nr6bzf7v8kw7yzww8zanxgw2qi1ll1sbv6kj6zpd";
                                       }) {}
}:

with pkgs;

stdenv.mkDerivation rec {
  name = "pkg";
  env = buildEnv { name = name; paths = buildInputs; };
  buildInputs = [
    dep
    go
  ];
}
