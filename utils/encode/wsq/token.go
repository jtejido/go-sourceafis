package wsq

import (
	"fmt"
	"io"
)

const (
	max_dht_tables     = 8
	max_subbands       = 64
	max_huffbits       = 16
	max_huffcounts_wsq = 256
	w_treelen          = 20
	q_treelen          = 64

	/* wsq marker definitions */
	soi_wsq = 0xffa0
	eoi_wsq = 0xffa1
	sof_wsq = 0xffa2
	sob_wsq = 0xffa3
	dtt_wsq = 0xffa4
	dqt_wsq = 0xffa5
	dht_wsq = 0xffa6
	drt_wsq = 0xffa7
	com_wsq = 0xffa8

	strt_subband_2     = 19
	strt_subband_3     = 52
	num_subbands       = 60
	strt_subband_del   = num_subbands
	strt_size_region_2 = 4
	strt_size_region_3 = 51

	/* case for getting any marker. */
	any_wsq    = 0xffff
	tbls_n_sof = 2
	tbls_n_sob = tbls_n_sof + 2
)

var bitmask = []int{0x00, 0x01, 0x03, 0x07, 0x0f, 0x1f, 0x3f, 0x7f, 0xff}

type tableDTT struct {
	lofilt                   []float32
	hifilt                   []float32
	losz, hisz, lodef, hidef int
}
type tableDQT struct {
	binCenter float32
	qBin      [max_subbands]float32
	zBin      [max_subbands]float32
	dqtDef    rune
}

type huffmanTable struct {
	tableLen             int
	bytesLeft            int
	tableId              int
	huffbits, huffvalues []int
}

type tableDHT struct {
	tabdef     byte
	huffbits   []int // MAX_HUFFBITS
	huffvalues []int // MAX_HUFFCOUNTS_WSQ + 1
}
type wavletTree struct {
	x, y, lenx, leny, invrw, invcl int
}

type quantTree struct {
	x, y, lenx, leny int
}

type number interface {
	int | float32
}

type reference[V number] struct {
	value V
}

type headerFrm struct {
	black, white, width, height int
	mShift, rScale              float32
	wsqEncoder, software        int
}

type huffCode struct {
	size int
	code int
}

func readByte(r io.ByteReader) (byte, error) {
	b, err := r.ReadByte()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return b, err
}

type token struct {
	tableDHT []*tableDHT
	tableDTT tableDTT
	tableDQT tableDQT
	wtree    []wavletTree
	qtree    []quantTree
	buffer   []byte
	pointer  int
}

func newToken(buffer []byte) *token {
	tok := &token{
		buffer:   buffer,
		tableDTT: tableDTT{},
		tableDQT: tableDQT{},
		tableDHT: make([]*tableDHT, max_dht_tables),
	}
	for i := 0; i < max_dht_tables; i++ {
		tok.tableDHT[i] = new(tableDHT)
		tok.tableDHT[i].tabdef = 0
		tok.tableDHT[i].huffbits = make([]int, max_huffbits)
		tok.tableDHT[i].huffvalues = make([]int, max_huffcounts_wsq+1)
	}

	return tok
}

func (t *token) getCTableWSQ(marker int) error {
	switch marker {
	case dtt_wsq:
		t.getCTransformTable()
		return nil
	case dqt_wsq:
		t.getCQuantizationTable()
		return nil
	case dht_wsq:
		t.getCHuffmanTableWSQ()
		return nil
	case com_wsq:
		//shams: i don't use return value
		t.getCComment()
		return nil
	default:
		return fmt.Errorf("ERROR: getCTableWSQ : Invalid table defined : %d", marker)
	}
}

func (t *token) getCFrameHeaderWSQ() headerFrm {
	headerFrm := headerFrm{}

	t.readShort() /* header size */

	headerFrm.black = t.readByte()
	headerFrm.white = t.readByte()
	headerFrm.height = t.readShort()
	headerFrm.width = t.readShort()
	scale := t.readByte()    /* exponent scaling parameter */
	shrtDat := t.readShort() /* buffer pointer */
	headerFrm.mShift = float32(shrtDat)
	for scale > 0 {
		headerFrm.mShift /= 10.0
		scale--
	}

	scale = t.readByte()
	shrtDat = t.readShort()
	headerFrm.rScale = float32(shrtDat)
	for scale > 0 {
		headerFrm.rScale /= 10.0
		scale--
	}

	headerFrm.wsqEncoder = t.readByte()
	headerFrm.software = t.readShort()

	return headerFrm
}

func (t *token) getCComment() string {
	size := t.readShort() - 2
	return string(t.readBytes(size))
}

