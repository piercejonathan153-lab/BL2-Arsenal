# BL2 Arsenal

A [willow2 mod manager](https://github.com/bl-sdk/willow2-mod-manager) mod for **Borderlands: The Pre-Sequel** that brings
140 Borderlands 2 weapons, shields and grenade mods over to the moon, with their real parts, stats, effects, sounds and
talking-gun lines, under new names. They drop the way they did in BL2 (boss uniques, world legendaries, E-tech,
pearlescents, Seraph gear).

Full details, console commands and the complete item list: [docs/MOD_README.md](docs/MOD_README.md).

## Requirements
- Borderlands: The Pre-Sequel (PC, Steam) with **willow2-mod-manager** v3.8 or newer installed.
- **Borderlands 2 installed on the same PC.** Textures and sounds are rebuilt locally from your own BL2 files.
  Items from a BL2 DLC only unlock if you own that DLC.

## Install with the launcher (recommended)
1. Download **BL2-Arsenal-Launcher.zip** from the [latest release](https://github.com/piercejonathan153-lab/BL2-Arsenal/releases/latest)
   and unzip it anywhere (inside your TPS folder is handy).
2. Double-click **BL2 Arsenal Launcher.cmd**.
3. Click **Install** - the launcher finds TPS through Steam, downloads the latest version and installs it.

The launcher checks for updates each time it opens (and when you click **Check for updates**), shows the patch notes,
and updates in one click. **PLAY** starts the game through Steam. It also has **Install from file...**, **Uninstall**
and **Choose game folder...** for non-standard installs. It only uses the PowerShell built into Windows 10/11.

## Manual install
Download `bl2_arsenal_vX.Y.Z.sdkmod` from the [latest release](https://github.com/piercejonathan153-lab/BL2-Arsenal/releases/latest)
and put it in `<TPS>/sdk_mods/` (remove any older `bl2_arsenal` folder or `.sdkmod` first).

## Notes
- Turning the mod off removes its items from your saves until you turn it back on.
- Co-op: every player needs the mod and the same BL2 DLC.
- `b2p drops` in the console checks every item's loot entries and odds in your game.
