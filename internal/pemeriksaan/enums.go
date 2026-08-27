package pemeriksaan

type Kesadaran string

const (
	KesadaranComposMentis Kesadaran = "Compos Mentis"
	KesadaranSomnolen     Kesadaran = "Somnolen"
	KesadaranSopor        Kesadaran = "Sopor"
	KesadaranKoma         Kesadaran = "Koma"
	KesadaranAlert        Kesadaran = "Alert"
	KesadaranConfusion    Kesadaran = "Confusion"
	KesadaranVoice        Kesadaran = "Voice"
	KesadaranPain         Kesadaran = "Pain"
	KesadaranUnresponsive Kesadaran = "Unresponsive"
	KesadaranApatis       Kesadaran = "Apatis"
	KesadaranDelirium     Kesadaran = "Delirium"
)

var ListKesadaran = []Kesadaran{
	KesadaranComposMentis,
	KesadaranSomnolen,
	KesadaranSopor,
	KesadaranKoma,
	KesadaranAlert,
	KesadaranConfusion,
	KesadaranVoice,
	KesadaranPain,
	KesadaranUnresponsive,
	KesadaranApatis,
	KesadaranDelirium,
}

func (k Kesadaran) IsValid() bool {
	for _, item := range ListKesadaran {
		if k == item {
			return true
		}
	}
	return false
}
