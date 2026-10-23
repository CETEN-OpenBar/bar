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
              # svelte dev (same versions as mise.toml)
              nodejs_24
              typescript
              typescript-language-server
              svelte-language-server
              svelte-check

              # go dev
              go_1_27
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
