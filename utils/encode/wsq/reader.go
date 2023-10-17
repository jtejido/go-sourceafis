package wsq

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
)

func Decode(r io.Reader) (image.Image, error) {
	return decode(r)
}

func decode(r io.Reader) (image.Image, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	token := newToken(buf.Bytes())

	/* Read the SOI marker. */
	_, err := token.getCMarkerWSQ(soi_wsq)
	if err != nil {
		return nil, err
	}

	/* Read in supporting tables up to the SOF marker. */
	marker, err := token.getCMarkerWSQ(tbls_n_sof)
	if err != nil {
		return nil, err
	}
	for marker != sof_wsq {
		if err := token.getCTableWSQ(marker); err != nil {
			return nil, err
		}
		marker, err = token.getCMarkerWSQ(tbls_n_sof)
		if err != nil {
			return nil, err
		}
	}

	/* Read in the Frame Header. */
	frmHeaderWSQ := token.getCFrameHeaderWSQ()
	width := frmHeaderWSQ.width
	height := frmHeaderWSQ.height

	getCPpiWSQ()

	/* Build WSQ decomposition trees. */
	buildWSQTrees(token, width, height)

	/* Decode the Huffman encoded buffer blocks. */

	qdata, err := huffmanDecodeDataMem(token, width*height)
	if err != nil {
		return nil, err
	}

	/* Decode the quantize wavelet subband buffer. */
	fdata, err := unquantize(token, qdata, width, height)
	if err != nil {
		return nil, err
	}

	/* Done with quantized wavelet subband buffer. */
	//noinspection UnusedAssignment
	qdata = nil

	wsqReconstruct(token, fdata, width, height)

	/* Convert floating point pixels to unsigned char pixels. */
	img := image.NewGray(image.Rect(0, 0, width, height))
	var idx int
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			pixel := (fdata[idx] * frmHeaderWSQ.rScale) + frmHeaderWSQ.mShift
			pixel += 0.5
			if pixel < 0.0 {
				img.SetGray(c, r, color.Gray{Y: 0})
			} else if pixel > 255.0 {
				img.SetGray(c, r, color.Gray{Y: 255})
			} else {
				img.SetGray(c, r, color.Gray{Y: uint8(pixel)})
			}

			idx++
		}
	}

	return img, nil
}

func getCPpiWSQ() int {
	return -1
}

func intSign(power int) int { /* "sign" power */
	num := -1 /* sign return value */

	if power == 0 {
		return 1
	}

	for cnt := 1; cnt < power; cnt++ {
		num *= -1
	}

	return num
}

func buildWSQTrees(token *token, width, height int) {
	/* Build a W-TREE structure for the image. */
	buildWTree(token, w_treelen, width, height)
	/* Build a Q-TREE structure for the image. */
	buildQTree(token, q_treelen)
}

func buildWTree(token *token, wtreelen, width, height int) {
	var lenx, lenx2, leny, leny2 int /* starting lengths of sections of
	the image being split into subbands */
	token.wtree = make([]wavletTree, wtreelen)
	for i := 0; i < wtreelen; i++ {
		token.wtree[i] = wavletTree{
			invrw: 0,
			invcl: 0,
		}
	}

	token.wtree[2].invrw = 1
	token.wtree[4].invrw = 1
	token.wtree[7].invrw = 1
	token.wtree[9].invrw = 1
	token.wtree[11].invrw = 1
	token.wtree[13].invrw = 1
	token.wtree[16].invrw = 1
	token.wtree[18].invrw = 1
	token.wtree[3].invcl = 1
	token.wtree[5].invcl = 1
	token.wtree[8].invcl = 1
	token.wtree[9].invcl = 1
	token.wtree[12].invcl = 1
	token.wtree[13].invcl = 1
	token.wtree[17].invcl = 1
	token.wtree[18].invcl = 1

	wtree4(token, 0, 1, width, height, 0, 0, 1)

	if (token.wtree[1].lenx % 2) == 0 {
		lenx = token.wtree[1].lenx / 2
		lenx2 = lenx
	} else {
		lenx = (token.wtree[1].lenx + 1) / 2
		lenx2 = lenx - 1
	}

	if (token.wtree[1].leny % 2) == 0 {
		leny = token.wtree[1].leny / 2
		leny2 = leny
	} else {
		leny = (token.wtree[1].leny + 1) / 2
		leny2 = leny - 1
	}

	wtree4(token, 4, 6, lenx2, leny, lenx, 0, 0)
	wtree4(token, 5, 10, lenx, leny2, 0, leny, 0)
	wtree4(token, 14, 15, lenx, leny, 0, 0, 0)

	token.wtree[19].x = 0
	token.wtree[19].y = 0
	if (token.wtree[15].lenx % 2) == 0 {
		token.wtree[19].lenx = token.wtree[15].lenx / 2
	} else {
		token.wtree[19].lenx = (token.wtree[15].lenx + 1) / 2
	}

	if (token.wtree[15].leny % 2) == 0 {
		token.wtree[19].leny = token.wtree[15].leny / 2
	} else {
		token.wtree[19].leny = (token.wtree[15].leny + 1) / 2
	}
}

