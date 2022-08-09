package main

import (
	"io"
	gokmp "leslack/src/kmp"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// var FileNames = [250]string{"000_001.txt", "001_002.txt", "002_003.txt", "003_004.txt", "004_005.txt", "005_006.txt", "006_007.txt", "007_008.txt", "008_009.txt", "009_010.txt", "010_011.txt", "011_012.txt", "012_013.txt", "013_014.txt", "014_015.txt", "015_016.txt", "016_017.txt", "017_018.txt", "018_019.txt", "019_020.txt", "020_021.txt", "021_022.txt", "022_023.txt", "023_024.txt", "024_025.txt", "025_026.txt", "026_027.txt", "027_028.txt", "028_029.txt", "029_030.txt", "030_031.txt", "031_032.txt", "032_033.txt", "033_034.txt", "034_035.txt", "035_036.txt", "036_037.txt", "037_038.txt", "038_039.txt", "039_040.txt", "040_041.txt", "041_042.txt", "042_043.txt", "043_044.txt", "044_045.txt", "045_046.txt", "046_047.txt", "047_048.txt", "048_049.txt", "049_050.txt", "050_051.txt", "051_052.txt", "052_053.txt", "053_054.txt", "054_055.txt", "055_056.txt", "056_057.txt", "057_058.txt", "058_059.txt", "059_060.txt", "060_061.txt", "061_062.txt", "062_063.txt", "063_064.txt", "064_065.txt", "065_066.txt", "066_067.txt", "067_068.txt", "068_069.txt", "069_070.txt", "070_071.txt", "071_072.txt", "072_073.txt", "073_074.txt", "074_075.txt", "075_076.txt", "076_077.txt", "077_078.txt", "078_079.txt", "079_080.txt", "080_081.txt", "081_082.txt", "082_083.txt", "083_084.txt", "084_085.txt", "085_086.txt", "086_087.txt", "087_088.txt", "088_089.txt", "089_090.txt", "090_091.txt", "091_092.txt", "092_093.txt", "093_094.txt", "094_095.txt", "095_096.txt", "096_097.txt", "097_098.txt", "098_099.txt", "099_100.txt", "100_101.txt", "101_102.txt", "102_103.txt", "103_104.txt", "104_105.txt", "105_106.txt", "106_107.txt", "107_108.txt", "108_109.txt", "109_110.txt", "110_111.txt", "111_112.txt", "112_113.txt", "113_114.txt", "114_115.txt", "115_116.txt", "116_117.txt", "117_118.txt", "118_119.txt", "119_120.txt", "120_121.txt", "121_122.txt", "122_123.txt", "123_124.txt", "124_125.txt", "125_126.txt", "126_127.txt", "127_128.txt", "128_129.txt", "129_130.txt", "130_131.txt", "131_132.txt", "132_133.txt", "133_134.txt", "134_135.txt", "135_136.txt", "136_137.txt", "137_138.txt", "138_139.txt", "139_140.txt", "140_141.txt", "141_142.txt", "142_143.txt", "143_144.txt", "144_145.txt", "145_146.txt", "146_147.txt", "147_148.txt", "148_149.txt", "149_150.txt", "150_151.txt", "151_152.txt", "152_153.txt", "153_154.txt", "154_155.txt", "155_156.txt", "156_157.txt", "157_158.txt", "158_159.txt", "159_160.txt", "160_161.txt", "161_162.txt", "162_163.txt", "163_164.txt", "164_165.txt", "165_166.txt", "166_167.txt", "167_168.txt", "168_169.txt", "169_170.txt", "170_171.txt", "171_172.txt", "172_173.txt", "173_174.txt", "174_175.txt", "175_176.txt", "176_177.txt", "177_178.txt", "178_179.txt", "179_180.txt", "180_181.txt", "181_182.txt", "182_183.txt", "183_184.txt", "184_185.txt", "185_186.txt", "186_187.txt", "187_188.txt", "188_189.txt", "189_190.txt", "190_191.txt", "191_192.txt", "192_193.txt", "193_194.txt", "194_195.txt", "195_196.txt", "196_197.txt", "197_198.txt", "198_199.txt", "199_200.txt", "200_201.txt", "201_202.txt", "202_203.txt", "203_204.txt", "204_205.txt", "205_206.txt", "206_207.txt", "207_208.txt", "208_209.txt", "209_210.txt", "210_211.txt", "211_212.txt", "212_213.txt", "213_214.txt", "214_215.txt", "215_216.txt", "216_217.txt", "217_218.txt", "218_219.txt", "219_220.txt", "220_221.txt", "221_222.txt", "222_223.txt", "223_224.txt", "224_225.txt", "225_226.txt", "226_227.txt", "227_228.txt", "228_229.txt", "229_230.txt", "230_231.txt", "231_232.txt", "232_233.txt", "233_234.txt", "234_235.txt", "235_236.txt", "236_237.txt", "237_238.txt", "238_239.txt", "239_240.txt", "240_241.txt", "241_242.txt", "242_243.txt", "243_244.txt", "244_245.txt", "245_246.txt", "246_247.txt", "247_248.txt", "248_249.txt", "249_250.txt"}
var FileNames = [11]string{"000_001.txt", "001_002.txt", "002_003.txt", "003_004.txt", "004_005.txt", "005_006.txt", "006_007.txt", "007_008.txt", "008_009.txt", "009_010.txt", "010_011.txt"}

const Limit = 16

var findChannel [Limit]chan Patterns
var kmp *gokmp.KMP
var wg *sync.WaitGroup

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("开始时间")
	pattern := "13908262671"
	kmp, _ = gokmp.NewKMP(pattern)
	for i := 0; i < Limit; i++ {
		findChannel[i] = createWorker()
	}
	wg = new(sync.WaitGroup)
	for i := range FileNames {
		file := "/Users/leslack/Downloads/pai_250/" + FileNames[i]
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		j := 0
		end := make([]byte, len(pattern))
		for {
			buf := make([]byte, 1024)
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}
			wg.Add(1)
			findChannel[rand.Intn(Limit)] <- Patterns{strs: string(append(end, buf[:n]...)), base: i*100000000 + j - len(end)}
			end = buf[1024-len(pattern):]
			j += 1024
		}
	}
	wg.Wait()
	log.Println("结束时间")
}

type Patterns struct {
	strs string
	base int
}

func createWorker() chan Patterns {
	worker := make(chan Patterns)

	go func() {
		for {
			result := <-worker
			//index := strings.Index(result, "13594796150")
			//if index != -1 {
			ints := kmp.FindAllStringIndex(result.strs)
			if len(ints) > 0 {
				log.Println(result.base + ints[0] + 1)
				log.Println("结束时间")
				os.Exit(1)
			}
			wg.Done()
		}
	}()

	return worker
}
