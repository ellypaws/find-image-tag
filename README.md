# Find Image Captions

Find Image Captions is a collection of tools designed to organize caption files for images. It is built with Go and uses
the `roggy` and `stemp` libraries for logging and string templating respectively. The project also uses
the `charmbracelet/bubbles` library for creating terminal user interfaces.

## Features

- Finds matching .txt files for a group of images in a directory.
- Provides options to move, copy, or hardlink the matched caption files.
- Displays detailed statistics and actions in a terminal user interface.
- Allows for the addition of new files to the dataset.
- Supports checking if each image has a caption.
- Provides options to append text files to matching images, merge new captions to existing ones, and replace spaces with
  underscores in captions.

## Usage

To use the Find Image Captions tool, you need to have Go installed on your machine. Once you have Go installed, you can
clone this repository and run the main.go file.

```bash
git clone https://github.com/ellypaws/find-image-tag.git
cd find-image-captions
go run main.go
```