func wtree4(token *token, start1, start2, lenx, leny, x, y, stop1 int) {
	var evenx, eveny int /* Check length of subband for even or odd */
	var p1, p2 int       /* w_tree locations for storing subband sizes and locations */

	p1 = start1
	p2 = start2

	evenx = lenx % 2
	eveny = leny % 2

	token.wtree[p1].x = x
	token.wtree[p1].y = y
	token.wtree[p1].lenx = lenx
	token.wtree[p1].leny = leny

	token.wtree[p2].x = x
	token.wtree[p2+2].x = x
	token.wtree[p2].y = y
	token.wtree[p2+1].y = y

	if evenx == 0 {
		token.wtree[p2].lenx = lenx / 2
		token.wtree[p2+1].lenx = token.wtree[p2].lenx
	} else {
		if p1 == 4 {
			token.wtree[p2].lenx = (lenx - 1) / 2
			token.wtree[p2+1].lenx = token.wtree[p2].lenx + 1
		} else {
			token.wtree[p2].lenx = (lenx + 1) / 2
			token.wtree[p2+1].lenx = token.wtree[p2].lenx - 1
		}
	}
	token.wtree[p2+1].x = token.wtree[p2].lenx + x
	if stop1 == 0 {
		token.wtree[p2+3].lenx = token.wtree[p2+1].lenx
		token.wtree[p2+3].x = token.wtree[p2+1].x
	}
	token.wtree[p2+2].lenx = token.wtree[p2].lenx

	if eveny == 0 {
		token.wtree[p2].leny = leny / 2
		token.wtree[p2+2].leny = token.wtree[p2].leny
	} else {
		if p1 == 5 {
			token.wtree[p2].leny = (leny - 1) / 2
			token.wtree[p2+2].leny = token.wtree[p2].leny + 1
		} else {
			token.wtree[p2].leny = (leny + 1) / 2
			token.wtree[p2+2].leny = token.wtree[p2].leny - 1
		}
	}
	token.wtree[p2+2].y = token.wtree[p2].leny + y
	if stop1 == 0 {
		token.wtree[p2+3].leny = token.wtree[p2+2].leny
		token.wtree[p2+3].y = token.wtree[p2+2].y
	}
	token.wtree[p2+1].leny = token.wtree[p2].leny
}

func buildQTree(token *token, qtreelen int) {
	token.qtree = make([]quantTree, qtreelen)
	for i := 0; i < len(token.qtree); i++ {
		token.qtree[i] = quantTree{}
	}

	qtree16(token, 3, token.wtree[14].lenx, token.wtree[14].leny, token.wtree[14].x, token.wtree[14].y, 0, 0)
	qtree16(token, 19, token.wtree[4].lenx, token.wtree[4].leny, token.wtree[4].x, token.wtree[4].y, 0, 1)
	qtree16(token, 48, token.wtree[0].lenx, token.wtree[0].leny, token.wtree[0].x, token.wtree[0].y, 0, 0)
	qtree16(token, 35, token.wtree[5].lenx, token.wtree[5].leny, token.wtree[5].x, token.wtree[5].y, 1, 0)
	qtree4(token, 0, token.wtree[19].lenx, token.wtree[19].leny, token.wtree[19].x, token.wtree[19].y)
}

