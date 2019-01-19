package analyzer

import (
	"fmt"
	"io"
	"os"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
)

func findPlayerID(replay *rep.Replay, names map[string]struct{}) byte {
	for _, p := range replay.Header.PIDPlayers {
		if _, ok := names[p.Name]; ok {
			return p.ID
		}
	}
	return 127 // On a byte field and for a player id, this will be a poor man's None
}

// If the command is a specific building/unit creation/evolution of a specific player id, it returns the second
// that it happened. Returns true if it was.
// Should be used on ProcessCommand.
func maybePlayersUnitSeconds(command repcmd.Cmd, playerID byte, unitID uint16) (string, bool) {
	if command.BaseCmd().PlayerID == playerID {
		switch c := command.(type) {
		case *repcmd.BuildCmd: // N.B due to a limitation in Go's type system this code cannot be simpler -_-
			if c.Unit.ID == unitID {
				return fmt.Sprintf("%v", int(c.Frame.Seconds())), true
			}
		case *repcmd.BuildingMorphCmd:
			if c.Unit.ID == unitID {
				return fmt.Sprintf("%v", int(c.Frame.Seconds())), true
			}
		case *repcmd.TrainCmd:
			if c.Unit.ID == unitID {
				return fmt.Sprintf("%v", int(c.Frame.Seconds())), true
			}
		}
	}
	return "-1", false
}

