# BL2 Arsenal

**BL2 Arsenal is a mod for Borderlands: The Pre-Sequel that aims to bring almost everything from Borderlands 2 that
isn't in TPS** - gear, characters and more - working the way it did in BL2.

Included so far:
- **Weapons, shields and grenade mods:** 140 items from Borderlands 2 and its DLC that never made it into The Pre-Sequel,
  with their real parts, stats, effects, sounds and talking-gun lines, under new moon-themed names. They drop the way they
  did in BL2: boss uniques, world legendaries, E-tech, pearlescents, Seraph gear and the rainbow Effervescent tier.
- **Krieg the Psycho** as a playable seventh Vault Hunter, with his skills, Buzz Axe Rampage, class mods, heads and skins
  (the KriegTPS mod, built in with its author's permission).

More Borderlands 2 content is added over time.

No Borderlands 2 files are included: everything is rebuilt on your PC from your own Borderlands 2 install, so you get the
content of the BL2 DLC you own. It runs on [willow2-mod-manager](https://github.com/bl-sdk/willow2-mod-manager).

Full details, console commands and the complete item list: [docs/MOD_README.md](docs/MOD_README.md).

## Requirements
- Borderlands: The Pre-Sequel (PC, Steam) with **willow2-mod-manager** v3.8 or newer installed.
- **Borderlands 2 installed on the same PC.** Textures, sounds and Krieg's files are rebuilt locally from your own BL2
  files. Items from a BL2 DLC only unlock if you own that DLC; Krieg needs the Psycho Pack.

## Install with the launcher (recommended)
1. Download **[BL2ArsenalLauncher.exe](https://github.com/piercejonathan153-lab/BL2-Arsenal/raw/main/releases/BL2ArsenalLauncher.exe)**
   and put it anywhere (inside your TPS folder is handy).
2. Run it. Windows may say "Windows protected your PC" because the launcher isn't code-signed: click **More info -> Run anyway**.
3. Click **Install** - the launcher finds TPS and Borderlands 2 through Steam, downloads the latest version and installs it.

The launcher checks for updates every time it opens (and when you click **Check for updates**), shows the patch notes and
updates in one click. It shows where TPS and Borderlands 2 are installed, whether the mod manager is set up and which BL2 Arsenal version you have, and **PLAY** starts the
game through Steam. It also has **Install from file...**, **Change game folder...**, **Open mod folder** and **Uninstall**,
and it updates itself when a new launcher is released. Needs Windows 10/11 (uses the built-in Edge WebView2).

## Manual install
Download the newest `bl2_arsenal_vX.Y.Z.sdkmod` from the [releases folder](releases/)
and put it in `<TPS>/sdk_mods/` (remove any older `bl2_arsenal` folder or `.sdkmod` first).

## Notes
- Turning the mod off removes its items from your saves until you turn it back on.
- Co-op: every player needs the mod and the same BL2 DLC.
- `b2p drops` in the console checks every item's loot entries and odds in your game.