func qtree16(token *token, start, lenx, leny, x, y, rw, cl int) {
	var tempx, temp2x int /* temporary x values */
	var tempy, temp2y int /* temporary y values */
	var evenx, eveny int  /* Check length of subband for even or odd */
	var p int             /* indicates subband information being stored */

	p = start
	evenx = lenx % 2
	eveny = leny % 2

	if evenx == 0 {
		tempx = lenx / 2
		temp2x = tempx
	} else {
		if cl != 0 {
			temp2x = (lenx + 1) / 2
			tempx = temp2x - 1
		} else {
			tempx = (lenx + 1) / 2
			temp2x = tempx - 1
		}
	}

	if eveny == 0 {
		tempy = leny / 2
		temp2y = tempy
	} else {
		if rw != 0 {
			temp2y = (leny + 1) / 2
			tempy = temp2y - 1
		} else {
			tempy = (leny + 1) / 2
			temp2y = tempy - 1
		}
	}

	evenx = tempx % 2
	eveny = tempy % 2

	token.qtree[p].x = x
	token.qtree[p+2].x = x
	token.qtree[p].y = y
	token.qtree[p+1].y = y
	if evenx == 0 {
		token.qtree[p].lenx = tempx / 2
		token.qtree[p+1].lenx = token.qtree[p].lenx
		token.qtree[p+2].lenx = token.qtree[p].lenx
		token.qtree[p+3].lenx = token.qtree[p].lenx
	} else {
		token.qtree[p].lenx = (tempx + 1) / 2
		token.qtree[p+1].lenx = token.qtree[p].lenx - 1
		token.qtree[p+2].lenx = token.qtree[p].lenx
		token.qtree[p+3].lenx = token.qtree[p+1].lenx
	}
	token.qtree[p+1].x = x + token.qtree[p].lenx
	token.qtree[p+3].x = token.qtree[p+1].x
	if eveny == 0 {
		token.qtree[p].leny = tempy / 2
		token.qtree[p+1].leny = token.qtree[p].leny
		token.qtree[p+2].leny = token.qtree[p].leny
		token.qtree[p+3].leny = token.qtree[p].leny
	} else {
		token.qtree[p].leny = (tempy + 1) / 2
		token.qtree[p+1].leny = token.qtree[p].leny
		token.qtree[p+2].leny = token.qtree[p].leny - 1
		token.qtree[p+3].leny = token.qtree[p+2].leny
	}
	token.qtree[p+2].y = y + token.qtree[p].leny
	token.qtree[p+3].y = token.qtree[p+2].y

	evenx = temp2x % 2

	token.qtree[p+4].x = x + tempx
	token.qtree[p+6].x = token.qtree[p+4].x
	token.qtree[p+4].y = y
	token.qtree[p+5].y = y
	token.qtree[p+6].y = token.qtree[p+2].y
	token.qtree[p+7].y = token.qtree[p+2].y
	token.qtree[p+4].leny = token.qtree[p].leny
	token.qtree[p+5].leny = token.qtree[p].leny
	token.qtree[p+6].leny = token.qtree[p+2].leny
	token.qtree[p+7].leny = token.qtree[p+2].leny
	if evenx == 0 {
		token.qtree[p+4].lenx = temp2x / 2
		token.qtree[p+5].lenx = token.qtree[p+4].lenx
		token.qtree[p+6].lenx = token.qtree[p+4].lenx
		token.qtree[p+7].lenx = token.qtree[p+4].lenx
	} else {
		token.qtree[p+5].lenx = (temp2x + 1) / 2
		token.qtree[p+4].lenx = token.qtree[p+5].lenx - 1
		token.qtree[p+6].lenx = token.qtree[p+4].lenx
		token.qtree[p+7].lenx = token.qtree[p+5].lenx
	}
	token.qtree[p+5].x = token.qtree[p+4].x + token.qtree[p+4].lenx
	token.qtree[p+7].x = token.qtree[p+5].x

	eveny = temp2y % 2

	token.qtree[p+8].x = x
	token.qtree[p+9].x = token.qtree[p+1].x
	token.qtree[p+10].x = x
	token.qtree[p+11].x = token.qtree[p+1].x
	token.qtree[p+8].y = y + tempy
	token.qtree[p+9].y = token.qtree[p+8].y
	token.qtree[p+8].lenx = token.qtree[p].lenx
	token.qtree[p+9].lenx = token.qtree[p+1].lenx
	token.qtree[p+10].lenx = token.qtree[p].lenx
	token.qtree[p+11].lenx = token.qtree[p+1].lenx
	if eveny == 0 {
		token.qtree[p+8].leny = temp2y / 2
		token.qtree[p+9].leny = token.qtree[p+8].leny
		token.qtree[p+10].leny = token.qtree[p+8].leny
		token.qtree[p+11].leny = token.qtree[p+8].leny
	} else {
		token.qtree[p+10].leny = (temp2y + 1) / 2
		token.qtree[p+11].leny = token.qtree[p+10].leny
		token.qtree[p+8].leny = token.qtree[p+10].leny - 1
		token.qtree[p+9].leny = token.qtree[p+8].leny
	}
	token.qtree[p+10].y = token.qtree[p+8].y + token.qtree[p+8].leny
	token.qtree[p+11].y = token.qtree[p+10].y

	token.qtree[p+12].x = token.qtree[p+4].x
	token.qtree[p+13].x = token.qtree[p+5].x
	token.qtree[p+14].x = token.qtree[p+4].x
	token.qtree[p+15].x = token.qtree[p+5].x
	token.qtree[p+12].y = token.qtree[p+8].y
	token.qtree[p+13].y = token.qtree[p+8].y
	token.qtree[p+14].y = token.qtree[p+10].y
	token.qtree[p+15].y = token.qtree[p+10].y
	token.qtree[p+12].lenx = token.qtree[p+4].lenx
	token.qtree[p+13].lenx = token.qtree[p+5].lenx
	token.qtree[p+14].lenx = token.qtree[p+4].lenx
	token.qtree[p+15].lenx = token.qtree[p+5].lenx
	token.qtree[p+12].leny = token.qtree[p+8].leny
	token.qtree[p+13].leny = token.qtree[p+8].leny
	token.qtree[p+14].leny = token.qtree[p+10].leny
	token.qtree[p+15].leny = token.qtree[p+10].leny
}

