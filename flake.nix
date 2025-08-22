{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems f;
    in
    {
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              # svelte dev
              nodejs
              nodePackages.typescript
              nodePackages.typescript-language-server
              nodePackages_latest.svelte-language-server
              nodePackages_latest.svelte-check

              # go dev
              go
              gopls
              gofumpt

              go-task

              openapi-generator-cli
            ];
            shellHook = "";
          };
        }
      );
    };
}
