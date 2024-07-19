# inkwell
The official source code repository for the inkwell manuscript compiler

## Installation
To install inkwell, simply run the following command in your terminal:

```bash
go install github.com/inkwell/inkwell
```

## Usage
To use inkwell, simply run the following command in your terminal:

```bash
./inkwell --config=.inkwell.yaml
```

## Configuration
To configure inkwell, create a file named `.inkwell.yaml` in the root of your project. Here is an example configuration file:

```yaml
# .inkwell.yaml
title: Your Project Name
authors: 
  - Author 1
  - Author 2
summary: A summary of your project
output_filename: path/to/full-manuscript.md
summary_filename: path/to/summary.md
chapters:
  - title: Chapter 1
    output_filename: path/to/chapter-01.md # optional
    number_paragraphs: true # adds <#> to the end of each paragraph in the chapter output, but not the full manuscript output
    scenes:
      - files: # First Scene
        - path/to/chapter 1/file 1.md
        - path/to/chapter 1/file 2.md
      - files: # Second Scene
        - path/to/chapter 1/file 3.md
  - title: Chapter 2
    output_filename: path/to/chapter-02.md # optional
    number_paragraphs: true
    scenes:
      - files:
        - path/to/chapter 2/file 1.md
```

## Contributing
To contribute to inkwell, simply fork this repository, make your changes, and submit a pull request.

## License
This license is licensed under GPL-3.0. For more information, see the [license](LICENSE.md) file.
