package backend

type (
	Config struct {
		Colors Colors
		Labels Labels
	}
)

func defaultConfig() Config {
	return Config{
		Colors: Colors{
			ColorID("blue"):    HexCode("#6C93C7"),
			ColorID("green"):   HexCode("#7FB380"),
			ColorID("red"):     HexCode("#C77373"),
			ColorID("magenta"): HexCode("#B080B0"),
			ColorID("orange"):  HexCode("#D8A15C"),
			ColorID("teal"):    HexCode("#6BB3A3"),

			ColorID("valid"):   HexCode("#C0C0C0"),
			ColorID("invalid"): HexCode("#808080"),
		},

		Labels: Labels{
			LabelID("task"): LabelInfo{
				Name:    "Tasks",
				ColorID: ColorID("red"),
			},
			"university": LabelInfo{
				Name:    "University",
				ColorID: ColorID("blue"),
			},
		},
	}
}

func MockConfig() (Labels, Colors) {
	c := defaultConfig()
	return c.Labels, c.Colors
}
