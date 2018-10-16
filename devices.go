package goinsta

import (
	"fmt"
	"github.com/alok87/goutils/pkg/random"
	"strings"
)


const (
	UserLocale = "en_US"
	AppVersion = "27.0.0.7.97"
	UserAgentFormat = "Instagram %s Android (%s/%s; %s; %s; %s; %s; %s; %s; %s)"
	UserAgentFormatIos = "Instagram %s %s"
)

var iosDevices = []string {

	"(iPhone7,1; iOS 10_2; en_RS; en-RS; scale=2.00; gamut=normal; 640x1136) AppleWebKit/420+",

	"(iPhone7,1; iOS 10_2; en_RS; en-RS; scale=2.61; gamut=normal; 1080x1920) AppleWebKit/420+",

	"(iPhone9,2; iOS 11_2_1; en_RS; en-RS; scale=2.61; gamut=wide; 1080x1920) AppleWebKit/420+",

	"(iPhone9,4; iOS 11_0; en_RS; en-RS; scale=2.61; gamut=wide; 1080x1920) AppleWebKit/420+",

	"(iPhone9,3; iOS 11_2_2; en_RS; en-RS; scale=2.00; gamut=wide; 750x1334) AppleWebKit/420+",

	"(iPhone6,1; iOS 10_3_2; en_RS; en-RS; scale=2.00; gamut=normal; 640x1136) AppleWebKit/420+",

	"(iPhone8,1; iOS 11_3_1; en_RS; en-RS; scale=2.00; gamut=normal; 750x1334) AppleWebKit/420+",

	"(iPhone8,1; iOS 11_2_1; en_RS; en-RS; scale=2.00; gamut=normal; 750x1334) AppleWebKit/420+",

	"(iPhone7,1; iOS 9_3_1; en_RS; en-RS; scale=2.88; gamut=normal; 1080x1920) AppleWebKit/420+",

	"(iPhone8,4; iOS 11_2_6; en_RS; en-RS; scale=2.00; gamut=normal; 640x1136) AppleWebKit/420+",
}


var androidDevices = []string {
	/* OnePlus 3T. Released: November 2016.
	 * https://www.amazon.com/OnePlus-A3010-64GB-Gunmetal-International/dp/B01N4H00V8
	 * https://www.handsetdetection.com/properties/devices/OnePlus/A3010
	 */
	"24/7.0; 380dpi; 1080x1920; OnePlus; ONEPLUS A3010; OnePlus3T; qcom",

	/* LG G5. Released: April 2016.
	 * https://www.amazon.com/LG-Unlocked-Phone-Titan-Warranty/dp/B01DJE22C2
	 * https://www.handsetdetection.com/properties/devices/LG/RS988
	 */
	"23/6.0.1; 640dpi; 1440x2392; LGE/lge; RS988; h1; h1",

	/* Huawei Mate 9 Pro. Released: January 2017.
	 * https://www.amazon.com/Huawei-Dual-Sim-Titanium-Unlocked-International/dp/B01N9O1L6N
	 * https://www.handsetdetection.com/properties/devices/Huawei/LON-L29
	 */
	"24/7.0; 640dpi; 1440x2560; HUAWEI; LON-L29; HWLON; hi3660",

	/* ZTE Axon 7. Released: June 2016.
	 * https://www.frequencycheck.com/models/OMYDK/zte-axon-7-a2017u-dual-sim-lte-a-64gb
	 * https://www.handsetdetection.com/properties/devices/ZTE/A2017U
	 */
	"23/6.0.1; 640dpi; 1440x2560; ZTE; ZTE A2017U; ailsa_ii; qcom",

	/* Samsung Galaxy S7 Edge SM-G935F. Released: March 2016.
	 * https://www.amazon.com/Samsung-SM-G935F-Factory-Unlocked-Smartphone/dp/B01C5OIINO
	 * https://www.handsetdetection.com/properties/devices/Samsung/SM-G935F
	 */
	"23/6.0.1; 640dpi; 1440x2560; samsung; SM-G935F; hero2lte; samsungexynos8890",

	/* Samsung Galaxy S7 SM-G930F. Released: March 2016.
	 * https://www.amazon.com/Samsung-SM-G930F-Factory-Unlocked-Smartphone/dp/B01J6MS6BC
	 * https://www.handsetdetection.com/properties/devices/Samsung/SM-G930F
	 */
	"23/6.0.1; 640dpi; 1440x2560; samsung; SM-G930F; herolte; samsungexynos8890",
}


func getAndroidUserAgent() string {

	devices := getRandomAndroidDevice()

	parts := strings.Split(devices, "; ")

	// If part is less than 7, return default user agent.
	if len(parts) <= 7 {
		return goInstaUserAgent
	}

	androidOS := strings.Split(parts[0], "/")


	manufacturerAndBrand := strings.Split(parts[3], "/" )

	androidVersion := androidOS[0]
	androidRelease := androidOS[1]

	dpi := parts[1]
	resolution := parts[2]
	manufacturer := manufacturerAndBrand[0]
	model := parts[4]
	device := parts[5]
	cpu := parts[6]

	return fmt.Sprintf(UserAgentFormat, AppVersion, androidVersion, androidRelease, dpi, resolution,manufacturer,
		model, device, cpu, UserLocale)
}

func getRandomAndroidDevice() string {
	random := random.RangeInt(0, len(androidDevices), 1)[0]
	return androidDevices[random]
}

func getRandomIosUserAgent() string {
	random := random.RangeInt(0, len(iosDevices), 1)[0]
	return iosDevices[random]
}

func getIosUserAgent() string {
	return fmt.Sprintf(UserAgentFormatIos,AppVersion, getRandomIosUserAgent())
}

// Returns UserAgent and Device
func UserAgentAndDevice(username string, password string) (string, string){

	choice := random.RangeInt(1,200, 1)[0] % 2

	// Return user agent for ios.
	if choice == 0 {
		return getIosUserAgent(), generateUUID()
	}else // Return user agent for android
	{
		return getAndroidUserAgent(), generateDeviceID(
			generateMD5Hash(username + password),
		)
	}
}