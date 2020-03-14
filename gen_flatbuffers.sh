flatc --ts --no-fb-import --short-names -o /Users/alexakimov/marsgame/marsgame-client/src/flatbuffers flatbuffers/schemes/*.fbs
flatc --go -o flatbuffers/generated flatbuffers/schemes/*.fbs