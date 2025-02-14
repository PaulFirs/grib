package griblib

import (
	"image"
	"log"

	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
	"reflect"
)

func ExportMessagesAsPngs(messages []*Message) {
	for i, message := range messages {
		dataImage, err := imageFromMessage(message)
		if err != nil {
			log.Printf("Message could not be converted to image: %v\n", err)
		} else {
			writeImageToFilename(dataImage, imageFileName(i, message))
		}
	}
}

func ExportMessageAsPng(message *Message, filename string) error {
	dataImage, err := imageFromMessage(message)
	if err != nil {
		return err
	}
	return writeImageToFilename(dataImage, filename)
}

func imageFileName(messageNumber int, message *Message) string {
	dataname := ReadProductDisciplineParameters(message.Section0.Discipline, message.Section4.ProductDefinitionTemplate.ParameterCategory)
	return fmt.Sprintf("%s - discipline%d category%d messageIndex%d.png",
		dataname,
		message.Section0.Discipline,
		message.Section4.ProductDefinitionTemplate.ParameterCategory,
		messageNumber)
}

func writeImageToFilename(img image.Image, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func imageFromMessage(message *Message) (image.Image, error) {

	grid0, ok := message.Section3.Definition.(*Grid0)

	if !ok {
		err := fmt.Errorf("Currently not supporting definition of type %s ", reflect.TypeOf(message.Section3.Definition))
		return nil, err
	}

	height := int(grid0.Nj)
	width := int(grid0.Ni)

	maxValue, minValue := MaxMin(message.Section7.Data)

	rgbaImage := image.NewNRGBA(image.Rect(0, 0, width, height))
	length := len(message.Section7.Data)
	if length == width*height {
		//		log.Printf("d=%d , w=%d, h=%d, wxh=%d\n", length, width, height, width*height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				value := message.Section7.Data[y*width+x]
				r,g,b := RGBValue(value, maxValue, minValue)
				rgbaImage.Set(x, (height - y), color.NRGBA{
					R: r,
					G: g,
					B: b,
					A: uint8(254),
				})
			}
		}
	}
	return rgbaImage, nil
}


// RGBValue returns a number between 0 and 255
func RGBValue(value float64, maxValue float64, minValue float64) (uint8, uint8, uint8) {
	//value  = value - 273
	part := (maxValue - minValue) / 12
	value = value - minValue
	red := 0.0
	green := 0.0
	blue := 0.0
	if value < part * 4 {
		red = 0.0
		green = 0.0
		blue = 128.0
	} else if value < part * 5 {
		red = 0
		green = 0
		blue = 127 * math.Pow((value - (part * 4)) / part, 2) + 127
	} else if value < part * 6 {
		red = 0
		green = -127 * math.Pow((value - (part * 6)) / part, 2) + 127
		blue = 255
	} else if value < part * 7 {
		red = 0
		green = 127 * math.Pow((value - (part * 6)) / part, 2) + 127
		blue = 255
	} else if value < part * 8 {
		red = 127 * math.Pow((value - (part * 8)) / part, 2) + 127
		green = 255
		blue = 127 * math.Pow((value - (part * 8)) / part, 2) + 127
	} else if value < part * 9 {
		red = -127 * math.Pow((value - (part * 8)) / part, 2) + 127
		green = 255
		blue = -127 * math.Pow((value - (part * 8)) / part, 2) + 127
	} else if value < part * 10 {
		red = 255
		green = 127 * math.Pow((value - (part * 10)) / part, 2) + 127
		blue = 0
	} else if value < part * 11 {
		red = 255
		green = -127 * math.Pow((value - (part * 10)) / part, 2) + 127
		blue = 0
	} else if value <= part * 12 {
		red = 127 * math.Pow((value - (part * 12)) / part, 2) + 127
		green = 0
		blue = 0
	}

	return uint8(red), uint8(green), uint8(blue)
}
/*// RGBValue returns a number between 0 and 255
func RGBValue(value float64, maxValue float64, minValue float64) (uint8, uint8, uint8) {
	//value  = value - 273
	part := (maxValue - minValue) / 12
	value = value - minValue
	red := 0.0
	green := 0.0
	blue := 0.0
	if value < part * 4 {
		red = 0.0
		green = 0.0
		blue = 128.0
	} else if value < part * 5 {
		red = 0
		green = 0
		blue = ((value - part * 4) / part + 1) * 128
	} else if value < part * 6 {
		red = 0
		green = ((value - part * 5) / part) * 128
		blue = 255
	} else if value < part * 7 {
		red = 0
		green = ((value - part * 6) / part + 1) * 128
		blue = 255
	} else if value < part * 8 {
		red = ((value - part * 7) / part) * 128
		green = 255
		blue = (-(value - part * 7) / part + 2) * 128
	} else if value < part * 9 {
		red = ((value - part * 8) / part + 1) * 128
		green = 255
		blue = (-(value - part * 8) / part + 1) * 128
	} else if value < part * 10 {
		red = 255
		green = (-(value - part * 9) / part + 2) * 128
		blue = 0
	} else if value < part * 11 {
		red = 255
		green = (-(value - part * 10) / part + 1) * 128
		blue = 0
	} else if value <= part * 12 {
		red = (-(value - part * 11) / part + 2) * 128
		green = 0
		blue = 0
	}

	return uint8(red), uint8(green), uint8(blue)
}
*/
// returns a number between 0 and 255
func blueValue(value float64, maxValue float64, minValue float64) uint8 {
	//value  = value - 273
	if value < 0 {
		percentOfMaxValue := (math.Abs(value) + math.Abs(minValue)) / (math.Abs(maxValue) + math.Abs(minValue))
		return uint8(percentOfMaxValue * 255.0)
	}
	return 0
}

// RedValue returns a number between 0 and 255
func RedValue(value float64, maxValue float64, minValue float64) uint8 {
	//value  = value - 273
	len := maxValue - minValue
	if value > 0 {
		percentOfMaxValue := value / len
		return uint8(percentOfMaxValue * 255.0)
	}
	return 0
}

func MaxMin(float64s []float64) (float64, float64) {
	max, min := -9999999.0, 999999.0
	for _, v := range float64s {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return max, min
}