func (t *token) getCTransformTable() {
	// read header Size;
	t.readShort()

	t.tableDTT.hisz = t.readByte()
	t.tableDTT.losz = t.readByte()

	t.tableDTT.hifilt = make([]float32, t.tableDTT.hisz)
	t.tableDTT.lofilt = make([]float32, t.tableDTT.losz)

	var aSize int
	if t.tableDTT.hisz%2 != 0 {
		aSize = (t.tableDTT.hisz + 1) / 2
	} else {
		aSize = t.tableDTT.hisz / 2
	}

	aLofilt := make([]float32, aSize)

	aSize--
	for cnt := 0; cnt <= aSize; cnt++ {
		sign := t.readByte()
		scale := t.readByte()
		shrtDat := t.readLong()
		aLofilt[cnt] = float32(shrtDat)

		for scale > 0 {
			aLofilt[cnt] /= 10.0
			scale--
		}

		if sign != 0 {
			aLofilt[cnt] *= -1.0
		}

		if t.tableDTT.hisz%2 != 0 {
			t.tableDTT.hifilt[cnt+aSize] = float32(intSign(cnt)) * aLofilt[cnt]
			if cnt > 0 {
				t.tableDTT.hifilt[aSize-cnt] = t.tableDTT.hifilt[cnt+aSize]
			}
		} else {
			t.tableDTT.hifilt[cnt+aSize+1] = float32(intSign(cnt)) * aLofilt[cnt]
			t.tableDTT.hifilt[aSize-cnt] = -1 * t.tableDTT.hifilt[cnt+aSize+1]
		}
	}

	if t.tableDTT.losz%2 != 0 {
		aSize = (t.tableDTT.losz + 1) / 2
	} else {
		aSize = t.tableDTT.losz / 2
	}

	aHifilt := make([]float32, aSize)

	aSize--
	for cnt := 0; cnt <= aSize; cnt++ {
		sign := t.readByte()
		scale := t.readByte()
		shrtDat := t.readLong()

		aHifilt[cnt] = float32(shrtDat)

		for scale > 0 {
			aHifilt[cnt] /= 10.0
			scale--
		}

		if sign != 0 {
			aHifilt[cnt] *= -1.0
		}

		if t.tableDTT.losz%2 != 0 {
			t.tableDTT.lofilt[cnt+aSize] = float32(intSign(cnt)) * aHifilt[cnt]
			if cnt > 0 {
				t.tableDTT.lofilt[aSize-cnt] = t.tableDTT.lofilt[cnt+aSize]
			}
		} else {
			t.tableDTT.lofilt[cnt+aSize+1] = float32(intSign(cnt+1)) * aHifilt[cnt]
			t.tableDTT.lofilt[aSize-cnt] = t.tableDTT.lofilt[cnt+aSize+1]
		}
	}

	t.tableDTT.lodef = 1
	t.tableDTT.hidef = 1
}

func (t *token) getCQuantizationTable() {
	t.readShort()            /* header size */
	scale := t.readByte()    /* scaling parameter */
	shrtDat := t.readShort() /* counter and temp short buffer */

	t.tableDQT.binCenter = float32(shrtDat)
	for scale > 0 {
		t.tableDQT.binCenter /= 10.0
		scale--
	}

	for cnt := 0; cnt < max_subbands; cnt++ {
		scale = t.readByte()
		shrtDat = t.readShort()
		t.tableDQT.qBin[cnt] = float32(shrtDat)
		for scale > 0 {
			t.tableDQT.qBin[cnt] /= 10.0
			scale--
		}

		scale = t.readByte()
		shrtDat = t.readShort()
		t.tableDQT.zBin[cnt] = float32(shrtDat)
		for scale > 0 {
			t.tableDQT.zBin[cnt] /= 10.0
			scale--
		}
	}

	t.tableDQT.dqtDef = 1
}

func (t *token) getCHuffmanTableWSQ() error {
	/* First time, read table len. */
	firstHuffmanTable, err := t.getCHuffmanTable(max_huffcounts_wsq, 0, true)
	if err != nil {
		return err
	}

	/* Store table into global structure list. */
	tableId := firstHuffmanTable.tableId
	copy(t.tableDHT[tableId].huffbits, firstHuffmanTable.huffbits)
	copy(t.tableDHT[tableId].huffvalues, firstHuffmanTable.huffvalues)
	t.tableDHT[tableId].tabdef = 1

	bytesLeft := firstHuffmanTable.bytesLeft
	for bytesLeft != 0 {
		/* Read next table without reading table len. */
		huffmantable, err := t.getCHuffmanTable(max_huffcounts_wsq, bytesLeft, false)
		if err != nil {
			return err
		}

		/* If table is already defined ... */
		tableId = huffmantable.tableId
		if t.tableDHT[tableId].tabdef != 0 {
			return fmt.Errorf("ERROR : getCHuffmanTableWSQ : huffman table already defined.")
		}

		/* Store table into global structure list. */
		copy(t.tableDHT[tableId].huffbits, huffmantable.huffbits)
		copy(t.tableDHT[tableId].huffvalues, huffmantable.huffvalues)
		t.tableDHT[tableId].tabdef = 1
		bytesLeft = huffmantable.bytesLeft
	}

	return nil
}

