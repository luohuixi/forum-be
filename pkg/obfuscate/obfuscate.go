package obfuscate

import (
	"errors"

	hashids "github.com/speps/go-hashids/v2"
)

type Obfuscator struct {
	*hashids.HashID
}

func NewObfuscator(salt string, minLength int) *Obfuscator {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	return &Obfuscator{
		h,
	}
}

func (o *Obfuscator) Obfuscate(id uint) string {
	hid, err := o.Encode([]int{int(id)})
	if err != nil {
		return ""
	}
	return hid
}

func (o *Obfuscator) Deobfuscate(hid string) (uint, error) {
	if hid == "" {
		return 0, errors.New("hid is empty")
	}

	ids, err := o.DecodeWithError(hid)
	if err != nil {
		return 0, err
	}
	return uint(ids[0]), nil
}

func (o *Obfuscator) DeobfuscateHids(hids []string) ([]uint, error) {
	ids := make([]uint, len(hids))
	for key, hid := range hids {
		if hid == "" {
			return nil, errors.New("hid is an empty string")
		}
		id, err := o.Deobfuscate(hid)
		if err != nil {
			return nil, err
		}
		ids[key] = id
	}
	return ids, nil
}
