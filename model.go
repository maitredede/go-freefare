package freefare

type Tag struct {
	TagValue string  `json:"tag"`
	TagType  TagType `json:"type"`
	Device   string  `json:"device"`
}

type TagType string

const (
	TagTypeFelica      TagType = "felica"
	TagTypeMini        TagType = "mini"
	TagTypeClassic1k   TagType = "classic1k"
	TagTypeClassic4k   TagType = "classic4k"
	TagTypeDESFire     TagType = "desfire"
	TagTypeUltralight  TagType = "ultralight"
	TagTypeUltralightC TagType = "ultralightc"
	TagTypeNtag21x     TagType = "ntagg21x"
)

// from https://github.com/clausecker/freefare
// var (
// 	tagTypeMap map[int]TagType = map[int]TagType{
// 		freefare.Felica:      TagTypeFelica,
// 		freefare.Mini:        TagTypeMini,
// 		freefare.Classic1k:   TagTypeClassic1k,
// 		freefare.Classic4k:   TagTypeClassic4k,
// 		freefare.DESFire:     TagTypeDESFire,
// 		freefare.Ultralight:  TagTypeUltralight,
// 		freefare.UltralightC: TagTypeUltralightC,
// 		freefare.Ntag21x:     TagTypeNtag21x,
// 	}
// )
