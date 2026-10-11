# BL2 Arsenal (for Borderlands: The Pre-Sequel)

Brings 140 Borderlands 2 weapons, shields and grenade mods that never made it into TPS over to the moon:
their real parts, stats, behaviours, skins and effects, under new names.

## Requirements
- Borderlands: The Pre-Sequel (PC) with the **willow2 mod manager / Python SDK** (v3.8 or newer).
- **Borderlands 2 installed on the same PC.** Textures are rebuilt locally from your own BL2 files the first
  time the mod runs (takes a few seconds). Nothing from BL2 is redistributed beyond gameplay data.
- Items from BL2 DLC only unlock if that DLC is installed in your Borderlands 2 copy.

## Install
Put `bl2_arsenal.sdkmod` in `<TPS>/sdk_mods/`. Start the game. The mod enables itself at the main menu.
If your Borderlands 2 lives somewhere unusual, create `sdk_mods/settings/bl2_arsenal/config.json` with
`{"bl2_path": "D:/SteamLibrary/steamapps/common/Borderlands 2"}`.

## How you get them
Each item drops the way it did in BL2, translated to TPS:
- BL2 world-drop legendaries -> the same TPS legendary pools, at BL2's weights.
- E-tech -> world drops at very-rare tier. Pearlescents -> rare world drops at legendary tier.
- Boss / badass / mission uniques -> a matching TPS boss (Terramorphous gear -> Iwajira, Hector's gun -> Colonel Zarpedon, ...).
- Seraph gear -> the Invincible Empyrean Sentinel raid, and Seraph-vendor stock also turns up in Moonstone loot.
- BL2 vendor items -> TPS vending machines (item of the day).

Drop weights use the same scale as TPS's own loot. In any one loot list, the BL2 items together never outweigh what
that list already drops, so a boss still drops its own TPS legendary at least half the time. `b2p drops` shows the
exact odds for every item.

Rarity colours: Seraph pink, Pearlescent cyan and E-tech magenta use TPS's existing tiers; Effervescent items (the 11 rainbow
items from Commander Lilith) get their own animated rainbow tier (rarity level 507). TPS already uses BL2's
rainbow level 506 for its own orange tier, which is why they used to show up as legendaries or other tiers. The colour cycles live on the item card (inventory and ground pickups), the backpack tiles, the equipped slots, the loot beam on the ground and the weapon-type icon shown over loot you aim at. While you hold a rainbow gun, the manufacturer header on the HUD weapon card cycles too. The colour starts cycling the moment each of these appears.

## Console commands
- `b2p list` - every item and whether it is enabled for you
- `b2p give <BL2 or TPS name> [level]` - spawn one into your backpack (e.g. `b2p give Svetlana`).
  Punctuation and spaces don't matter, so leave out apostrophes: `b2p give Zarpedons Paradise`
- `b2p drops` - check every item's loot entries in the game's live loot lists and its chance per roll; writes
  `sdk_mods/settings/bl2_arsenal/drops_report.txt` (read-only: spawns nothing, changes nothing)
- `b2p ammo` - refill all ammo and grenades
- `b2p sound <event>` - play one of the ported BL2 sounds (e.g. `b2p sound Play_VOCT_HotMama_Reload`)
- `b2p giveall` - spawn every enabled item
- `b2p diag` - write a diagnostics report to `sdk_mods/settings/bl2_arsenal/diagnostics.json`
- `b2p test` - developer self-test: gives, equips and fires/throws/equips every item in turn and records the
  result in `sdk_mods/settings/bl2_arsenal/autotest.json` (`b2p test stop` to abort). It adds the items to your
  backpack, so back up your save first.

## NEIGHPLOSION!!! (Unicornsplosion)

BL2's flying Butt Stallion and its rainbow-poop chunks can't be rendered in TPS, so the shot flies TPS's own
Pandora Unicorn (from the Claptastic Voyage DLC) at BL2's size, drawn with TPS's Butt Stallion material, and the
poop chunks rain down as mini unicorns. BL2's rainbow puke decal is rebuilt on a TPS decal material with a puddle
texture made for this mod.

## Minecraft guns

