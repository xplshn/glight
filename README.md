# glight

**glight** is a simple yet powerful tool designed to automatically adjust your laptop's backlight brightness based on the ambient light detected by your webcam. It ensures a smooth transition between brightness levels and offers extensive configuration options to suit your preferences.

## Features

- **Automatic Brightness Adjustment**: Uses your webcam to detect ambient light and adjust the screen brightness accordingly.
- **Smooth Transitions**: Gradual changes in brightness levels to avoid abrupt shifts.

## Usage

To use **glight**, ensure that your system has the necessary permissions and devices available:

- `/dev/video*` should be available and writable.
- `/sys/class/backlight/*/brightness` should be writable and `/sys/class/backlight/*/max_brightness` should be accessible.

Here is an example of how you can set up these permissions in your `/etc/rc.local`:

```sh
#!/bin/sh

# For AppBundle/AppImage/FlatImage support
modprobe fuse

# Change group ownership to 'video'
chown root.video /dev/video*
chown root.video /sys/class/backlight/*/brightness

# Set group write permissions
chmod g+rw /dev/video*
chmod g+rw /sys/class/backlight/*/brightness
```

Of course, you may do this using an mdev script, or an smdev helper, or a udev rule, etc. There are many ways to do this.

### Command-Line Options

```sh
]~/Documents/TrulyMine/glight.v4l@ glight -h

 Copyright (c) 2025: xplshn and contributors
 For more details refer to https://github.com/xplshn/glight

  Synopsis
    glight <|--webcam [filepath](/dev/video*)|--brightness [filepath](/sys/class/backlight/*/brightness)|--max-brightness [filepath](/sys/class/backlight/*/max_brightness)|--min-brightness [1-100](10)|--set [1-100]|--max [1-100]|--scale [1-100](120)> [FILE/s]
  Description:
    Lets you control your laptop's backlight easily
  Options:
    --brightness: Path to the brightness control file (/sys/class/backlight/*/brightness)
    --every: Time interval to capture a frame and adjust brightness (30s)
    --max: Show maximum brightness value and exit
    --max-brightness: Path to the max brightness control file (/sys/class/backlight/*/max_brightness)
    --min-brightness: Minimum brightness percentage [1-100] (10)
    --scale: Scale factor for brightness transition (120)
    --set: Set brightness directly [1-100]
    --webcam: Path to the webcam device (/dev/video*)
```

## Installation

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/xplshn/glight.git
   cd glight
   ```

2. **Build the Project**:
   ```sh
   go build -o glight glight.go
   ```

3. **Run the Program**:
   ```sh
   ./glight &
   ```

## Contributing

Contributions are welcome! If you have any ideas, suggestions, or bug reports, please open an issue or submit a pull request.

#### TODO: Add a config file, to make this portable to freeBSD, openBSD and other platforms.

## License

**scriptfs** is licensed under the following licenses: Choose whichever fits your needs best:
- ISC (Pre-2007 ISC License, the one the OpenBSD project uses)
- MIT-0 (MIT ZERO attribution)
- Unlicense (The Unlicense)
