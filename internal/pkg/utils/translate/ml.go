package translate

func ClothesTypeToTryOnCategory(clothesType string) string {
	switch clothesType {
	case "Верх":
		return "upper_body"
	case "Низ":
		return "lower_body"
	case "Платья":
		return "dresses"
	default:
		return ""
	}
}
