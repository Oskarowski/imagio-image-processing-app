# Flow Diagram for Bandpass Filter

## Applying a _bandpass_ filter to an image:

![Flow Diagram for Bandpass Filter](/assets/bandpass_filter_flow.png)

```D2
B: "Convert Image to Complex"
D: "Compute 2D FFT"
DC: "Extract DC Component"
E: "Swap Quadrants"
I: "Apply Bandpass Filter"
J: "Swap Quadrants Back"
R: "Restore DC Component"
N: "Reconstruct Image with Inverse FFT"
O: "Convert Reconstructed Image to BMP"

B -> D
D -> DC
D -> E
E -> I
I -> J
J -> R
R -> N
N -> O
```

## Applying a _lowpass_ filter to an image:

![Flow Diagram for Bandpass Filter](/assets/lowpass_filter_flow.png)

```D2
B: "Convert Image to Complex"
D: "Compute 2D FFT"
DC: "Extract DC Component"
E: "Swap Quadrants"
LP: "Apply Lowpass Filter"
J: "Swap Quadrants Back"
R: "Restore DC Component"
N: "Reconstruct Image with Inverse FFT"
O: "Convert Reconstructed Image to BMP"

B -> D
D -> DC
D -> E
E -> LP
LP -> J
J -> R
R -> N
N -> O
```


## Applying a _Bandpass with edge detection_ filter to an image:

![Flow Diagram for Bandpass Edge Detection Filter](/assets/bandpass_with_edge_detection_filter_flow.png)

```D2
input_image: {
  shape: oval
  text: "Input Frequency Spectrum\n(n x m)"
}

input_mask: {
  shape: oval
  text: "Binary Mask\n(p x q)"
}

scale_mask: {
  shape: parallelogram
  text: "Scale Binary Mask\n(to n x m)"
}

filter_process: {
  shape: diamond
  text: "Apply Filter:\nIf scaledMask(x, y) == 0,\nset spectrum(x, y) = 0\nElse retain spectrum(x, y)"
}

output_spectrum: {
  shape: oval
  text: "Filtered Frequency Spectrum\n(n x m)"
}

input_image -> scale_mask: "Image Spectrum"
input_mask -> scale_mask: "Mask"
scale_mask -> filter_process: "Scaled Mask"
filter_process -> output_spectrum: "Filtered Spectrum"

```