BL2 draws the Pixel Pummeler (Blockhead) pellets and the Zero-G Longshot (Longbow) arrow with two tiny BL2 models.
The pellets use TPS's cube, sized to BL2's blocks, with BL2's Minecraft fire texture. The arrow is BL2's own 3D voxel
arrow model, converted at build time into the vertex format TPS uses, with BL2's arrow colours.

## Sound
BL2's own sounds that TPS doesn't have (the talking guns - Hot Mama, Ego Tripper/Overcompensator, Aegis-1340,
Feedback Loop, Loose Lips' guilt lines - plus DLC weapon/grenade/shield sounds and the Tiny Tina spell/drunk lines)
are rebuilt from your own Borderlands 2 install each time the game starts and loaded straight into TPS's audio
engine. Nothing is written to the game folder. Voice lines are English only.


### Subtitles

The talking-gun lines show BL2's subtitles (turn Subtitles on in TPS's options). On each start the mod builds
`B2P_*.<lang>` subtitle files in TPS's `WillowGame/Localization/<LANG>` folders from your own BL2 install's subtitle
files (English, German, Spanish, French, Italian and Japanese; BL2 has none for Russian and only some for
Korean/Chinese). They are new files named after the mod's packages, so none of TPS's own files are changed; delete the
`B2P_*` files there if you remove the mod.

Captain Scarlett's line for Greed (Lunar Larceny) plays the first time you equip it each session, as in BL2. The Tiny
Tina DLC's player-character shouts (spell grenades, Grog Nozzle, Crit) are left silent, since TPS's characters never
recorded them.