func qtree4(token *token, start, lenx, leny, x, y int) {
	var evenx, eveny int /* Check length of subband for even or odd */
	var p int            /* indicates subband information being stored */

	p = start
	evenx = lenx % 2
	eveny = leny % 2

	token.qtree[p].x = x
	token.qtree[p+2].x = x
	token.qtree[p].y = y
	token.qtree[p+1].y = y
	if evenx == 0 {
		token.qtree[p].lenx = lenx / 2
		token.qtree[p+1].lenx = token.qtree[p].lenx
		token.qtree[p+2].lenx = token.qtree[p].lenx
		token.qtree[p+3].lenx = token.qtree[p].lenx
	} else {
		token.qtree[p].lenx = (lenx + 1) / 2
		token.qtree[p+1].lenx = token.qtree[p].lenx - 1
		token.qtree[p+2].lenx = token.qtree[p].lenx
		token.qtree[p+3].lenx = token.qtree[p+1].lenx
	}
	token.qtree[p+1].x = x + token.qtree[p].lenx
	token.qtree[p+3].x = token.qtree[p+1].x
	if eveny == 0 {
		token.qtree[p].leny = leny / 2
		token.qtree[p+1].leny = token.qtree[p].leny
		token.qtree[p+2].leny = token.qtree[p].leny
		token.qtree[p+3].leny = token.qtree[p].leny
	} else {
		token.qtree[p].leny = (leny + 1) / 2
		token.qtree[p+1].leny = token.qtree[p].leny
		token.qtree[p+2].leny = token.qtree[p].leny - 1
		token.qtree[p+3].leny = token.qtree[p+2].leny
	}
	token.qtree[p+2].y = y + token.qtree[p].leny
	token.qtree[p+3].y = token.qtree[p+2].y
}

