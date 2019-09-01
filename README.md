# Tuile  
  
Tuile (french for tile) is a 2D graphics engine inspired from old
hardware and based on layers, tiles sets, tile maps and sprites. Its
scanline rendering pipeline makes it perfect for raster effects.

## Just a rendering pipeline

Tuile rendering pipeline outputs everything to a framebuffer. You can
therefore use your favorite game engine to display the resulting frames.
Note that this repository uses [Ebiten](https://ebiten.org/) for all the
samples.

## Samples

You can check the available samples in the `samples` directory. For
instance, you can run the curvature sample using the following commands:

```
cd samples/curvature
go run .
``` 
