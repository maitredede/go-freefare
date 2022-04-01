package freefare

import (
	"fmt"

	nfc "github.com/clausecker/nfc/v2"
)

func GetTags(device nfc.Device) ([]SimpleTag, error) {
	if err := device.InitiatorInit(); err != nil {
		return nil, fmt.Errorf("intiator init failed: %w", err)
	}
	// Drop the field for a while
	if err := device.SetPropertyBool(nfc.ActivateField, false); err != nil {
		return nil, fmt.Errorf("field drop failed: %w", err)
	}

	// Configure the CRC and Parity settings
	if err := device.SetPropertyBool(nfc.HandleCRC, true); err != nil {
		return nil, fmt.Errorf("set handle crc failed: %w", err)
	}
	if err := device.SetPropertyBool(nfc.HandleParity, true); err != nil {
		return nil, fmt.Errorf("set handle crc failed: %w", err)
	}
	if err := device.SetPropertyBool(nfc.AutoISO14443_4, true); err != nil {
		return nil, fmt.Errorf("set handle crc failed: %w", err)
	}

	// Enable field so more power consuming cards can power themselves up
	if err := device.SetPropertyBool(nfc.ActivateField, true); err != nil {
		return nil, fmt.Errorf("field activate failed: %w", err)
	}
	mods := []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
		{Type: nfc.Felica, BaudRate: nfc.Nbr424},
	}
	tags := make([]SimpleTag, 0)
	for _, m := range mods {
		targets, err := device.InitiatorListPassiveTargets(m)
		if err != nil {
			continue
		}
		for _, t := range targets {
			if FelicaTaste(device, t) {
				t := SimpleTag{
					uid:        "",
					target:     t,
					modulation: m,
					tagType:    TagTypeFelica,
				}
				tags = append(tags, t)
				continue
			}
			if nti, ok := MifareMiniTaste(device, t); ok {
				t := SimpleTag{
					uid:        nti.String(),
					target:     t,
					modulation: m,
					tagType:    TagTypeMini,
				}
				tags = append(tags, t)
				continue

			}
			if nti, ok := MifareClassic1kTaste(device, t); ok {
				t := SimpleTag{
					uid:        nti.String(),
					target:     t,
					modulation: m,
					tagType:    TagTypeClassic1k,
				}
				tags = append(tags, t)
				continue
			}
			if nti, ok := MifareClassic4kTaste(device, t); ok {
				t := SimpleTag{
					uid:        nti.String(),
					target:     t,
					modulation: m,
					tagType:    TagTypeClassic4k,
				}
				tags = append(tags, t)
				continue
			}
			if nti, ok := MifareDesfireTaste(device, t); ok {
				t := SimpleTag{
					uid:        nti.String(),
					target:     t,
					modulation: m,
					tagType:    TagTypeDESFire,
				}
				tags = append(tags, t)
				continue
			}
			if nti, ok := NTag21xTaste(device, t); ok {
				t := SimpleTag{
					uid:        nti.String(),
					target:     t,
					modulation: m,
					tagType:    TagTypeNtag21x,
				}
				tags = append(tags, t)
				continue
			}
			if nti, ok := MifareUltralightcTaste(device, t); ok {
				t := SimpleTag{
					uid:        nti.String(),
					target:     t,
					modulation: m,
					tagType:    TagTypeUltralightC,
				}
				tags = append(tags, t)
				continue
			}
			if nti, ok := MifareUltralightTaste(device, t); ok {
				t := SimpleTag{
					uid:        nti.String(),
					target:     t,
					modulation: m,
					tagType:    TagTypeUltralight,
				}
				tags = append(tags, t)
				continue
			}
		}
	}
	return tags, nil
}

type SimpleTag struct {
	uid        string
	modulation nfc.Modulation
	target     nfc.Target
	tagType    TagType
}

func (t *SimpleTag) UID() string {
	return t.uid
}

func (t *SimpleTag) String() string {
	return t.uid
}

func (t *SimpleTag) Type() TagType {
	return t.tagType
}

func FelicaTaste(device nfc.Device, target nfc.Target) bool {
	return target.Modulation().Type == nfc.Felica
}

func MifareMiniTaste(device nfc.Device, target nfc.Target) (*nfc.ISO14443aTarget, bool) {
	if target.Modulation().Type != nfc.ISO14443a {
		return nil, false
	}

	nti, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		panic("cast failed")
	}
	return nti, nti.Sak == 0x09
}