func huffmanDecodeDataMem(token *token, size int) ([]int, error) {
	qdata := make([]int, size)

	maxcode := make([]int, max_huffbits+1)
	mincode := make([]int, max_huffbits+1)
	valptr := make([]int, max_huffbits+1)

	ref, err := token.getCMarkerWSQ(tbls_n_sob)
	if err != nil {
		return nil, err
	}
	marker := &reference[int]{ref}

	bitCount := &reference[int]{0} /* bit count for getc_nextbits_wsq routine */
	nextByte := &reference[int]{0} /*next byte of buffer*/
	var hufftableId int            /* huffman table number */
	var ip int

	for marker.value != eoi_wsq {
		if marker.value != 0 {
			for marker.value != sob_wsq {
				if err := token.getCTableWSQ(marker.value); err != nil {
					return nil, err
				}
				marker.value, err = token.getCMarkerWSQ(tbls_n_sob)
				if err != nil {
					return nil, err
				}
				if marker.value == eoi_wsq {
					break
				}
			}
			if marker.value == eoi_wsq {
				break
			}
			hufftableId = getCBlockHeader(token) /* huffman table number */

			if token.tableDHT[hufftableId].tabdef != 1 {
				return nil, fmt.Errorf("ERROR : huffmanDecodeDataMem : huffman table undefined.")
			}

			/* the next two routines reconstruct the huffman tables */
			hufftable := buildHuffsizes(token.tableDHT[hufftableId].huffbits, max_huffcounts_wsq)
			buildHuffcodes(hufftable)

			/* this routine builds a set of three tables used in decoding */
			/* the compressed buffer*/
			genDecodeTable(hufftable, maxcode, mincode, valptr, token.tableDHT[hufftableId].huffbits)

			bitCount.value = 0
			marker.value = 0
		}

		/* get next huffman category code from compressed input buffer stream */
		nodeptr, err := decodeDataMem(token, mincode, maxcode, valptr, token.tableDHT[hufftableId].huffvalues, bitCount, marker, nextByte)
		if err != nil {
			return nil, err
		}
		/* nodeptr  pointers for decoding */

		if nodeptr == -1 {
			continue
		}

		if nodeptr > 0 && nodeptr <= 100 {
			for n := 0; n < nodeptr; n++ {
				qdata[ip] = 0 /* z run */
				ip++
			}
		} else if nodeptr > 106 && nodeptr < 0xff {
			qdata[ip] = nodeptr - 180
			ip++
		} else if nodeptr == 101 {
			v, err := getCNextbitsWSQ(token, marker, bitCount, 8, nextByte)
			if err != nil {
				return nil, err
			}
			qdata[ip] = v
			ip++
		} else if nodeptr == 102 {
			v, err := getCNextbitsWSQ(token, marker, bitCount, 8, nextByte)
			if err != nil {
				return nil, err
			}
			qdata[ip] = -v
			ip++
		} else if nodeptr == 103 {
			v, err := getCNextbitsWSQ(token, marker, bitCount, 16, nextByte)
			if err != nil {
				return nil, err
			}
			qdata[ip] = v
			ip++
		} else if nodeptr == 104 {
			v, err := getCNextbitsWSQ(token, marker, bitCount, 16, nextByte)
			if err != nil {
				return nil, err
			}
			qdata[ip] = -v
			ip++
		} else if nodeptr == 105 {
			n, err := getCNextbitsWSQ(token, marker, bitCount, 8, nextByte)
			if err != nil {
				return nil, err
			}
			for n > 0 {
				n--
				qdata[ip] = 0
				ip++
			}
		} else if nodeptr == 106 {
			n, err := getCNextbitsWSQ(token, marker, bitCount, 16, nextByte)
			if err != nil {
				return nil, err
			}
			for n > 0 {
				n--
				qdata[ip] = 0
				ip++
			}
		} else {
			return nil, fmt.Errorf("ERROR: huffman_decode_data_mem : Invalid code (%d)", nodeptr)
		}
	}

	return qdata, nil
}

func getCBlockHeader(token *token) int {
	token.readShort() /* block header size */
	return token.readByte()
}

func buildHuffsizes(huffbits []int, maxHuffcounts int) []huffCode {
	var huffcodeTable []huffCode /*table of huffman codes and sizes*/
	numberOfCodes := 1           /*the number codes for a given code size*/

	huffcodeTable = make([]huffCode, maxHuffcounts+1)

	var tempSize int

	for codeSize := 1; codeSize <= max_huffbits; codeSize++ {
		for numberOfCodes <= huffbits[codeSize-1] {
			huffcodeTable[tempSize] = huffCode{}
			huffcodeTable[tempSize].size = codeSize
			tempSize++
			numberOfCodes++
		}
		numberOfCodes = 1
	}

	huffcodeTable[tempSize] = huffCode{}
	huffcodeTable[tempSize].size = 0

	return huffcodeTable
}

func buildHuffcodes(huffcodeTable []huffCode) {
	var tempCode int /*used to construct code word*/
	var pointer int  /*pointer to code word information*/

	tempSize := huffcodeTable[0].size
	if huffcodeTable[pointer].size == 0 {
		return
	}

	for {
		for {
			huffcodeTable[pointer].code = tempCode
			tempCode++
			pointer++
			if huffcodeTable[pointer].size != tempSize {
				break
			}
		}

		if huffcodeTable[pointer].size == 0 {
			return
		}

		for {
			tempCode <<= 1
			tempSize++
			if huffcodeTable[pointer].size == tempSize {
				break
			}
		}
		if huffcodeTable[pointer].size != tempSize {
			break
		}
	}
}

func genDecodeTable(huffcodeTable []huffCode, maxcode, mincode, valptr, huffbits []int) {
	for i := 0; i <= max_huffbits; i++ {
		maxcode[i] = 0
		mincode[i] = 0
		valptr[i] = 0
	}

	var i2 int
	for i := 1; i <= max_huffbits; i++ {
		if huffbits[i-1] == 0 {
			maxcode[i] = -1
			continue
		}
		valptr[i] = i2
		mincode[i] = huffcodeTable[i2].code
		i2 = i2 + huffbits[i-1] - 1
		maxcode[i] = huffcodeTable[i2].code
		i2++
	}
}

