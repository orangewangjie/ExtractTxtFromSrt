# ExtractTxtFromSrt
Subtitle files in videos contain valuable information, particularly when you're learning new skills. I've created a software tool to extract pure content from .srt subtitle files, and I hope it can help others with similar needs. I plan to make improvements when I have more time. If you have any suggestions, please feel free to share them.

The main function of the software is to extract pure content from subtitle files (.srt). You can use the command below to extract .txt files from a specified folder (e.g., ./aws):
```
./extract.exe -f ./aws
```
The extracted content will be saved in text format. Additionally, the program works even if the folder is nested.