## Known limitations
- A handful of BL2-only effect materials (some E-tech beams, slag smoke, Unicornsplosion's fireworks) render with the
  closest TPS effect material, because BL2's compiled shaders don't run in TPS.
- Ported items are saved through their own asset-library set: if you disable or remove the mod, they disappear from saves.
- Co-op: every player needs the mod and the same BL2 DLC.
- BL2's PhysX debris particles (sparks/rock bits on a few grenade and shield effects) are left out.

## Item list
| BL2 name | TPS name | Type | Rarity | BL2 DLC needed | Where to get it |
|---|---|---|---|---|---|
| Boom Puppy | Pupulsion | Assault Rifle | Rare | Mr. Torgue | Rooster Booster |
| CHOPPER | The Bogan Mower | Assault Rifle | Rare | Sir Hammerlock | Swagman |
| Damned Cowboy | Vacuum Vaquero | Assault Rifle | Rare | Sir Hammerlock | Drongo Bones |
| Evil Smasher | Malice Masher | Assault Rifle | Rare | — | Meg |
| Kitten | Purrpetrator | Assault Rifle | Rare | Mr. Torgue | RedBelly |
| Rapier | Spacebuckler | Assault Rifle | Rare | Captain Scarlett | Red |
| Scorpio | Scorpio Zero | Assault Rifle | Rare | — | Corporal Bob |
| Stinkpot | Toxic Payload | Assault Rifle | Rare | Captain Scarlett | Bosun |
| Stomper | Crater Stomper | Assault Rifle | Rare | — | Bruce |
| Lead Storm | Regolith Rain | Assault Rifle | Seraph | Sir Hammerlock | Invincible Empyrean Sentinel (raid) |
| Seeker | Pursuit Protocol | Assault Rifle | Seraph | Tiny Tina (Dragon Keep) | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Seraphim | Archangel-0 | Assault Rifle | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| M2828 Thumpson | Tommy Moonshot | Assault Rifle | Legendary | Commander Lilith | World drop (legendary pool) |
| Madhous! | The Lunatic Ward | Assault Rifle | Legendary | — | World drop (legendary pool) |
| Ogre | Ogreload | Assault Rifle | Legendary | Tiny Tina (Dragon Keep) | Even-More-Disgusting Tork |
| Bearcat | Badgernaut | Assault Rifle | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| Bekah | Bekah Prime | Assault Rifle | Pearlescent | UVHM Pack 1 | Pearlescent world drop (very rare/legendary pools) |
| Sawbar | Orbit Saw | Assault Rifle | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| Peak Opener | Cratercracker | Assault Rifle | Effervescent | Commander Lilith | World drop (legendary pool) |
| Toothpick | Moonspike | Assault Rifle | Effervescent | Commander Lilith | Giant Shuggurath of Ice |
| Greed | Lunar Larceny | Pistol | Rare | Captain Scarlett | Nel |
| Grog Nozzle | Moonshine Nozzle | Pistol | Rare | Tiny Tina (Dragon Keep) | Lazlo |
| Judge | The Magistrate | Pistol | Rare | — | X-STLK-23 |
| Law | Sheriff's Decree | Pistol | Rare | — | Deadlift |
| Little Evie | E.V.I.E. | Pistol | Rare | Captain Scarlett | Meg |
| Pocket Rocket | Micro Missile | Pistol | Rare | Mr. Torgue | Vending machine item of the day |
| Rubi | Scarlet Sip | Pistol | Rare | — | Nel |
| Teapot | Corro-Kettle | Pistol | Rare | — | Bruce |
| Tinderbox | Oxygen Optional | Pistol | Rare | — | Flame Knuckle |
| Veritas | Brutal Truth | Pistol | Rare | — | Red Shirt |
| Dahlminator | DAHL-ETE | Pistol | E-tech | — | Corporal Bob |
| Devastator | Depopulator | Pistol | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Infection | Quarantine Breach | Pistol | Seraph | Sir Hammerlock | Invincible Empyrean Sentinel (raid) |
| Stinger | Needletail | Pistol | Seraph | Tiny Tina (Dragon Keep) | Invincible Empyrean Sentinel (raid) |
| Hector's Paradise | Zarpedon's Paradise | Pistol | Legendary | Commander Lilith | Colonel Zarpedon |
| Infinity | Event Horizon | Pistol | Legendary | — | World drop (legendary pool) |
| Stalker | Shadow Tracker | Pistol | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| Unforgiven | No Pardons | Pistol | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| Wanderlust | Wanderer Zero | Pistol | Pearlescent | UVHM Pack 1 | Pearlescent world drop (very rare/legendary pools) |
| Bane | Feedback Loop | SMG | Rare | — | Lazlo |
| Chulainn | Double Jeopardy | SMG | Rare | — | Bosun |
| Commerce | Hostile Takeover | SMG | Rare | — | X-STLK-23 |
| Crit | Butterfingers | SMG | Rare | Tiny Tina (Dragon Keep) | RedBelly |
| Lascaux | Cave Graffiti | SMG | Rare | — | Drongo Bones |
| Orc | Gruntling | SMG | Rare | Tiny Tina (Dragon Keep) | Tork Dredger |
| Sand Hawk | Regolith Raptor | SMG | Rare | Captain Scarlett | Red |
| Yellow Jacket | Stingray Swarm | SMG | E-tech | Sir Hammerlock | Boomer |
| Actualizer | Reality Check | SMG | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Florentine | Twin Eclipse | SMG | Seraph | Tiny Tina (Dragon Keep) | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Tattler | Loose Lips | SMG | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Emperor | The Regent | SMG | Legendary | — | World drop (legendary pool) |
| Avenger | Payback Program | SMG | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| Infection Cleaner | Antivirus.exe | SMG | Effervescent | Commander Lilith | Eghood |
| Nirvana | Astral Bliss | SMG | Effervescent | Commander Lilith | Tiny Destroyer |
| Blockhead | Pixel Pummeler | Shotgun | Rare | — | Badass Kraggon |
| Dog | Scrap Hound | Shotgun | Rare | — | X-STLK-23 |
| Hydra | Many-Headed Mayhem | Shotgun | Rare | Sir Hammerlock | Meg |
| Jolly Roger | Void Corsair | Shotgun | Rare | Captain Scarlett | Bosun |
| Orphan Maker | Deadbeat Dad | Shotgun | Rare | Captain Scarlett | Red |
| RokSalt | Rocksalt Reckoning | Shotgun | Rare | — | Rabid Adams |
| Shotgun 1340 | Boomstick 0-G | Shotgun | Rare | — | Bruce |
| Teeth of Terramorphous | Iwajira's Fangs | Shotgun | Rare | — | Iwajira |
| Triquetra | Three's a Crowd | Shotgun | Rare | — | Nel |
| Twister | Dust Devil | Shotgun | Rare | Sir Hammerlock | Flame Knuckle |
| Slow Hand | Steady Now | Shotgun | E-tech | Mr. Torgue | Rooster Booster |
| SWORDSPLOSION!!! | SABERSPLOSION!!! | Shotgun | E-tech | Tiny Tina (Dragon Keep) | Lazlo |
| Interfacer | Facebreaker Protocol | Shotgun | Seraph | Sir Hammerlock | Invincible Empyrean Sentinel (raid) |
| Omen | Bad Moon Omen | Shotgun | Seraph | Tiny Tina (Dragon Keep) | Invincible Empyrean Sentinel (raid) |
| Retcher | The Hurlitzer | Shotgun | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Deliverance | Express Delivery | Shotgun | Legendary | — | World drop (legendary pool) |
| Overcompensator | Ego Tripper | Shotgun | Legendary | Commander Lilith | World drop (legendary pool) |
| Butcher | Bot Butcher | Shotgun | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| Carnage | Full System Failure | Shotgun | Pearlescent | UVHM Pack 1 | Pearlescent world drop (very rare/legendary pools) |
| Unicornsplosion | NEIGHPLOSION!!! | Shotgun | Effervescent | Commander Lilith | World drop (legendary pool) |
| Morningstar | First Light | Sniper Rifle | Rare | — | Red Shirt |
| Pimpernel | Regolith Rose | Sniper Rifle | Rare | Captain Scarlett | Meg |
| Hawk Eye | Orbital Eyesight | Sniper Rifle | Seraph | Sir Hammerlock | Invincible Empyrean Sentinel (raid) |
| Patriot | Legion's Pride | Sniper Rifle | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Longbow | Zero-G Longshot | Sniper Rifle | Legendary | — | Badass Kraggon |
| Lyudmila | Svetlana | Sniper Rifle | Legendary | — | World drop (legendary pool) |
| Godfinger | Titan's Touch | Sniper Rifle | Pearlescent | UVHM Pack 1 | Pearlescent world drop (very rare/legendary pools) |
| Storm | Ion Tempest | Sniper Rifle | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| Hot Mama | Helios Hottie | Sniper Rifle | Effervescent | Commander Lilith | Raum-Kampfjet Mark V |
| 12 Pounder | Zero-G Dozen | Rocket Launcher | Rare | Captain Scarlett | Bosun |
| Creamer | Frothy Fallout | Rocket Launcher | Rare | — | Meg |
| Hive | Stinger Nest | Rocket Launcher | Rare | — | Giant Shuggurath of Ice |
| Roaster | Barbecue Protocol | Rocket Launcher | Rare | — | Flame Knuckle |
| Ahab | Space Whale | Rocket Launcher | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Bunny | Moon Hare | Rocket Launcher | Legendary | — | World drop (legendary pool) |
| Norfleet | Nova Armada | Rocket Launcher | Legendary | — | Invincible Empyrean Sentinel (raid) |
| Tunguska | Regolith Event | Rocket Launcher | Pearlescent | UVHM Pack 2 (Digistruct Peak) | Pearlescent world drop (very rare/legendary pools) |
| World Burn | Sunburner | Rocket Launcher | Effervescent | Commander Lilith | Colonel Zarpedon |
| BlASSter | Scav Blasstah! | Assault Rifle | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Blaster | Proto-Blitz | Assault Rifle | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Dart | Needle Driver | Pistol | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Spiker | Spike Array | Pistol | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Plasma Caster | Plasma Flicker | SMG | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Splasher Blashter | Splashcaster | Shotgun | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Splatgun | Splat-O-Matic | Shotgun | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Hybridfication | Hybridizer | Sniper Rifle | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Moloko | Moonmilk | Sniper Rifle | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Railer | Rail Driver | Sniper Rifle | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Launcher | Orbit Lobber | Rocket Launcher | E-tech | — | E-tech world drop (very rare/legendary pools) |
| PBFG | Big Freakin' Prototype | Rocket Launcher | E-tech | — | E-tech world drop (very rare/legendary pools) |
| PRAZMA CANON | PRIZMA CANNON | Rocket Launcher | E-tech | — | E-tech world drop (very rare/legendary pools) |
| Topneaa | Torpnado | Rocket Launcher | E-tech | — | E-tech world drop (very rare/legendary pools) |
| 1340 Shield | Aegis-1340 | Shield | Rare | — | Bruce |
| Aequitas | Equalizer Field | Shield | Rare | — | Red Shirt |
| Captain Blade's Manly Man Shield | Captain Scav's Macho Mantle | Shield | Rare | Captain Scarlett | Bosun |
| Order | Lawman's Oath | Shield | Rare | — | Deadlift |
| Pot O' Gold | Prospector's Prize | Shield | Rare | — | Nel |
| The Rough Rider | No Safety Net | Shield | Rare | Sir Hammerlock | Meg |
| Cracked Sash | Fractured Halo | Shield | Epic | — | Tork Dredger |
| Antagonist | The Provocateur | Shield | Seraph | Tiny Tina (Dragon Keep) | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Big Boom Blaster | Blast Dividend | Shield | Seraph | Mr. Torgue | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Blockade | Hard Vacuum | Shield | Seraph | Tiny Tina (Dragon Keep) | Invincible Empyrean Sentinel (raid) |
| Evolution | Adaptation Protocol | Shield | Seraph | Captain Scarlett | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Hoplite | Phalanx Zero | Shield | Seraph | Mr. Torgue | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Pun-chee | Punchline | Shield | Seraph | Mr. Torgue | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Sponge | Absorb-O-Matic | Shield | Seraph | Mr. Torgue | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Flame of the Firehawk | Helios Halo | Shield | Legendary | — | Flame Knuckle |
| Hide of Terramorphous | Iwajira's Hide | Shield | Legendary | — | Iwajira |
| Impaler | Meteor Spike | Shield | Legendary | — | Pool_Shields_SpikeShields_Corrosive; World drop (legendary pool) |
| The Bee | Moonquito | Shield | Legendary | — | World drop (legendary pool) |
| Easy Mode | Training Wheels | Shield | Effervescent | Commander Lilith | World drop (legendary pool) |
| Retainer | Spikekeeper | Shield | Effervescent | Commander Lilith | Giant Shuggurath of Ice |
| Captain Blade's Midnight Star | Captain Blackstar's Nightcap | Grenade Mod | Rare | Captain Scarlett | Bosun |
| Fireball | Spellburner | Grenade Mod | Rare | Tiny Tina (Dragon Keep) | Flame Knuckle |
| Fuster Cluck | Cluster Duck | Grenade Mod | Rare | — | Rabid Adams |
| Lightning Bolt | Arc Bolt | Grenade Mod | Rare | Tiny Tina (Dragon Keep) | Badass Kraggon |
| Magic Missile | Hex-Seeker | Grenade Mod | Rare | Tiny Tina (Dragon Keep) | Eghood |
| Breath of Terramorphous | Iwajira's Breath | Grenade Mod | E-tech | — | Iwajira |
| Crossfire | Crosscurrent | Grenade Mod | Seraph | Mr. Torgue | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Meteor Shower | Asteroid Rain | Grenade Mod | Seraph | Mr. Torgue | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| O-Negative | Universal Donor | Grenade Mod | Seraph | Mr. Torgue | Invincible Empyrean Sentinel (raid); Moonstone loot (Grinder/chests) |
| Chain Lightning | Daisy Chain | Grenade Mod | Legendary | Tiny Tina (Dragon Keep) | Rooster Booster |
| Fastball | Orbital Fastpitch | Grenade Mod | Legendary | — | World drop (legendary pool) |
| Fire Storm | Solar Flare | Grenade Mod | Legendary | Tiny Tina (Dragon Keep) | Flame Knuckle |
| Antifection | Antivirus Bomb | Grenade Mod | Effervescent | Commander Lilith | Eghood |
| The Electric Chair | Shock Therapy | Grenade Mod | Effervescent | Commander Lilith | Boomer |