func decodeDataMem(token *token, mincode, maxcode, valptr, huffvalues []int, bitCount, marker, nextByte *reference[int]) (int, error) {
	code, err := getCNextbitsWSQ(token, marker, bitCount, 1, nextByte) /* becomes a huffman code word  (one bit at a time)*/
	if err != nil {
		return 0, err
	}
	if marker.value != 0 {
		return -1, nil
	}

	var inx int
	for inx = 1; code > maxcode[inx]; inx++ {
		tbits, err := getCNextbitsWSQ(token, marker, bitCount, 1, nextByte) /* becomes a huffman code word  (one bit at a time)*/
		if err != nil {
			return 0, err
		}
		code = ((code << 1) + tbits)

		if marker.value != 0 {
			return -1, nil
		}
	}

	inx2 := valptr[inx] + code - mincode[inx] /*increment variables*/
	return huffvalues[inx2], nil
}

func getCNextbitsWSQ(token *token, marker, bitCount *reference[int], bitsReq int, nextByte *reference[int]) (int, error) {
	if bitCount.value == 0 {
		nextByte.value = token.readByte()

		bitCount.value = 8
		if nextByte.value == 0xFF {
			code2 := token.readByte() /*stuffed byte of buffer*/

			if code2 != 0x00 && bitsReq == 1 {
				marker.value = (nextByte.value << 8) | code2
				return 1, nil
			}
			if code2 != 0x00 {
				return 0, fmt.Errorf("ERROR: getCNextbitsWSQ : No stuffed zeros.")
			}
		}
	}

	var bits, tbits int /*bits of current buffer byte requested*/
	var bitsNeeded int  /*additional bits required to finish request*/
	var err error
	if bitsReq <= bitCount.value {
		bits = (nextByte.value >> (bitCount.value - bitsReq)) & (bitmask[bitsReq])
		bitCount.value -= bitsReq
		nextByte.value &= bitmask[bitCount.value]
	} else {
		bitsNeeded = bitsReq - bitCount.value /*additional bits required to finish request*/
		bits = nextByte.value << bitsNeeded
		bitCount.value = 0
		tbits, err = getCNextbitsWSQ(token, marker, bitCount, bitsNeeded, nextByte)
		if err != nil {
			return 0, err
		}
		bits |= tbits
	}

	return bits, nil
}

func unquantize(token *token, sip []int, width, height int) ([]float32, error) {
	fip := make([]float32, width*height) /* floating point image */

	if token.tableDQT.dqtDef != 1 {
		return nil, fmt.Errorf("ERROR: unquantize : quantization table parameters not defined!")
	}

	binCenter := token.tableDQT.binCenter /* quantizer bin center */

	var sptr int
	for cnt := 0; cnt < num_subbands; cnt++ {
		if token.tableDQT.qBin[cnt] == 0.0 {
			continue
		}

		fptr := (token.qtree[cnt].y * width) + token.qtree[cnt].x

		for row := 0; row < token.qtree[cnt].leny; row++ {
			for col := 0; col < token.qtree[cnt].lenx; col++ {
				if sip[sptr] == 0 {
					fip[fptr] = 0
				} else if sip[sptr] > 0 {
					fip[fptr] = (token.tableDQT.qBin[cnt] * (float32(sip[sptr]) - binCenter)) + (token.tableDQT.zBin[cnt] / 2.0)
				} else if sip[sptr] < 0 {
					fip[fptr] = (token.tableDQT.qBin[cnt] * (float32(sip[sptr]) + binCenter)) - (token.tableDQT.zBin[cnt] / 2.0)
				} else {
					return nil, fmt.Errorf("ERROR : unquantize : invalid quantization pixel value")
				}
				fptr++
				sptr++
			}
			fptr += width - token.qtree[cnt].lenx
		}
	}

	return fip, nil
}