func (t *token) getCHuffmanTable(maxHuffcounts, bytesLeft int, readTableLen bool) (*huffmanTable, error) {
	huffmanTable := &huffmanTable{}

	/* table_len */
	if readTableLen {
		huffmanTable.tableLen = t.readShort()
		huffmanTable.bytesLeft = huffmanTable.tableLen - 2
		bytesLeft = huffmanTable.bytesLeft
	} else {
		huffmanTable.bytesLeft = bytesLeft
	}

	/* If no bytes left ... */
	if bytesLeft <= 0 {
		return nil, fmt.Errorf("ERROR : getCHuffmanTable : no huffman table bytes remaining")
	}

	/* Table ID */
	huffmanTable.tableId = t.readByte()
	huffmanTable.bytesLeft--

	huffmanTable.huffbits = make([]int, max_huffbits)
	var numHufvals int
	/* L1 ... L16 */
	for i := 0; i < max_huffbits; i++ {
		huffmanTable.huffbits[i] = t.readByte()
		numHufvals += huffmanTable.huffbits[i]
	}
	huffmanTable.bytesLeft -= max_huffbits

	if numHufvals > maxHuffcounts+1 {
		return nil, fmt.Errorf("ERROR : getCHuffmanTable : numHufvals is larger than MAX_HUFFCOUNTS")
	}

	/* Could allocate only the amount needed ... then we wouldn't */
	/* need to pass MAX_HUFFCOUNTS. */
	huffmanTable.huffvalues = make([]int, maxHuffcounts+1)

	/* V1,1 ... V16,16 */
	for i := 0; i < numHufvals; i++ {
		huffmanTable.huffvalues[i] = t.readByte()
	}
	huffmanTable.bytesLeft -= numHufvals

	return huffmanTable, nil
}

func (t *token) getCMarkerWSQ(tt int) (int, error) {
	if t.pointer >= len(t.buffer) {
		return 0, fmt.Errorf("invalid pointer : %d", t.pointer)
	}

	marker := t.readShort()

	switch tt {
	case soi_wsq:
		if marker != soi_wsq {
			return 0, fmt.Errorf("ERROR : getCMarkerWSQ : No SOI marker : %d", marker)
		}

		return marker, nil

	case tbls_n_sof:
		if marker != dtt_wsq && marker != dqt_wsq && marker != dht_wsq && marker != sof_wsq && marker != com_wsq && marker != eoi_wsq {
			return 0, fmt.Errorf("ERROR : getc_marker_wsq : No SOF, Table, or comment markers : %d", marker)
		}

		return marker, nil

	case tbls_n_sob:
		if marker != dtt_wsq && marker != dqt_wsq && marker != dht_wsq && marker != sob_wsq && marker != com_wsq && marker != eoi_wsq {
			return 0, fmt.Errorf("ERROR : getc_marker_wsq : No SOB, Table, or comment markers :  %d", marker)
		}
		return marker, nil
	case any_wsq:
		if (marker & 0xff00) != 0xff00 {
			return 0, fmt.Errorf("ERROR : getc_marker_wsq : no marker found : %d", marker)
		}

		/* Added by MDG on 03-07-05 */
		if (marker < soi_wsq) || (marker > com_wsq) {
			return 0, fmt.Errorf("ERROR : getc_marker_wsq : not a valid marker : %d", marker)
		}

		return marker, nil
	default:
		return 0, fmt.Errorf("ERROR : getc_marker_wsq : Invalid marker : %d", marker)
	}
}

func (t *token) readLong() int64 {
	if t.pointer+4 > len(t.buffer) {
		// Handle buffer underflow or other error
		return 0
	}

	byte1 := int64(t.buffer[t.pointer])
	byte2 := int64(t.buffer[t.pointer+1])
	byte3 := int64(t.buffer[t.pointer+2])
	byte4 := int64(t.buffer[t.pointer+3])

	t.pointer += 4

	value := (0xff&byte1)<<24 | (0xff&byte2)<<16 | (0xff&byte3)<<8 | (0xff & byte4)
	return value
}

func (t *token) readShort() int {
	if t.pointer+2 > len(t.buffer) {
		// Handle buffer underflow or other error
		return 0
	}

	byte1 := int(t.buffer[t.pointer])
	byte2 := int(t.buffer[t.pointer+1])

	t.pointer += 2

	value := (0xff&byte1)<<8 | (0xff & byte2)
	return value
}

func (t *token) readByte() int {
	if t.pointer >= len(t.buffer) {
		// Handle buffer underflow or other error
		return 0
	}

	byte1 := int(t.buffer[t.pointer])
	t.pointer++

	value := 0xff & byte1
	return value
}

func (t *token) readBytes(size int) []byte {
	bytes := make([]byte, size)

	for i := 0; i < size; i++ {
		if t.pointer >= len(t.buffer) {
			// Handle buffer underflow or other error
			return nil
		}
		bytes[i] = t.buffer[t.pointer]
		t.pointer++
	}

	return bytes
}
