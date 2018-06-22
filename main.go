package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/deepch/dcodec"
	bits "github.com/deepch/old_bits"
	darknet "github.com/gyonluks/go-darknet"
)

const (
	NALU_RAW = iota
	NALU_AVCC
	NALU_ANNEXB
)

func main() {

	n := darknet.YOLONetwork{
		GPUDeviceIndex:           0,
		DataConfigurationFile:    "cfg/coco.data",
		NetworkConfigurationFile: "yolov3-tiny.cfg",
		WeightsFile:              "yolov3-tiny.weights",
		Threshold:                .5,
	}
	n2 := darknet.YOLONetwork{
		GPUDeviceIndex:           0,
		DataConfigurationFile:    "cfg/coco.data",
		NetworkConfigurationFile: "yolov3-tiny.cfg",
		WeightsFile:              "yolov3-tiny.weights",
		Threshold:                .5,
	}
	n3 := darknet.YOLONetwork{
		GPUDeviceIndex:           0,
		DataConfigurationFile:    "cfg/coco.data",
		NetworkConfigurationFile: "yolov3-tiny.cfg",
		WeightsFile:              "yolov3-tiny.weights",
		Threshold:                .5,
	}
	n4 := darknet.YOLONetwork{
		GPUDeviceIndex:           0,
		DataConfigurationFile:    "cfg/coco.data",
		NetworkConfigurationFile: "yolov3-tiny.cfg",
		WeightsFile:              "yolov3-tiny.weights",
		Threshold:                .5,
	}

	//log.Println(n)
	if err := n.Init(); err != nil {
		log.Println(err)
		return
	}
	if err := n2.Init(); err != nil {
		log.Println(err)
		return
	}
	if err := n3.Init(); err != nil {
		log.Println(err)
		return
	}
	if err := n4.Init(); err != nil {
		log.Println(err)
		return
	}
	//log.Println("deep")
	defer n.Close()
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
			log.Println("encode", err, "need more feed encoder?")
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
					naluTypefrec := nalu[0] & 0x1f
					log.Println("decode", naluTypefrec, err, "if 7 or 8 sps and pps it normal")
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
	imgs, err := darknet.ImageFromPath("139.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	imgs2, err := darknet.ImageFromPath("139.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	imgs3, err := darknet.ImageFromPath("139.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	imgs4, err := darknet.ImageFromPath("139.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	defer imgs.Close()
	go func() {
		for {
			//ddd := time.Now()

			dr, err := n2.Detect(imgs2)
			if err != nil || dr == nil {
				log.Println(err)
				return
			}
		}
	}()
	go func() {
		for {
			//ddd := time.Now()

			dr, err := n3.Detect(imgs3)
			if err != nil || dr == nil {
				log.Println(err)
				return
			}
		}
	}()
	go func() {
		for {
			//ddd := time.Now()

			dr, err := n4.Detect(imgs4)
			if err != nil || dr == nil {
				log.Println(err)
				return
			}
		}
	}()
	for {
		ddd := time.Now()

		dr, err := n.Detect(imgs)
		if err != nil || dr == nil {
			log.Println(err)
			return
		}
		//	var deep := dr
		//	log.Println("===>deep")
		//	log.Println("Network-only time taken:", dr.NetworkOnlyTimeTaken)
		//	log.Println("Overall time taken:", dr.OverallTimeTaken)
		//		for _, d := range dr.Detections {
		//			for i := range d.ClassIDs {
		//				bBox := d.BoundingBox
		//		fmt.Printf("%s (%d): %.4f%% | start point: (%d,%d) | end point: (%d, %d)\n",
		//		d.ClassNames[i], d.ClassIDs[i],
		//		d.Probabilities[i],
		//		bBox.StartPoint.X, bBox.StartPoint.Y,
		//		bBox.EndPoint.X, bBox.EndPoint.Y,
		//	)
		//		}
		//	}
		log.Println("====>", time.Now().Sub(ddd), len(dr.Detections))
		//	time.Sleep(100 * time.Millisecond)
	}
	//	log.Println("====>", time.Now().Sub(ddd), len(dr.Detections))
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
