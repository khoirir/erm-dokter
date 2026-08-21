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

func (k Kesadaran) IsValid() bool {
	switch k {
	case KesadaranComposMentis, KesadaranSomnolen, KesadaranSopor, KesadaranKoma,
		KesadaranAlert, KesadaranConfusion, KesadaranVoice, KesadaranPain,
		KesadaranUnresponsive, KesadaranApatis, KesadaranDelirium:
		return true
	default:
		return false
	}
}