func MifareClassic1kTaste(device nfc.Device, target nfc.Target) (*nfc.ISO14443aTarget, bool) {
	if target.Modulation().Type != nfc.ISO14443a {
		return nil, false
	}

	nti, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		panic("cast failed")
	}
	return nti, nti.Sak == 0x08 || nti.Sak == 0x28 || nti.Sak == 0x68 || nti.Sak == 0x88
}

func MifareClassic4kTaste(device nfc.Device, target nfc.Target) (*nfc.ISO14443aTarget, bool) {
	if target.Modulation().Type != nfc.ISO14443a {
		return nil, false
	}

	nti, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		panic("cast failed")
	}
	return nti, nti.Sak == 0x18 || nti.Sak == 0x38
}

func MifareDesfireTaste(device nfc.Device, target nfc.Target) (*nfc.ISO14443aTarget, bool) {
	if target.Modulation().Type != nfc.ISO14443a {
		return nil, false
	}

	nti, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		panic("cast failed")
	}
	if nti.Sak != 0x20 {
		return nti, false
	}
	panic("todo")
}

func NTag21xTaste(device nfc.Device, target nfc.Target) (*nfc.ISO14443aTarget, bool) {
	if target.Modulation().Type != nfc.ISO14443a {
		return nil, false
	}

	nti, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		panic("cast failed")
	}
	if nti.Sak != 0x00 {
		return nti, false
	}
	isAuthSupported := ntag21xIsAuthSupported(device, nti)
	return nti, isAuthSupported
}

func ntag21xIsAuthSupported(device nfc.Device, nai *nfc.ISO14443aTarget) bool {
	modulation := nfc.Modulation{
		Type:     nfc.ISO14443a,
		BaudRate: nfc.Nbr106,
	}

	_ /*pnti*/, err := device.InitiatorSelectPassiveTarget(modulation, nil)
	if err != nil {
		panic(err)
	}
	if err := device.SetPropertyBool(nfc.EasyFraming, false); err != nil {
		panic(err)
	}
	var cmd_step1 []byte = []byte{0x60}
	var res_step1 []byte = make([]byte, 8)
	_, err = device.InitiatorTransceiveBytes(cmd_step1, res_step1, 0)
	if err := device.SetPropertyBool(nfc.EasyFraming, true); err != nil {
		panic(err)
	}
	if err := device.InitiatorDeselectTarget(); err != nil {
		panic(err)
	}
	return err != nil
}

func MifareUltralightcTaste(device nfc.Device, target nfc.Target) (*nfc.ISO14443aTarget, bool) {
	if target.Modulation().Type != nfc.ISO14443a {
		return nil, false
	}

	nti, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		panic("cast failed")
	}
	if nti.Sak != 0x00 {
		return nti, false
	}
	isMUcOnReader := isMifateUltralightCOnReader(device, nti)
	return nti, isMUcOnReader
}

func MifareUltralightTaste(device nfc.Device, target nfc.Target) (*nfc.ISO14443aTarget, bool) {
	if target.Modulation().Type != nfc.ISO14443a {
		return nil, false
	}

	nti, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		panic("cast failed")
	}
	if nti.Sak != 0x00 {
		return nti, false
	}
	isMUcOnReader := isMifateUltralightCOnReader(device, nti)
	return nti, !isMUcOnReader
}

func isMifateUltralightCOnReader(device nfc.Device, nai *nfc.ISO14443aTarget) bool {
	modulation := nfc.Modulation{
		Type:     nfc.ISO14443a,
		BaudRate: nfc.Nbr106,
	}

	_ /*pnti*/, err := device.InitiatorSelectPassiveTarget(modulation, nai.UID[:])
	if err != nil {
		panic(err)
	}
	if err := device.SetPropertyBool(nfc.EasyFraming, false); err != nil {
		panic(err)
	}
	var cmd_step1 []byte = []byte{0x1A, 0x00}
	var res_step1 []byte = make([]byte, 9)
	_, err = device.InitiatorTransceiveBytes(cmd_step1, res_step1, 0)
	if err := device.SetPropertyBool(nfc.EasyFraming, true); err != nil {
		panic(err)
	}
	if err := device.InitiatorDeselectTarget(); err != nil {
		panic(err)
	}
	return err != nil
}
