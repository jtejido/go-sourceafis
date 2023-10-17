package wsq

// import (
// 	"image"
// 	"io"
// )

// func Encode(w io.Writer, m image.Image) error {
// 	return encode(w)
// }

// func encode(w io.Writer, m image.Image) error {
// 	b := m.Bounds()
// 	var e encoder
// 	if ww, ok := w.(writer); ok {
// 		e.w = ww
// 	} else {
// 		e.w = bufio.NewWriter(w)
// 	}
// 	var fdata []float32	/* floating point pixel image  */
// 	var data []int	/* quantized image pointer     */
// 	/* quantized block sizes */
// 	qsize := reference[int]{0}
// 	qsize1 := reference[int]{0}
// 	qsize2 := reference[int]{0}
// 	qsize3 := reference[int]{0}
// 	hufftable := &huffmanTable{}
// 	huffbits := make([]reference[int],0)
// 	huffvalues := make([]reference[int],0) /* huffman code parameters */

// 	/* Convert image pixels to floating point. */
// 	m_shift := reference[float32]
// 	 r_scale := reference[float32]
// 	fdata = convertImageToFloat(w, b.Width, b.Height, m_shift, r_scale);

// 	 token = new Token();

// 	/* Build WSQ decomposition trees */
// 	WSQHelper.buildWSQTrees(token, bitmap.getWidth(), bitmap.getHeight());

// 	/* WSQ decompose the image */
// 	wsqDecompose(token, fdata, bitmap.getWidth(), bitmap.getHeight(), token.tableDTT.hifilt, MAX_HIFILT, token.tableDTT.lofilt, MAX_LOFILT);

// 	/* Set compression ratio and 'q' to zero. */
// 	token.quant_vals.cr = 0;
// 	token.quant_vals.q = 0.0f;

// 	/* Assign specified r-bitrate into quantization structure. */
// 	token.quant_vals.r = (float)bitRate;

// 	/* Compute subband variances. */
// 	variance(token, fdata, bitmap.getWidth(), bitmap.getHeight());

// 	/* Quantize the floating point pixmap. */

// 	qdata = quantize(token, qsize, fdata, bitmap.getWidth(), bitmap.getHeight());

// 	/* Compute quantized WSQ subband block sizes */
// 	quant_block_sizes(token, qsize1, qsize2, qsize3);

// 	if (qsize.value != qsize1.value + qsize2.value + qsize3.value) {
// 		throw new IllegalStateException("ERROR : wsq_encode_1 : problem w/quantization block sizes");
// 	}

// 	/* Add a Start Of Image (SOI) marker to the WSQ buffer. */
// 	dataOutput.writeShort(SOI_WSQ);

// 	putc_nistcom_wsq(dataOutput, bitmap, (float)bitRate, metadata, comments);

// 	/* Store the Wavelet filter taps to the WSQ buffer. */
// 	putc_transform_table(dataOutput, token.tableDTT.lofilt, MAX_LOFILT, token.tableDTT.hifilt, MAX_HIFILT);

// 	/* Store the quantization parameters to the WSQ buffer. */
// 	putc_quantization_table(dataOutput, token);

// 	/* Store a frame header to the WSQ buffer. */
// 	putc_frame_header_wsq(dataOutput, bitmap.getWidth(), bitmap.getHeight(), m_shift.value, r_scale.value);

// 	/* ENCODE Block 1 */

// 	/* Compute Huffman table for Block 1. */
// 	hufftable = gen_hufftable_wsq(token, huffbits, huffvalues, qdata, 0, new int[] { qsize1.value });

// 	/* Store Huffman table for Block 1 to WSQ buffer. */
// 	putc_huffman_table(dataOutput, DHT_WSQ, 0, huffbits.value, huffvalues.value);

// 	/* Store Block 1's header to WSQ buffer. */
// 	putc_block_header(dataOutput, 0);

// 	/* Compress Block 1 data. */
// 	compress_block(dataOutput, qdata, 0, qsize1.value, MAX_HUFFCOEFF, MAX_HUFFZRUN, hufftable);

// 	/* ENCODE Block 2 */

// 	/* Compute  Huffman table for Blocks 2 & 3. */
// 	hufftable = gen_hufftable_wsq(token, huffbits, huffvalues, qdata, qsize1.value, new int[] { qsize2.value, qsize3.value });

// 	/* Store Huffman table for Blocks 2 & 3 to WSQ buffer. */
// 	putc_huffman_table(dataOutput, DHT_WSQ, 1, huffbits.value, huffvalues.value);

// 	/* Store Block 2's header to WSQ buffer. */
// 	putc_block_header(dataOutput, 1);

// 	/* Compress Block 2 data. */
// 	compress_block(dataOutput, qdata, qsize1.value, qsize2.value, MAX_HUFFCOEFF, MAX_HUFFZRUN, hufftable);

// 	/* ENCODE Block 3 */

// 	/* Store Block 3's header to WSQ buffer. */
// 	putc_block_header(dataOutput, 1);

// 	/* Compress Block 3 data. */
// 	compress_block(dataOutput, qdata, qsize1.value + qsize2.value, qsize3.value, MAX_HUFFCOEFF, MAX_HUFFZRUN, hufftable);

// 	/* Add a End Of Image (EOI) marker to the WSQ buffer. */
// 	dataOutput.writeShort(EOI_WSQ);
// }
