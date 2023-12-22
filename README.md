# cic

cic (pronounced *kick*, *sis* or *sick* [as you
prefer](https://en.wikipedia.org/wiki/Hard_and_soft_C)) is an experimental
colouring-in creator.

cic transforms an image – whether a drawing, a poster, or a screenshot from a TV
show – into a colouring sheet. It is a work in process (hence experimental
above), and though fully functional, it is hoped to refine the quality of output
for colouring sheets over time.

cic uses a [Canny edge
detector](https://en.wikipedia.org/wiki/Canny_edge_detector) to perform edge
detection, though future versions may adapt this and add alternative methods.

cic is implemented in [Go](https://go.dev/), and uses
[Cobra](https://cobra.dev/) to provide a CLI. I hope to have a web version soon,
to enable use across platforms and by non-commmandline literate users.

## Motivation

Duck Duck Go (other search engines are available and have the same limitations)
only brings up a limited number of free and good colouring-in pictures of
tractors, helicopters, and other interesting objects, and my son keeps asking
for new ones. I'd also like to be able to make him colouring sheets of his
favourite characters from TV shows he likes. These are for personal use only,
with no intent to distribute, and no desire to infringe copyright of those
making creative and educational content for kids.

Of course, an alternative method would be to use generative AI to produce
colouring sheets of a description provided in a prompt. But that isn't the
approach of this project.

## Compilation

    go build -o cic main.go
    
## Usage

    cic [flags] filename
    
Valid flags are:
- `-h`, `--help`: print help
- `-s`, `--stddev float`: Standard deviation of Gaussian blur (see below).
  Default is 1.0.
- `-l`, `--lower int`: Lower threshold for edge suppression (see below). Default
  is 10.
- `-u`, `--upper int`: Upper threshold for edge suppression (see below). Default
  is 100.
  
## Parameters and tuning

As in any edge-detection problem, creating a colouring sheet requires finding
and including the edges we want, while ignoring image edges which display
texture or other noise, in order to create a clean, shape based colouring sheet.
This is a difficult problem in all computer vision tasks, and I hope to improve
cic's performance for colouring sheets.

Among the techniques used in a Canny edge detector for distinguishing between
signal and noise, two have tunable parameters in cic: Gaussian blurring and
threshold-based suppression.

Gaussian blurring is a low-pass filter applied to the original image, before
edge detection. The standard deviation of a [Gaussian
distribution](https://en.wikipedia.org/wiki/Normal_distribution) determines how
wide the blur is, and therefore how much smoothing is applied. This can be set
by use of the `-s` flag (or long version `--stddev`). A higher number will blur
the image more before performing edge detection, which will reduce the noise of
textures and fine details. A lower number will provide less blurring, and so
edges will be retained with greater sharpness. Setting the standard deviation to
zero will result in no blurring applied, which is recommended only for images
with no texture or shading.

Threshold-based suppression is performed after edge detection, and aims to keep
significant edges while discarding noise and other insignificant edges. Applying
a single threshold is not effective, because an important edge can vary in
sharpness over its length, and the smoother parts of the edge would be discarded
by a simple threshold. Instead, two thresholds are used. Edges which are very
smooth fall below the lower threshold and are discarded as noise. Edges which
are very sharp lie above the upper threshold and are retained. Edges which fall
between the lower and upper threshold are kept if they are connected to a sharp
edge above the upper threshold, but discarded if they are not connected to a
sharp edge. The tuning of the upper (`-u` or long version `--upper`)and lower
(`-l` or long version `--lower`) thresholds are quite sensitive, and the best
results vary quite widely from one image to another. If the default values do
not give good results, trial and error is required to tune these values to
improve results -- though they are limited in what they can do.

Future versions of cic will aim to incorporate more methods for retaining
important edges and discarding noise.
