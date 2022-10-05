# Markdown style formatter

This is a very primitive markdown title tag style formatter written in go as a challenge for myself and to quickly address the issue of non-standardized formatting in markdown documentation. 

It was originally developed to address [this issue](https://github.com/hasura/graphql-engine/issues/8885) and therefore follows their chosen styling requirements of Title Case for headings and Sentence case for subheadings. 

This program could definitely be better optimized and abstracted to allow for more control over the functionality but this can be left as a future exercise. 

# Usage

To use the program run as below where filePath is the directory containing your mdx files.

```
go run . <filePath>
```
