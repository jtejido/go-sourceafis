package utils

func ExtensionFromMime(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpeg"
	case "image/jp2":
		return ".jp2"
	case "image/png":
		return ".png"
	case "image/svg+xml":
		return ".svg"
	case "image/tiff":
		return ".tiff"
	case "image/bmp":
		return ".bmp"
	/*
	 * WSQ doesn't have a MIME type. We will invent one.
	 */
	case "image/x-wsq":
		return ".wsq"
	/*
	 * No extension for unknown formats. This is safe. We don't want null.
	 */
	default:
		return ""
	}
}