func wsqReconstruct(token *token, fdata []float32, width, height int) error {
	if token.tableDTT.lodef != 1 {
		return fmt.Errorf("ERROR: wsq_reconstruct : Lopass filter coefficients not defined")
	}

	if token.tableDTT.hidef != 1 {
		return fmt.Errorf("ERROR: wsq_reconstruct : Hipass filter coefficients not defined")
	}

	numPix := width * height
	/* Allocate temporary floating point pixmap. */
	fdataTemp := make([]float32, numPix)

	/* Reconstruct floating point pixmap from wavelet subband buffer. */
	for node := w_treelen - 1; node >= 0; node-- {
		fdataBse := (token.wtree[node].y * width) + token.wtree[node].x
		joinLets(fdataTemp, fdata, 0, fdataBse, token.wtree[node].lenx, token.wtree[node].leny,
			1, width,
			token.tableDTT.hifilt, token.tableDTT.hisz,
			token.tableDTT.lofilt, token.tableDTT.losz,
			token.wtree[node].invcl)
		joinLets(fdata, fdataTemp, fdataBse, 0, token.wtree[node].leny, token.wtree[node].lenx,
			width, 1,
			token.tableDTT.hifilt, token.tableDTT.hisz,
			token.tableDTT.lofilt, token.tableDTT.losz,
			token.wtree[node].invrw)
	}

	return nil
}

