package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/deepch/dcodec"
	bits "github.com/deepch/old_bits"
)

const (
	NALU_RAW = iota
	NALU_AVCC
	NALU_ANNEXB
)

func main() {
	//encode decode example
	//open test img
	inJPEG, err := os.Open("in.jpeg")
	defer inJPEG.Close()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	//open out img
	outJPEG, err := os.Create("out.jpeg")
	defer outJPEG.Close()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	//decode image
	img, err := jpeg.Decode(inJPEG)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	//assert image.YCbCr
	imgYCbCr := img.(*image.YCbCr)
	//create out file
	outH264, err := os.Create("out.h264")
	defer outH264.Close()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	//init encoder
	encoder, err := dcodec.NewEncoder()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	//init decoder
	decoder, err := dcodec.NewDecoder()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	//encode 1000 image to h264 video
	i := 0
	var LastFrame *image.YCbCr
	for i <= 100 {
		buf, err := encoder.Encode(imgYCbCr)
		if err != nil {
			log.Println(err)
			continue
		}
		nalus, _ := SplitNALUs(buf)
		for _, nalu := range nalus {
			if len(nalu) > 0 {
				lastkeys := append([]byte("\000\000\000\001"), nalu...)
				outH264.Write(lastkeys)
				//ok test decode image
				decImage, err := decoder.Decode(lastkeys)
				if err != nil {
					log.Println(err)
					continue
				}
				//save last frame
				LastFrame = decImage
			}
		}
		i++
	}
	imgbuffer := new(bytes.Buffer)
	if err := jpeg.Encode(imgbuffer, LastFrame, nil); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	outJPEG.Write(imgbuffer.Bytes())
	//ffplay out.h264
}

func SplitNALUs(b []byte) (nalus [][]byte, typ int) {
	if len(b) < 4 {
		return [][]byte{b}, NALU_RAW
	}
	val3 := bits.GetUIntBE(b, 24)
	val4 := bits.GetUIntBE(b, 32)
	// maybe AVCC
	if val4 <= uint(len(b)) {
		_val4 := val4
		_b := b[4:]
		nalus := [][]byte{}
		for {
			nalus = append(nalus, _b[:_val4])
			if _val4 > uint(len(_b)) {
				break
			}
			_b = _b[_val4:]
			if len(_b) < 4 {
				break
			}
			_val4 = bits.GetUIntBE(_b, 32)
			_b = _b[4:]
			if _val4 > uint(len(_b)) {
				break
			}
		}
		if len(_b) == 0 {
			return nalus, NALU_AVCC
		}
	}
	// is Annex B
	if val3 == 1 || val4 == 1 {
		_val3 := val3
		_val4 := val4
		start := 0
		pos := 0
		for {
			if start != pos {
				nalus = append(nalus, b[start:pos])
			}
			if _val3 == 1 {
				pos += 3
			} else if _val4 == 1 {
				pos += 4
			}
			start = pos
			if start == len(b) {
				break
			}
			_val3 = 0
			_val4 = 0
			for pos < len(b) {
				if pos+2 < len(b) && b[pos] == 0 {
					_val3 = bits.GetUIntBE(b[pos:], 24)
					if _val3 == 0 {
						if pos+3 < len(b) {
							_val4 = uint(b[pos+3])
							if _val4 == 1 {
								break
							}
						}
					} else if _val3 == 1 {
						break
					}
					pos++
				} else {
					pos++
				}
			}
		}
		typ = NALU_ANNEXB
		return
	}
	return [][]byte{b}, NALU_RAW
}