var (
	nameToUnitID = map[string]uint16{
		"Marine":                        0x00,
		"Ghost":                         0x01,
		"Vulture":                       0x02,
		"Goliath":                       0x03,
		"Goliath Turret":                0x04,
		"Siege Tank (Tank Mode)":        0x05,
		"Siege Tank Turret (Tank Mode)": 0x06,
		"SCV":                                 0x07,
		"Wraith":                              0x08,
		"Science Vessel":                      0x09,
		"Gui Motang (Firebat)":                0x0A,
		"Dropship":                            0x0B,
		"Battlecruiser":                       0x0C,
		"Spider Mine":                         0x0D,
		"Nuclear Missile":                     0x0E,
		"Terran Civilian":                     0x0F,
		"Sarah Kerrigan (Ghost)":              0x10,
		"Alan Schezar (Goliath)":              0x11,
		"Alan Schezar Turret":                 0x12,
		"Jim Raynor (Vulture)":                0x13,
		"Jim Raynor (Marine)":                 0x14,
		"Tom Kazansky (Wraith)":               0x15,
		"Magellan (Science Vessel)":           0x16,
		"Edmund Duke (Tank Mode)":             0x17,
		"Edmund Duke Turret (Tank Mode)":      0x18,
		"Edmund Duke (Siege Mode)":            0x19,
		"Edmund Duke Turret (Siege Mode)":     0x1A,
		"Arcturus Mengsk (Battlecruiser)":     0x1B,
		"Hyperion (Battlecruiser)":            0x1C,
		"Norad II (Battlecruiser)":            0x1D,
		"Terran Siege Tank (Siege Mode)":      0x1E,
		"Siege Tank Turret (Siege Mode)":      0x1F,
		"Firebat":                             0x20,
		"Scanner Sweep":                       0x21,
		"Medic":                               0x22,
		"Larva":                               0x23,
		"Egg":                                 0x24,
		"Zergling":                            0x25,
		"Hydralisk":                           0x26,
		"Ultralisk":                           0x27,
		"Drone":                               0x29,
		"Overlord":                            0x2A,
		"Mutalisk":                            0x2B,
		"Guardian":                            0x2C,
		"Queen":                               0x2D,
		"Defiler":                             0x2E,
		"Scourge":                             0x2F,
		"Torrasque (Ultralisk)":               0x30,
		"Matriarch (Queen)":                   0x31,
		"Infested Terran":                     0x32,
		"Infested Kerrigan (Infested Terran)": 0x33,
		"Unclean One (Defiler)":               0x34,
		"Hunter Killer (Hydralisk)":           0x35,
		"Devouring One (Zergling)":            0x36,
		"Kukulza (Mutalisk)":                  0x37,
		"Kukulza (Guardian)":                  0x38,
		"Yggdrasill (Overlord)":               0x39,
		"Valkyrie":                            0x3A,
		"Mutalisk Cocoon":                     0x3B,
		"Corsair":                             0x3C,
		"Dark Templar":                        0x3D,
		"Devourer":                            0x3E,
		"Dark Archon":                         0x3F,
		"Probe":                               0x40,
		"Zealot":                              0x41,
		"Dragoon":                             0x42,
		"High Templar":                        0x43,
		"Archon":                              0x44,
		"Shuttle":                             0x45,
		"Scout":                               0x46,
		"Arbiter":                             0x47,
		"Carrier":                             0x48,
		"Interceptor":                         0x49,
		"Protoss Dark Templar (Hero)":         0x4A,
		"Zeratul (Dark Templar)":              0x4B,
		"Tassadar/Zeratul (Archon)":           0x4C,
		"Fenix (Zealot)":                      0x4D,
		"Fenix (Dragoon)":                     0x4E,
		"Tassadar (Templar)":                  0x4F,
		"Mojo (Scout)":                        0x50,
		"Warbringer (Reaver)":                 0x51,
		"Gantrithor (Carrier)":                0x52,
		"Reaver":                              0x53,
		"Observer":                            0x54,
		"Scarab":                              0x55,
		"Danimoth (Arbiter)":                  0x56,
		"Aldaris (Templar)":                   0x57,
		"Artanis (Scout)":                     0x58,
		"Rhynadon (Badlands Critter)":         0x59,
		"Bengalaas (Jungle Critter)":          0x5A,
		"Cargo Ship (Unused)":                 0x5B,
		"Mercenary Gunship (Unused)":          0x5C,
		"Scantid (Desert Critter)":            0x5D,
		"Kakaru (Twilight Critter)":           0x5E,
		"Ragnasaur (Ashworld Critter)":        0x5F,
		"Ursadon (Ice World Critter)":         0x60,
		"Lurker Egg":                          0x61,
		"Raszagal (Corsair)":                  0x62,
		"Samir Duran (Ghost)":                 0x63,
		"Alexei Stukov (Ghost)":               0x64,
		"Map Revealer":                        0x65,
		"Gerard DuGalle (BattleCruiser)":      0x66,
		"Lurker": 0x67,
		"Infested Duran (Infested Terran)":     0x68,
		"Disruption Web":                       0x69,
		"Command Center":                       0x6A,
		"ComSat":                               0x6B,
		"Nuclear Silo":                         0x6C,
		"Supply Depot":                         0x6D,
		"Refinery":                             0x6E,
		"Barracks":                             0x6F,
		"Academy":                              0x70,
		"Factory":                              0x71,
		"Starport":                             0x72,
		"Control Tower":                        0x73,
		"Science Facility":                     0x74,
		"Covert Ops":                           0x75,
		"Physics Lab":                          0x76,
		"Machine Shop":                         0x78,
		"Repair Bay (Unused)":                  0x79,
		"Engineering Bay":                      0x7A,
		"Armory":                               0x7B,
		"Missile Turret":                       0x7C,
		"Bunker":                               0x7D,
		"Norad II (Crashed)":                   0x7E,
		"Ion Cannon":                           0x7F,
		"Uraj Crystal":                         0x80,
		"Khalis Crystal":                       0x81,
		"Infested CC":                          0x82,
		"Hatchery":                             0x83,
		"Lair":                                 0x84,
		"Hive":                                 0x85,
		"Nydus Canal":                          0x86,
		"Hydralisk Den":                        0x87,
		"Defiler Mound":                        0x88,
		"Greater Spire":                        0x89,
		"Queens Nest":                          0x8A,
		"Evolution Chamber":                    0x8B,
		"Ultralisk Cavern":                     0x8C,
		"Spire":                                0x8D,
		"Spawning Pool":                        0x8E,
		"Creep Colony":                         0x8F,
		"Spore Colony":                         0x90,
		"Unused Zerg Building1":                0x91,
		"Sunken Colony":                        0x92,
		"Zerg Overmind (With Shell)":           0x93,
		"Overmind":                             0x94,
		"Extractor":                            0x95,
		"Mature Chrysalis":                     0x96,
		"Cerebrate":                            0x97,
		"Cerebrate Daggoth":                    0x98,
		"Unused Zerg Building2":                0x99,
		"Nexus":                                0x9A,
		"Robotics Facility":                    0x9B,
		"Pylon":                                0x9C,
		"Assimilator":                          0x9D,
		"Unused Protoss Building1":             0x9E,
		"Observatory":                          0x9F,
		"Gateway":                              0xA0,
		"Unused Protoss Building2":             0xA1,
		"Photon Cannon":                        0xA2,
		"Citadel of Adun":                      0xA3,
		"Cybernetics Core":                     0xA4,
		"Templar Archives":                     0xA5,
		"Forge":                                0xA6,
		"Stargate":                             0xA7,
		"Stasis Cell/Prison":                   0xA8,
		"Fleet Beacon":                         0xA9,
		"Arbiter Tribunal":                     0xAA,
		"Robotics Support Bay":                 0xAB,
		"Shield Battery":                       0xAC,
		"Khaydarin Crystal Formation":          0xAD,
		"Protoss Temple":                       0xAE,
		"Xel'Naga Temple":                      0xAF,
		"Mineral Field (Type 1)":               0xB0,
		"Mineral Field (Type 2)":               0xB1,
		"Mineral Field (Type 3)":               0xB2,
		"Cave (Unused)":                        0xB3,
		"Cave-in (Unused)":                     0xB4,
		"Cantina (Unused)":                     0xB5,
		"Mining Platform (Unused)":             0xB6,
		"Independent Command Center (Unused)":  0xB7,
		"Independent Starport (Unused)":        0xB8,
		"Independent Jump Gate (Unused)":       0xB9,
		"Ruins (Unused)":                       0xBA,
		"Khaydarin Crystal Formation (Unused)": 0xBB,
		"Vespene Geyser":                       0xBC,
		"Warp Gate":                            0xBD,
		"Psi Disrupter":                        0xBE,
		"Zerg Marker":                          0xBF,
		"Terran Marker":                        0xC0,
		"Protoss Marker":                       0xC1,
		"Zerg Beacon":                          0xC2,
		"Terran Beacon":                        0xC3,
		"Protoss Beacon":                       0xC4,
		"Zerg Flag Beacon":                     0xC5,
		"Terran Flag Beacon":                   0xC6,
		"Protoss Flag Beacon":                  0xC7,
		"Power Generator":                      0xC8,
		"Overmind Cocoon":                      0xC9,
		"Dark Swarm":                           0xCA,
		"Floor Missile Trap":                   0xCB,
		"Floor Hatch (Unused)":                 0xCC,
		"Left Upper Level Door":                0xCD,
		"Right Upper Level Door":               0xCE,
		"Left Pit Door":                        0xCF,
		"Right Pit Door":                       0xD0,
		"Floor Gun Trap":                       0xD1,
		"Left Wall Missile Trap":               0xD2,
		"Left Wall Flame Trap":                 0xD3,
		"Right Wall Missile Trap":              0xD4,
		"Right Wall Flame Trap":                0xD5,
		"Start Location":                       0xD6,
		"Flag":                                 0xD7,
		"Young Chrysalis":                      0xD8,
		"Psi Emitter":                          0xD9,
		"Data Disc":                            0xDA,
		"Khaydarin Crystal":                    0xDB,
		"Mineral Cluster Type 1":               0xDC,
		"Mineral Cluster Type 2":               0xDD,
		"Protoss Vespene Gas Orb Type 1":       0xDE,
		"Protoss Vespene Gas Orb Type 2":       0xDF,
		"Zerg Vespene Gas Sac Type 1":          0xE0,
		"Zerg Vespene Gas Sac Type 2":          0xE1,
		"Terran Vespene Gas Tank Type 1":       0xE2,
		"Terran Vespene Gas Tank Type 2":       0xE3,
	}
	raceNameTranslations = map[string]string{
		"zerg":    "Zerg",
		"z":       "Zerg",
		"protoss": "Protoss",
		"p":       "Protoss",
		"toss":    "Protoss",
		"terran":  "Terran",
		"t":       "Terran",
		"ran":     "Terran",
	}
	// trainedUnitIDs = map[uint16]struct{}{
	// 	0x00: struct{}{}, // Terran units
	// 	0x01: struct{}{},
	// 	0x02: struct{}{},
	// 	0x03: struct{}{},
	// 	0x05: struct{}{},
	// 	0x07: struct{}{},
	// 	0x08: struct{}{},
	// 	0x09: struct{}{},
	// 	0x0B: struct{}{},
	// 	0x0C: struct{}{},
	// 	0x20: struct{}{},
	// 	0x22: struct{}{},
	// 	0x3C: struct{}{}, // Protoss units
	// 	0x3D: struct{}{},
	// 	0x3F: struct{}{},
	// 	0x40: struct{}{},
	// 	0x41: struct{}{},
	// 	0x42: struct{}{},
	// 	0x43: struct{}{},
	// 	0x44: struct{}{},
	// 	0x45: struct{}{},
	// 	0x46: struct{}{},
	// 	0x47: struct{}{},
	// 	0x48: struct{}{},
	// 	0x53: struct{}{},
	// 	0x54: struct{}{},
	// }
	// morphedUnitIDs = map[uint16]struct{}{
	// 	0x25: struct{}{}, // Zerg units
	// 	0x26: struct{}{},
	// 	0x27: struct{}{},
	// 	0x29: struct{}{},
	// 	0x2A: struct{}{},
	// 	0x2B: struct{}{},
	// 	0x2C: struct{}{},
	// 	0x2D: struct{}{},
	// 	0x2E: struct{}{},
	// 	0x2F: struct{}{},
	// 	0x32: struct{}{}, // TODO: Infested Terran probably trained?
	// 	0x3A: struct{}{},
	// 	0x3B: struct{}{},
	// 	0x3E: struct{}{},
	// }
	// builtUnitIDs = map[uint16]struct{}{
	// 	0x6A: struct{}{}, // Terran buildings
	// 	0x6B: struct{}{},
	// 	0x6C: struct{}{},
	// 	0x6D: struct{}{},
	// 	0x6E: struct{}{},
	// 	0x6F: struct{}{},
	// 	0x70: struct{}{},
	// 	0x71: struct{}{},
	// 	0x72: struct{}{},
	// 	0x73: struct{}{},
	// 	0x74: struct{}{},
	// 	0x75: struct{}{},
	// 	0x76: struct{}{},
	// 	0x78: struct{}{},
	// 	0x7A: struct{}{},
	// 	0x7B: struct{}{},
	// 	0x7C: struct{}{},
	// 	0x7D: struct{}{},
	// 	0x9A: struct{}{}, // Protoss buildings
	// 	0x9B: struct{}{},
	// 	0x9C: struct{}{},
	// 	0x9D: struct{}{},
	// 	0x9F: struct{}{},
	// 	0xA0: struct{}{},
	// 	0xA2: struct{}{},
	// 	0xA3: struct{}{},
	// 	0xA4: struct{}{},
	// 	0xA5: struct{}{},
	// 	0xA6: struct{}{},
	// 	0xA7: struct{}{},
	// 	0xA9: struct{}{},
	// 	0xAA: struct{}{},
	// 	0xAB: struct{}{},
	// 	0xAC: struct{}{},
	// 	0x83: struct{}{}, // Zerg buildings
	// 	0x86: struct{}{},
	// 	0x87: struct{}{},
	// 	0x88: struct{}{},
	// 	0x8A: struct{}{},
	// 	0x8B: struct{}{},
	// 	0x8C: struct{}{},
	// 	0x8D: struct{}{},
	// 	0x8E: struct{}{},
	// 	0x8F: struct{}{},
	// 	0x95: struct{}{},
	// }
	// buldingMorphedUnits = map[uint16]struct{}{
	// 	0x84: struct{}{}, // Zerg building evolutions
	// 	0x85: struct{}{},
	// 	0x90: struct{}{},
	// 	0x92: struct{}{},
	// 	0x89: struct{}{},
	// }
)

// https://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-exists
// exists returns whether the given file or directory exists
func isFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func cloneStringSlice(ss []string) []string {
	var nss = make([]string, len(ss))
	for i := range ss {
		nss[i] = ss[i]
	}
	return nss
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

// https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
