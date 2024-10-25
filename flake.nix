{
  description = "home-control";

  inputs = {
    nixpkgs = {
      type = "github";
      owner = "NixOS";
      repo = "nixpkgs";
      rev = "86e1ad4ec007f4f0e9561886935fe9b278860de8";
    };
    flake-utils = {
      type = "github";
      owner = "numtide";
      repo = "flake-utils";
      rev = "b1d9ab70662946ef0850d488da1c9019f3a9752a";
    };
    pre-commit-hooks = {
      type = "github";
      owner = "cachix";
      repo = "git-hooks.nix";
      rev = "c7012d0c18567c889b948781bc74a501e92275d1";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      pre-commit-hooks,
      ...
    }:
    let
      utils = flake-utils;
    in
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        checks = {
          pre-commit-check = pre-commit-hooks.lib.${system}.run {
            src = ./.;
            hooks = {
              nixfmt = {
                enable = true;
                name = "nixfmt check";
                entry = "${pkgs.nixfmt-rfc-style}/bin/nixfmt -c ";
                types = [ "nix" ];
              };
            };
          };
        };

        packages.default = pkgs.buildGoModule {
          pname = "home-control";
          version = "0.1.0";
          vendorHash = "sha256-NoNqrx5EUF2DV8Qn0dabrsbuxrWvPMwqn3MyVTzSfLY=";
          src = ./.;
        };

        devShells = {
          default = pkgs.mkShell {
            inherit (self.checks.${system}.pre-commit-check) shellHook;
            buildInputs = self.checks.${system}.pre-commit-check.enabledPackages;

            packages = with pkgs; [
              go
              golangci-lint
            ];
          };
        };
      }
    );
}