func joinLets(
	newdata,
	olddata []float32,
	newIndex,
	oldIndex,
	len1, /* temporary length parameters */
	len2,
	pitch, /* pitch gives next row_col to filter */
	stride int, /*           stride gives next pixel to filter */
	hi []float32,
	hsz int, /* NEW */
	lo []float32, /* filter coefficients */
	lsz, /* NEW */
	inv int) /* spectral inversion? */ {
	var lp0, lp1 int
	var hp0, hp1 int
	var lopass, hipass int /* lo/hi pass image pointers */
	var limg, himg int
	var pix, cl_rw int /* pixel counter and column/row counter */
	var i, da_ev int   /* if "scanline" is even or odd and */
	var loc, hoc int
	var hlen, llen int
	var nstr, pstr int
	var tap int
	var fi_ev int
	var olle, ohle, olre, ohre int
	var lle, lle2, lre, lre2 int
	var hle, hle2, hre, hre2 int
	var lpx, lspx int
	var lpxstr, lspxstr int
	var lstap, lotap int
	var hpx, hspx int
	var hpxstr, hspxstr int
	var hstap, hotap int
	var asym, fhre, ofhre int
	var ssfac, osfac, sfac float32

	da_ev = len2 % 2
	fi_ev = lsz % 2
	pstr = stride
	nstr = -pstr
	if da_ev != 0 {
		llen = (len2 + 1) / 2
		hlen = llen - 1
	} else {
		llen = len2 / 2
		hlen = llen
	}

	if fi_ev != 0 {
		asym = 0
		ssfac = 1.0
		ofhre = 0
		loc = (lsz - 1) / 4
		hoc = (hsz+1)/4 - 1
		lotap = ((lsz - 1) / 2) % 2
		hotap = ((hsz + 1) / 2) % 2
		if da_ev != 0 {
			olle = 0
			olre = 0
			ohle = 1
			ohre = 1
		} else {
			olle = 0
			olre = 1
			ohle = 1
			ohre = 0
		}
	} else {
		asym = 1
		ssfac = -1.0
		ofhre = 2
		loc = lsz/4 - 1
		hoc = hsz/4 - 1
		lotap = (lsz / 2) % 2
		hotap = (hsz / 2) % 2
		if da_ev != 0 {
			olle = 1
			olre = 0
			ohle = 1
			ohre = 1
		} else {
			olle = 1
			olre = 1
			ohle = 1
			ohre = 1
		}

		if loc == -1 {
			loc = 0
			olle = 0
		}
		if hoc == -1 {
			hoc = 0
			ohle = 0
		}

		for i = 0; i < hsz; i++ {
			hi[i] *= -1.0
		}
	}

	for cl_rw = 0; cl_rw < len1; cl_rw++ {
		limg = newIndex + cl_rw*pitch
		himg = limg
		newdata[himg] = 0.0
		newdata[himg+stride] = 0.0
		if inv != 0 {
			hipass = oldIndex + cl_rw*pitch
			lopass = hipass + stride*hlen
		} else {
			lopass = oldIndex + cl_rw*pitch
			hipass = lopass + stride*llen
		}

		lp0 = lopass
		lp1 = lp0 + (llen-1)*stride
		lspx = lp0 + (loc * stride)
		lspxstr = nstr
		lstap = lotap
		lle2 = olle
		lre2 = olre

		hp0 = hipass
		hp1 = hp0 + (hlen-1)*stride
		hspx = hp0 + (hoc * stride)
		hspxstr = nstr
		hstap = hotap
		hle2 = ohle
		hre2 = ohre
		osfac = ssfac

		for pix = 0; pix < hlen; pix++ {
			for tap = lstap; tap >= 0; tap-- {
				lle = lle2
				lre = lre2
				lpx = lspx
				lpxstr = lspxstr

				newdata[limg] = olddata[lpx] * lo[tap]
				for i = tap + 2; i < lsz; i += 2 {
					if lpx == lp0 {
						if lle != 0 {
							lpxstr = 0
							lle = 0
						} else {
							lpxstr = pstr
						}
					}
					if lpx == lp1 {
						if lre != 0 {
							lpxstr = 0
							lre = 0
						} else {
							lpxstr = nstr
						}
					}
					lpx += lpxstr
					newdata[limg] += olddata[lpx] * lo[i]
				}
				limg += stride
			}
			if lspx == lp0 {
				if lle2 != 0 {
					lspxstr = 0
					lle2 = 0
				} else {
					lspxstr = pstr
				}
			}
			lspx += lspxstr
			lstap = 1

			for tap = hstap; tap >= 0; tap-- {
				hle = hle2
				hre = hre2
				hpx = hspx
				hpxstr = hspxstr
				fhre = ofhre
				sfac = osfac

				for i = tap; i < hsz; i += 2 {
					if hpx == hp0 {
						if hle != 0 {
							hpxstr = 0
							hle = 0
						} else {
							hpxstr = pstr
							sfac = 1.0
						}
					}
					if hpx == hp1 {
						if hre != 0 {
							hpxstr = 0
							hre = 0
							if asym != 0 && da_ev != 0 {
								hre = 1
								fhre--
								sfac = float32(fhre)
								if sfac == 0.0 {
									hre = 0
								}
							}
						} else {
							hpxstr = nstr
							if asym != 0 {
								sfac = -1.0
							}
						}
					}
					newdata[himg] += olddata[hpx] * hi[i] * sfac
					hpx += hpxstr
				}
				himg += stride
			}
			if hspx == hp0 {
				if hle2 != 0 {
					hspxstr = 0
					hle2 = 0
				} else {
					hspxstr = pstr
					osfac = 1.0
				}
			}
			hspx += hspxstr
			hstap = 1
		}

		if da_ev != 0 {
			if lotap != 0 {
				lstap = 1
			} else {
				lstap = 0
			}
		} else if lotap != 0 {
			lstap = 2
		} else {
			lstap = 1
		}

		for tap = 1; tap >= lstap; tap-- {
			lle = lle2
			lre = lre2
			lpx = lspx
			lpxstr = lspxstr

			newdata[limg] = olddata[lpx] * lo[tap]
			for i = tap + 2; i < lsz; i += 2 {
				if lpx == lp0 {
					if lle != 0 {
						lpxstr = 0
						lle = 0
					} else {
						lpxstr = pstr
					}
				}
				if lpx == lp1 {
					if lre != 0 {
						lpxstr = 0
						lre = 0
					} else {
						lpxstr = nstr
					}
				}
				lpx += lpxstr
				newdata[limg] += olddata[lpx] * lo[i]
			}
			limg += stride
		}

		if da_ev != 0 {
			if hotap != 0 {
				hstap = 1
			} else {
				hstap = 0
			}
			if hsz == 2 {
				hspx -= hspxstr
				fhre = 1
			}
		} else if hotap != 0 {
			hstap = 2
		} else {
			hstap = 1
		}

		for tap = 1; tap >= hstap; tap-- {
			hle = hle2
			hre = hre2
			hpx = hspx
			hpxstr = hspxstr
			sfac = osfac
			if hsz != 2 {
				fhre = ofhre
			}

			for i = tap; i < hsz; i += 2 {
				if hpx == hp0 {
					if hle != 0 {
						hpxstr = 0
						hle = 0
					} else {
						hpxstr = pstr
						sfac = 1.0
					}
				}
				if hpx == hp1 {
					if hre != 0 {
						hpxstr = 0
						hre = 0
						if asym != 0 && da_ev != 0 {
							hre = 1
							fhre--
							sfac = float32(fhre)
							if sfac == 0.0 {
								hre = 0
							}
						}
					} else {
						hpxstr = nstr
						if asym != 0 {
							sfac = -1.0
						}
					}
				}
				newdata[himg] += olddata[hpx] * hi[i] * sfac
				hpx += hpxstr
			}
			himg += stride
		}
	}

	if fi_ev == 0 {
		for i = 0; i < hsz; i++ {
			hi[i] *= -1.0
		}
	}
}
