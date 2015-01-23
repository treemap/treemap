package main

type SeaRise struct {
	Erosion  string `json:"erosion"`
	SeaLevel string `json:"sea_level"`
	GeomData string `json:"geom"`
}

func AllSeaRise() (sea_rise []SeaRise) {
	for _, i := range []string{"pacific_coasts", "atlantic_coasts", "gulf_coasts"} {
		var s []SeaRise

		db.Table(i).Select("ST_AsGeoJSON(geog) as geom_data, erosion").Scan(&s)
		sea_rise = append(sea_rise, s...)
	}

	return
}
