package backend

// type (
// 	Config struct {
// 		Colors Colors
// 		Labels Labels
// 	}
// )
//
// func defaultConfig() Config {
// 	return Config{
// 		Colors: Colors{
// 			ColorID("blue"):    HexCode("#6C93C7"),
// 			ColorID("green"):   HexCode("#7FB380"),
// 			ColorID("red"):     HexCode("#C77373"),
// 			ColorID("magenta"): HexCode("#B080B0"),
// 			ColorID("orange"):  HexCode("#D8A15C"),
// 			ColorID("teal"):    HexCode("#6BB3A3"),
//
// 			ColorID("valid"):   HexCode("#C0C0C0"),
// 			ColorID("invalid"): HexCode("#808080"),
// 		},
//
// 		Labels: Labels{
// 			Selected: "uni",
// 			Info: map[LabelID]LabelInfo{
// 				"uni":  {Name: "University", Color: HexCode("#81c8be")},
// 				"work": {Name: "Work", Color: HexCode("#ca9ee6")},
// 				"gym":  {Name: "Gym", Color: HexCode("#ef9f76")},
// 				"app":  {Name: "Application", Color: HexCode("#babbf1")},
// 				"vol":  {Name: "Volunteer", Color: HexCode("#f2d5cf")},
// 			},
// 			Order: []LabelID{
// 				"uni", "work", "gym", "app", "vol",
// 			},
// 		},
// 	}
// }
//
// func MockConfig() (Labels, Colors) {
// 	c := defaultConfig()
// 	return c.Labels, c.Colors
// }
