package main

type NationalPark struct {
	UnitName string `json:"name"`
	UnitCode string `json:"code"`
	GeomData string `json:"geom"`
}

func AllNationalParks() (parks []NationalPark) {
	db.Model(NationalPark{}).Select("ST_AsGeoJSON(ST_CollectionExtract(geom, 3)) as geom_data, unit_name, unit_code").Scan(&parks)

	return
}
