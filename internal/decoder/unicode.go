package decoder

import "strings"

var letters = map[[2]byte]string{
	[2]byte{208, 176}: "а",
	[2]byte{208, 177}: "б",
	[2]byte{208, 178}: "в",
	[2]byte{208, 179}: "г",
	[2]byte{208, 180}: "д",
	[2]byte{208, 181}: "е",
	[2]byte{208, 182}: "ж",
	[2]byte{208, 183}: "з",
	[2]byte{208, 184}: "и",
	[2]byte{208, 185}: "й",
	[2]byte{208, 186}: "к",
	[2]byte{208, 187}: "л",
	[2]byte{208, 188}: "м",
	[2]byte{208, 189}: "н",
	[2]byte{208, 190}: "о",
	[2]byte{208, 191}: "п",
	[2]byte{209, 128}: "р",
	[2]byte{209, 129}: "с",
	[2]byte{209, 130}: "т",
	[2]byte{209, 131}: "у",
	[2]byte{209, 132}: "ф",
	[2]byte{209, 133}: "х",
	[2]byte{209, 134}: "ц",
	[2]byte{209, 135}: "ч",
	[2]byte{209, 136}: "ш",
	[2]byte{209, 137}: "ъ",
	[2]byte{209, 138}: "ы",
	[2]byte{209, 139}: "ь",
	[2]byte{209, 140}: "ь",
	[2]byte{209, 141}: "э",
	[2]byte{209, 142}: "ю",
	[2]byte{209, 143}: "я",

	[2]byte{32, 0}: " ",

	[2]byte{208, 144}: "А",
	[2]byte{208, 145}: "Б",
	[2]byte{208, 146}: "В",
	[2]byte{208, 147}: "Г",
	[2]byte{208, 148}: "Д",
	[2]byte{208, 149}: "Е",
	[2]byte{208, 150}: "Ж",
	[2]byte{208, 151}: "З",
	[2]byte{208, 152}: "И",
	[2]byte{208, 153}: "Й",
	[2]byte{208, 154}: "К",
	[2]byte{208, 155}: "Л",
	[2]byte{208, 156}: "М",
	[2]byte{208, 157}: "Н",
	[2]byte{208, 158}: "О",
	[2]byte{208, 159}: "П",
	[2]byte{208, 160}: "Р",
	[2]byte{208, 161}: "С",
	[2]byte{208, 162}: "Т",
	[2]byte{208, 163}: "У",
	[2]byte{208, 164}: "Ф",
	[2]byte{208, 165}: "Х",
	[2]byte{208, 166}: "Ц",
	[2]byte{208, 167}: "Ч",
	[2]byte{208, 168}: "Ш",
	[2]byte{208, 169}: "Щ",
	[2]byte{208, 170}: "Ъ",
	[2]byte{208, 171}: "Ы",
	[2]byte{208, 172}: "Ь",
	[2]byte{208, 173}: "Э",
	[2]byte{208, 174}: "Ю",
	[2]byte{208, 175}: "Я",
}

func BytesToUnicodeString(b []byte) string {
	iter := 1
	var result strings.Builder
	var buf [2]byte
	for iter < len(b) {
		if b[iter] >= 208 {
			copy(buf[:], b[iter:iter+2])
			result.WriteString(DecodeUnicodeSymbol(buf))
			iter += 2
		} else {
			result.WriteString(string(b[iter]))
			iter += 1
		}
	}
	return result.String()
}

func DecodeUnicodeSymbol(b [2]byte) string {
	res, set := letters[b]
	if !set {
		return "!UNKNOWN!"
	}
	return res
}
