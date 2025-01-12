A simple program that uses your webcam to find out the optimal brightness for your notebook's backlight.
It changes the brightness smoothly and is entirely configurable.

```sh
]~/Documents/TrulyMine/glight.v4l@ glight -h

 Copyright (c) 2025: xplshn and contributors
 For more details refer to https://github.com/xplshn/glight

  Synopsis
    glight <|--webcam [filepath](/dev/video*)|--brightness [filepath](/sys/class/backlight/*/brightness)|--max-brightness [filepath](/sys/class/backlight/*/max_brightness)|--min-brightness [1-100](10)|--set [1-100]|--max [1-100]|--scale [1-100](120)> [FILE/s]
  Description:
    Lets you controls your laptop's backlight easily
  Options:
    --brightness: Path to the brightness control file
    --every: Time interval to capture a frame and adjust brightness
    --max: Show maximum brightness value and exit
    --max-brightness: Path to the max brightness control file
    --min-brightness: Minimum brightness percentage (1-100)
    --scale: Scale factor for brightness transition
    --set: Set brightness directly (1-100)
    --webcam: Path to the webcam device

]~/Documents/TrulyMine/glight.v4l@
```
