package master

type ItemMaster struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type Penjamin struct {
	ItemMaster
}

type Depo struct {
	ItemMaster
}

type Poliklinik struct {
	ItemMaster
}

type Bangsal struct {
	ItemMaster
}

type KelasKamar struct {
	ItemMaster
}

